package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	discordWebhookURL string
	discordDB         *gorm.DB
	discordHTTPClient = &http.Client{Timeout: 5 * time.Second}
	jakartaLoc        *time.Location

	hariIndo = map[time.Weekday]string{
		time.Sunday: "Minggu", time.Monday: "Senin", time.Tuesday: "Selasa",
		time.Wednesday: "Rabu", time.Thursday: "Kamis", time.Friday: "Jumat", time.Saturday: "Sabtu",
	}
	bulanIndo = map[time.Month]string{
		time.January: "Januari", time.February: "Februari", time.March: "Maret",
		time.April: "April", time.May: "Mei", time.June: "Juni", time.July: "Juli",
		time.August: "Agustus", time.September: "September", time.October: "Oktober",
		time.November: "November", time.December: "Desember",
	}
)

func init() {
	jakartaLoc, _ = time.LoadLocation("Asia/Jakarta")
	if jakartaLoc == nil {
		jakartaLoc = time.FixedZone("WIB", 7*60*60)
	}
}

// InitDiscordAlert mengaktifkan notifikasi error ke Discord webhook. Panggil
// sekali saat boot (internal/boot/http.go), setelah koneksi DB dibuka.
// Kalau webhookURL kosong (mis. environment development), notifikasi otomatis
// nonaktif -- aman dipanggil di semua environment tanpa efek samping.
func InitDiscordAlert(db *gorm.DB, webhookURL string) {
	discordDB = db
	discordWebhookURL = strings.TrimSpace(webhookURL)
}

func formatJakartaTime(t time.Time) string {
	t = t.In(jakartaLoc)
	return fmt.Sprintf("%s, %d %s %d - %02d:%02d:%02d WIB",
		hariIndo[t.Weekday()], t.Day(), bulanIndo[t.Month()], t.Year(),
		t.Hour(), t.Minute(), t.Second())
}

type discordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// truncate memotong string supaya tidak melebihi batas panjang field embed Discord.
func truncate(s string, max int) string {
	if s == "" {
		return "-"
	}
	if len(s) <= max {
		return s
	}
	return s[:max] + "\n...(terpotong)"
}

// sendDiscordEmbed mengirim satu embed ke Discord webhook secara async --
// tidak pernah memblokir atau menggagalkan request yang sedang berjalan.
func sendDiscordEmbed(title string, color int, fields []discordField) {
	if discordWebhookURL == "" {
		return
	}
	go func() {
		defer func() { recover() }() // kegagalan kirim notifikasi tidak boleh crash apa pun

		payload := map[string]interface{}{
			"embeds": []map[string]interface{}{
				{"title": title, "color": color, "fields": fields},
			},
		}
		body, err := json.Marshal(payload)
		if err != nil {
			return
		}
		req, err := http.NewRequest(http.MethodPost, discordWebhookURL, bytes.NewReader(body))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := discordHTTPClient.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
	}()
}

// lookupUserName cari nama+email user dari gold_id (best-effort, timeout 2s,
// tidak pernah menggagalkan alert kalau lookup gagal).
func lookupUserName(goldID int) string {
	if discordDB == nil || goldID == 0 {
		return "Tidak terautentikasi (tanpa token)"
	}
	var row struct {
		GoldNama  string
		GoldEmail string
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := discordDB.WithContext(ctx).
		Table("data_peserta").
		Select("gold_nama, gold_email").
		Where("gold_id = ?", goldID).
		Take(&row).Error
	if err != nil {
		return fmt.Sprintf("gold_id=%d (nama tidak ditemukan di DB)", goldID)
	}
	return fmt.Sprintf("%s (%s) — gold_id=%d", row.GoldNama, row.GoldEmail, goldID)
}

// lookupOutletName cari nama outlet dari kode outlet (best-effort, timeout 2s).
func lookupOutletName(outcode string) string {
	if outcode == "" {
		return "-"
	}
	if discordDB == nil {
		return outcode
	}
	var row struct {
		OutletName string
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := discordDB.WithContext(ctx).
		Table("outlet").
		Select("outlet_name").
		Where("outlet_code = ?", outcode).
		Take(&row).Error
	if err != nil {
		return fmt.Sprintf("%s (nama tidak ditemukan di DB)", outcode)
	}
	return fmt.Sprintf("%s (%s)", row.OutletName, outcode)
}

// extractOutletCode cari kode outlet dari query param umum (code/outcode),
// atau kalau tidak ada, dari potongan body JSON yang sudah di-buffer.
func extractOutletCode(c *gin.Context, bodySnippet string) string {
	if v := c.Query("code"); v != "" {
		return v
	}
	if v := c.Query("outcode"); v != "" {
		return v
	}
	for _, key := range []string{`"outcode"`, `"order_outcode"`, `"code"`} {
		idx := strings.Index(bodySnippet, key)
		if idx == -1 {
			continue
		}
		rest := strings.TrimLeft(bodySnippet[idx+len(key):], ": \"")
		if end := strings.IndexAny(rest, `",}`); end != -1 && end > 0 {
			return rest[:end]
		}
	}
	return ""
}

// bufferRequestBody membaca body request untuk keperluan alert, lalu
// mengembalikan body itu utuh supaya handler di belakang tetap bisa
// membacanya seperti biasa. Dilewati untuk multipart/form-data (upload foto
// bukti transfer) supaya tidak membengkakkan memory tiap request.
func bufferRequestBody(c *gin.Context) string {
	if c.Request.Body == nil || c.Request.Body == http.NoBody {
		return ""
	}
	if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
		return ""
	}
	if c.Request.ContentLength > 1<<20 { // >1MB, lewati (bukan payload JSON normal)
		return ""
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}
	c.Request.Body.Close()
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return truncate(string(bodyBytes), 1500)
}

// ErrorAlert mem-buffer body request lalu, kalau response akhirnya berstatus
// >=500, mengirim detail lengkap (user, outlet, endpoint, waktu, body, dsb)
// ke Discord. Panic ditangani terpisah lewat AlertPanic (lihat handler.go)
// karena panic melompati kode setelah c.Next() di sini.
func ErrorAlert() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		bodySnippet := bufferRequestBody(c)

		c.Next()

		status := c.Writer.Status()
		if status < 500 {
			return
		}

		errMsg := ""
		if len(c.Errors) > 0 {
			errMsg = c.Errors.String()
		}

		buildAndSendAlert(c, status, errMsg, "", bodySnippet, time.Since(start))
	}
}

// AlertPanic dipanggil dari dalam CustomRecovery (handler.go) karena panic
// melewati alur normal middleware -- termasuk kode setelah c.Next() di
// ErrorAlert -- sehingga perlu jalur pengiriman tersendiri.
func AlertPanic(c *gin.Context, recovered interface{}) {
	stack := truncate(string(debug.Stack()), 1000)
	buildAndSendAlert(c, http.StatusInternalServerError, fmt.Sprintf("%v", recovered), stack, "", 0)
}

func buildAndSendAlert(c *gin.Context, status int, errMsg, stack, bodySnippet string, latency time.Duration) {
	path := c.FullPath()
	if path == "" {
		path = c.Request.URL.Path
	}
	outcode := extractOutletCode(c, bodySnippet)

	var goldID int
	if v, ok := c.Get("user"); ok {
		if id, ok := v.(int); ok {
			goldID = id
		}
	}
	role := "-"
	if v, ok := c.Get("role"); ok {
		if r, ok := v.(string); ok {
			role = r
		}
	}

	title := fmt.Sprintf("🔴 Error %d - %s %s", status, c.Request.Method, path)

	fields := []discordField{
		{Name: "Waktu", Value: formatJakartaTime(time.Now()), Inline: false},
		{Name: "Endpoint", Value: truncate(fmt.Sprintf("%s %s?%s", c.Request.Method, path, c.Request.URL.RawQuery), 500), Inline: false},
		{Name: "Status", Value: fmt.Sprintf("%d", status), Inline: true},
		{Name: "Environment", Value: environment, Inline: true},
		{Name: "Latency", Value: latency.String(), Inline: true},
		{Name: "User", Value: truncate(lookupUserName(goldID), 500), Inline: false},
		{Name: "Role", Value: role, Inline: true},
		{Name: "Outlet", Value: truncate(lookupOutletName(outcode), 500), Inline: true},
		{Name: "Client IP", Value: c.ClientIP(), Inline: true},
		{Name: "Error", Value: fmt.Sprintf("```%s```", truncate(errMsg, 900)), Inline: false},
	}
	if bodySnippet != "" {
		fields = append(fields, discordField{Name: "Request Body", Value: fmt.Sprintf("```json\n%s\n```", truncate(bodySnippet, 900)), Inline: false})
	}
	if stack != "" {
		fields = append(fields, discordField{Name: "Stack Trace (panic)", Value: fmt.Sprintf("```%s```", stack), Inline: false})
	}

	sendDiscordEmbed(title, 15158332, fields)
}
