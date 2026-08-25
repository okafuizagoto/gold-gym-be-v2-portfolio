package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Rate limiter per-IP untuk endpoint yang memicu operasi argon2 (hash saat
// registrasi, compare saat login), supaya lonjakan request bersamaan tidak
// menjadi vektor penolakan-layanan (insiden 2026-08-18: 50 pendaftaran
// bersamaan -> kedua pod production restart) maupun celah brute-force
// password. Ini hanya membatasi laju masuk, bukan kekuatan hash/compare --
// tidak menyentuh keamanan password sama sekali.
type rateEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	rlMu    sync.Mutex
	rlByKey = make(map[string]*rateEntry)
)

func init() {
	go func() {
		for range time.Tick(10 * time.Minute) {
			rlMu.Lock()
			for k, e := range rlByKey {
				if time.Since(e.lastSeen) > 10*time.Minute {
					delete(rlByKey, k)
				}
			}
			rlMu.Unlock()
		}
	}()
}

func allow(key string, every time.Duration, burst int) bool {
	rlMu.Lock()
	defer rlMu.Unlock()

	e, ok := rlByKey[key]
	if !ok {
		e = &rateEntry{limiter: rate.NewLimiter(rate.Every(every), burst)}
		rlByKey[key] = e
	}
	e.lastSeen = time.Now()
	return e.limiter.Allow()
}

func tooManyRequests(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
		"error": gin.H{"status": true, "msg": msg, "code": http.StatusTooManyRequests},
	})
}

// RegistrationRateLimit membatasi type=registerbuyer & type=insertuser:
// 3 percobaan/menit per IP (burst 3, isi ulang tiap 20s).
func (h *Handler) RegistrationRateLimit(c *gin.Context) {
	t := c.Query("type")
	if t != "registerbuyer" && t != "insertuser" {
		c.Next()
		return
	}

	if !allow("register:"+c.ClientIP(), 20*time.Second, 3) {
		tooManyRequests(c, "Terlalu banyak percobaan pendaftaran dari alamat ini, coba lagi dalam beberapa saat")
		return
	}
	c.Next()
}

// LoginRateLimit membatasi percobaan login (type=loginuser di endpoint
// bersama, maupun route /login terpisah): 8 percobaan/menit per IP (burst 8,
// isi ulang tiap 8s) -- lebih longgar dari registrasi karena pengguna asli
// wajar salah ketik password beberapa kali. Melindungi dari brute-force
// password sekaligus dari lonjakan CompareHash (argon2) yang bisa
// menghabiskan CPU pod seperti insiden registrasi.
func (h *Handler) LoginRateLimit(c *gin.Context) {
	isLoginRoute := strings.HasSuffix(c.FullPath(), "/login") || c.Query("type") == "loginuser"
	if !isLoginRoute {
		c.Next()
		return
	}

	if !allow("login:"+c.ClientIP(), 8*time.Second, 8) {
		tooManyRequests(c, "Terlalu banyak percobaan login dari alamat ini, coba lagi dalam beberapa saat")
		return
	}
	c.Next()
}

// SensitiveUserActionRateLimit membatasi PUT & DELETE /v2/userdata (ubah
// password/nama/kartu via OTP, hapus subscription, dll) -- semuanya publik
// tanpa token. OTP reset password 6 digit (1 juta kemungkinan) tidak pernah
// kedaluwarsa di skema saat ini, jadi tanpa rate-limit bisa di-brute-force.
// 5 percobaan/menit per IP membuat brute-force perlu >100 hari, praktis
// tidak mungkin, tanpa mengubah alur OTP itu sendiri.
func (h *Handler) SensitiveUserActionRateLimit(c *gin.Context) {
	if !allow("sensitive:"+c.ClientIP(), 12*time.Second, 5) {
		tooManyRequests(c, "Terlalu banyak percobaan dari alamat ini, coba lagi dalam beberapa saat")
		return
	}
	c.Next()
}
