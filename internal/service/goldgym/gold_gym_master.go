package goldgym

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"gold-gym-be/internal/entity/auth/v2"
	goldAuthEntity "gold-gym-be/internal/entity/auth/v2"
	goldEntity "gold-gym-be/internal/entity/goldgym"
	goldTokenEntity "gold-gym-be/internal/entity/token"
	"gold-gym-be/pkg/errors"
	"gold-gym-be/pkg/response"
	"log"
	"math"
	"math/rand"
	"os"

	// "os"
	"strconv"
	"strings"
	"time"

	// "gold-gym-be/internal/entity/auth/v2"

	// "github.com/dgrijalva/jwt-go"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	// "go.opentelemetry.io/otel/attribute"
	// "go.opentelemetry.io/otel/trace"
	"gopkg.in/gomail.v2"
)

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(999999)
	return fmt.Sprintf("%06d", otp)
}

func sendOTP(email, otp string) error {
	// Create an email message
	message := gomail.NewMessage()
	// message.SetHeader("From", "your-email@example.com") // Replace with your email
	message.SetHeader("From", "playlistzr@gmail.com") // Replace with your email
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Login OTP")
	message.SetBody("text/plain", "Your OTP is: "+otp)

	// log.Println("masuk-SEND")
	// // Setup email server configuration
	// smtpServer := "smtp.example.com" // Replace with your SMTP server
	// smtpPort := 587                  // Replace with your SMTP port
	// smtpUsername := "your-username"  // Replace with your SMTP username
	// smtpPassword := "your-password"  // Replace with your SMTP password

	// Setup email server configuration
	smtpServer := "smtp.gmail.com" // Replace with your SMTP server
	smtpPort := 587                // Replace with your SMTP port
	smtpUsername := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// Dial the SMTP server
	dialer := gomail.NewDialer(smtpServer, smtpPort, smtpUsername, smtpPassword)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		return err
	}

	return nil
}

func (s Service) GetGoldUser(ctx context.Context) ([]goldEntity.GetGoldUser, error) {
	log.Println("service GetGoldUser object")

	users, err := s.goldgym.GetGoldUser(ctx)
	log.Println("servicegolduser", users)
	if err != nil {
		return users, errors.Wrap(err, "[Service][GetGoldUser]")
	}
	return users, nil
}

// func (s Service) GetGoldUserByEmail(ctx context.Context, email string) (goldEntity.GetGoldUser, error) {
func (s Service) GetGoldUserByEmail(ctx context.Context, email string) (string, error) {
	var (
		result string
	)
	log.Println("service GetGoldUserByEmail object")

	userss, err := s.goldgym.GetGoldUserByEmail(ctx, email)

	if userss != (goldEntity.GetGoldUserss{}) {
		result = "TERDAFTAR"
	}

	if userss == (goldEntity.GetGoldUserss{}) {
		result = "TIDAK TERDAFTAR"
	}

	if userss.GoldValidasiYN == "N" {
		result = "BELUM TERVALIDASI"
	}

	log.Println("servicegolduserbyemail", result)
	if err != nil {
		return result, errors.Wrap(err, "[Service][GetGoldUserByEmail]")
	}
	return result, nil
}

// GetGoldUserDataByEmail returns full user entity by email (for gRPC)
func (s Service) GetGoldUserDataByEmail(ctx context.Context, email string) (goldEntity.GetGoldUserss, error) {
	log.Println("service GetGoldUserDataByEmail object")

	user, err := s.goldgym.GetGoldUserByEmail(ctx, email)
	if err != nil {
		return user, errors.Wrap(err, "[Service][GetGoldUserDataByEmail]")
	}

	return user, nil
}

func (s Service) InsertGoldUser(ctx context.Context, user goldEntity.GetGoldUsers) (interface{}, error) {
	var (
		err    error
		result string
		users  goldEntity.GetGoldUserss
	)
	log.Println("service user object", user)

	// code, _ := strconv.Atoi(jadwal.JadwalData.JwlCode)
	users, err = s.goldgym.GetGoldUserByEmail(ctx, user.GoldEmail)
	if err != nil {
		fmt.Println("test", err.Error())
		if err.Error() == "record not found" {
			// if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		} else {
			return result, errors.Wrap(err, "[SERVICE][InsertGoldUser][GetGoldUserByEmail]")
		}
	}
	log.Println("data", user)
	log.Println("users", users)

	if users == (goldEntity.GetGoldUserss{}) {
		// Hash Password
		// if users.GoldPassword != "" {
		// 	hashedPassword, err := argon2pw.GenerateSaltedHash(user.GoldPassword)
		// 	if err != nil {
		// 		return errors.Wrap(err, "[SERVICE][EditUser]")
		// 	}

		// 	user.GoldPassword = hashedPassword
		// }

		// err = s.core.ChangePassword(ctx, _user.NIP, _user.Password)
		// if err != nil {
		// 	return errors.Wrap(err, "[SERVICE][ResetPassword]")
		// }

		hashedPassword, err := hashPassword(user.GoldPassword)
		if err != nil {
			return result, errors.Wrap(err, "[SERVICE][CreateUser]")
		}

		// Kolom kartu (nomor kartu/CVV/nama pemegang) legacy & opsional --
		// registrasi modern tidak mengisinya. Skip hashing kalau kosong,
		// supaya tidak error "Password length cannot be 0" (argon2pw
		// menganggap string kosong sebagai input tidak valid).
		hashedNomorKartu := user.GoldNomorKartu
		if hashedNomorKartu != "" {
			hashedNomorKartu, err = hashPassword(hashedNomorKartu)
			if err != nil {
				return result, errors.Wrap(err, "[SERVICE][CreateUser]")
			}
		}
		hashedNomorCvv := user.GoldCvv
		if hashedNomorCvv != "" {
			hashedNomorCvv, err = hashPassword(hashedNomorCvv)
			if err != nil {
				return result, errors.Wrap(err, "[SERVICE][CreateUser]")
			}
		}
		hashedNamaPemegangKartu := user.GoldPemegangKartu
		if hashedNamaPemegangKartu != "" {
			hashedNamaPemegangKartu, err = hashPassword(hashedNamaPemegangKartu)
			if err != nil {
				return result, errors.Wrap(err, "[SERVICE][CreateUser]")
			}
		}

		user.GoldPassword = hashedPassword
		user.GoldNomorKartu = hashedNomorKartu
		user.GoldCvv = hashedNomorCvv
		user.GoldPemegangKartu = hashedNamaPemegangKartu

		log.Println("user-service", user)

		log.Println("GoldNomorKartu-length", len(user.GoldNomorKartu))
		log.Println("GoldCvv-length", len(user.GoldCvv))

		otp := generateOTP()
		// sendOTP(user.GoldEmail, otp)
		// -----------------------------------------------------------------------------------------------------------------

		message := gomail.NewMessage()
		// message.SetHeader("From", "your-email@example.com") // Replace with your email
		message.SetHeader("From", "playlistzr@gmail.com") // Replace with your email
		message.SetHeader("To", user.GoldEmail)
		message.SetHeader("Subject", "Login OTP")
		message.SetBody("text/plain", "Your OTP is: "+otp)

		// imageUrl := "https://media.tenor.com/qebfaxdCiSIAAAAd/spareaccountv2-usopp-spare-account-v2-gifs-discord-gif-one-piece.gif"

		// // Define the HTML body with the embedded image
		// htmlBody := `
		// <html>
		//     <body>
		//         <p>This is an email with an embedded image:</p>
		//         <img src="` + imageUrl + `" alt="Embedded GIF">
		//     </body>
		// </html>`
		// message.SetBody("text/html", "test: "+htmlBody)

		// // Setup email server configuration
		// smtpServer := "smtp.example.com" // Replace with your SMTP server
		// smtpPort := 587                  // Replace with your SMTP port
		// smtpUsername := "your-username"  // Replace with your SMTP username
		// smtpPassword := "your-password"  // Replace with your SMTP password

		// Setup email server configuration
		smtpServer := "smtp.gmail.com" // Replace with your SMTP server
		smtpPort := 587                // Replace with your SMTP port
		smtpUsername := os.Getenv("SMTP_EMAIL")
		smtpPassword := os.Getenv("SMTP_PASSWORD")

		// Dial the SMTP server
		dialer := gomail.NewDialer(smtpServer, smtpPort, smtpUsername, smtpPassword)

		// Send the email
		if err := dialer.DialAndSend(message); err != nil {
			return result, errors.Wrap(err, "[Service][GetGoldUserByEmailLogin]")
		}

		// -----------------------------------------------------------------------------------------------------------------
		user.GoldOTP = otp
		result, err = s.goldgym.InsertGoldUser(ctx, user)
		if err != nil {
			return result, errors.Wrap(err, "[SERVICE][InsertGoldUser]")
		}
		result = "Sukses"
	} else {
		result = "Gagal - Email Sudah Terdaftar"
	}

	return result, err
}

func GenerateNumber(length int) (string, error) {
	const otpChars = "1234567890"
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}
	return string(buffer), nil
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s Service) LoginUser(ctx context.Context, _user, _password string, _host, device string) (auth.Token, map[string]interface{}, string, error) {
	var (
		refreshToken string
		err          error
	)
	token := auth.Token{}
	metadata := make(map[string]interface{})

	user, err := s.goldgym.GetGoldUserByEmail(ctx, _user)
	if err != nil {
		return token, metadata, refreshToken, errors.Wrap(err, "[SERVICE][Login]")
	}

	password := user.GoldPassword

	valid, err := compareHash(password, _password)
	if err != nil {
		return token, metadata, refreshToken, errors.Wrap(err, "[SERVICE][Login]")
	}
	if !valid {
		return token, metadata, refreshToken, errors.Wrap(errors.New("invalid password"), "[SERVICE][Login]")
	}

	t := time.Now()
	d := 15 * time.Minute
	dr := 7 * 24 * time.Hour
	e := t.Add(d)
	er := t.Add(dr)

	sessionID := uuid.New().String()

	// role default SELLER untuk user lama yang kolom rolenya kosong
	role := user.GoldRole
	if role == "" {
		role = "SELLER"
	}

	sign := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  goldAuthEntity.JwtApplicationName,
		"sub":  user.GoldId,
		"jti":  sessionID,
		"nbf":  t.Unix(),
		"iat":  t.Unix(),
		"exp":  e.Unix(),
		"role": role,
	})

	accessToken, err := sign.SignedString(goldAuthEntity.JwtSecret)
	if err != nil {
		return token, metadata, refreshToken, errors.Wrap(err, "[SERVICE][Login]")
	}
	user.GoldLastLoginHost = _host
	err = s.goldgym.UpdateLastLogin(ctx, user)
	if err != nil {
		return token, metadata, refreshToken, errors.Wrap(err, "[SERVICE][Login]")
	}

	token = auth.Token{
		AccessToken:         accessToken,
		ExpiresIn:           e.Unix() - t.Unix(),
		ExpiresAt:           e.Unix(),
		TokenType:           "Bearer",
		ForceChangePassword: user.GoldForceChangePassword,
	}

	refreshToken, err = GenerateRefreshToken()
	if err != nil {
		return token, metadata, refreshToken, errors.Wrap(err, "[SERVICE][Login]")
	}
	tokenHash := HashToken(refreshToken)
	tokenStruct := goldTokenEntity.TokenRedis{
		UserID:    user.GoldId,
		SessionId: sessionID,
		ExpiresAt: er.Unix(),
		IP:        _host,
		Device:    device,
		Role:      role,
	}
	err = s.redis.AddToRedis(ctx, tokenStruct, "refresh:"+tokenHash, 7*24*time.Hour)
	if err != nil {
		return token, metadata, refreshToken, errors.Wrap(err, "[SERVICE][Login]")
	}
	metadata["username"] = user.GoldNama
	metadata["role"] = role
	metadata["gold_id"] = user.GoldId
	// flag "sudah mendaftar sebagai pembeli" — FE pakai untuk menampilkan
	// menu Mode Pembeli (role BUYER selalu Y)
	buyerYN := user.GoldBuyerYN
	if role == "BUYER" {
		buyerYN = "Y"
	}
	if buyerYN == "" {
		buyerYN = "N"
	}
	metadata["buyer_yn"] = buyerYN
	log.Println("metadata", metadata)
	return token, metadata, refreshToken, err
}

// RegisterBuyer mendaftarkan akun tanpa OTP — langsung aktif.
// Role yang diizinkan hanya BUYER (default) atau SELLER; ADMIN diset manual di DB.
// gold_toko diisi ketika pembeli didaftarkan lewat menu penjual (tampil di nota).
func (s Service) RegisterBuyer(ctx context.Context, req goldEntity.RegisterBuyerRequest) (goldEntity.BuyerRow, error) {
	var buyer goldEntity.BuyerRow

	// Batas pendaftaran soft-launch: maksimal MaxRegisteredUsers akun
	// terdaftar. Dicek lebih dulu (sebelum validasi lain) supaya tidak
	// membebani DB dengan lookup lain kalau memang sudah penuh.
	total, err := s.goldgym.CountGoldUsers(ctx)
	if err != nil {
		return buyer, errors.Wrap(err, "[SERVICE][RegisterBuyer][CountGoldUsers]")
	}
	if total >= goldEntity.MaxRegisteredUsers {
		return buyer, errors.New("pendaftaran sudah mencapai batas maksimal, coba lagi nanti")
	}

	req.GoldEmail = strings.TrimSpace(req.GoldEmail)
	req.GoldNama = strings.TrimSpace(req.GoldNama)
	req.GoldNomorHp = strings.TrimSpace(req.GoldNomorHp)
	if req.GoldEmail == "" || req.GoldPassword == "" || req.GoldNama == "" {
		return buyer, errors.New("nama, email, dan password wajib diisi")
	}
	if !strings.Contains(req.GoldEmail, "@") || strings.ContainsAny(req.GoldEmail, " \t") {
		return buyer, errors.New("format email tidak valid")
	}
	if len(req.GoldPassword) < 6 {
		return buyer, errors.New("password minimal 6 karakter")
	}
	if len(req.GoldNama) > 100 {
		return buyer, errors.New("nama maksimal 100 karakter")
	}
	if len(strings.TrimSpace(req.GoldToko)) > 100 {
		return buyer, errors.New("nama toko maksimal 100 karakter")
	}

	role := strings.ToUpper(strings.TrimSpace(req.GoldRole))
	if role == "" {
		role = "BUYER"
	}
	if role != "BUYER" && role != "SELLER" {
		return buyer, errors.New("role harus BUYER atau SELLER")
	}

	existing, err := s.goldgym.GetGoldUserByEmail(ctx, req.GoldEmail)
	if err != nil && err.Error() != "record not found" {
		return buyer, errors.Wrap(err, "[SERVICE][RegisterBuyer][GetGoldUserByEmail]")
	}
	if existing.GoldEmail != "" {
		return buyer, errors.New("email sudah terdaftar")
	}

	hashedPassword, err := hashPassword(req.GoldPassword)
	if err != nil {
		return buyer, errors.Wrap(err, "[SERVICE][RegisterBuyer][HashPassword]")
	}

	buyer = goldEntity.BuyerRow{
		GoldEmail:         req.GoldEmail,
		GoldPassword:      hashedPassword,
		GoldNama:          req.GoldNama,
		GoldNomorHp:       req.GoldNomorHp,
		GoldValidasiYN:    "Y",
		GoldUpdatedBy:     "SYSTEM",
		GoldLastLoginHost: "-",
		GoldRole:          role,
		GoldStatus:        1,
	}
	if toko := strings.TrimSpace(req.GoldToko); toko != "" {
		buyer.GoldToko = &toko
	}

	buyer, err = s.goldgym.InsertBuyer(ctx, buyer)
	if err != nil {
		return buyer, errors.Wrap(err, "[SERVICE][RegisterBuyer][InsertBuyer]")
	}

	buyer.GoldPassword = ""
	return buyer, nil
}

func (s Service) GetAllSubscription(ctx context.Context) ([]goldEntity.Subscription, error) {
	log.Println("service GetAllSubscription object")

	users, err := s.goldgym.GetAllSubscription(ctx)
	log.Println("serviceGetAllSubscription", users)
	if err != nil {
		return users, errors.Wrap(err, "[Service][GetAllSubscription]")
	}
	return users, nil
}

func (s Service) InsertSubscriptionUser(ctx context.Context, subs goldEntity.InsertSubsAll) (string, error) {
	var (
		result           string
		err              error
		detailData       goldEntity.SubscriptionDetail
		insertDetailData []goldEntity.SubscriptionDetail
		totalHarga       float64
	)

	header, err := s.goldgym.GetAllSubscription(ctx)
	user, err := s.goldgym.GetGoldUserByEmail(ctx, subs.HeaderData.GoldEmail)
	if err != nil {
		// result = "Detail - Gagal - Email Tidak Tersedia"
		return result, errors.Wrap(err, "[Service][GetGoldUserByEmail]")
	}
	if user == (goldEntity.GetGoldUserss{}) {
		result = "Detail - Gagal - Email Tidak Tersedia"
		return result, errors.Wrap(err, "[Service][GetAllSubscription]")
	}

	subs.HeaderData.GoldId = user.GoldId
	// log.Println("len-detail", len(subs.DetailData))

	if len(subs.DetailData) == 1 && subs.DetailData[0].GoldMenuId == 1 {
		subs.DetailData[0].GoldNamaPaket = header[0].GoldNamaPaket
		subs.DetailData[0].GoldNamaLayanan = header[0].GoldNamaLayanan
		subs.DetailData[0].GoldHarga = header[0].GoldHarga
		subs.DetailData[0].GoldId = user.GoldId
		subs.DetailData[0].GoldJadwal = header[0].GoldJadwal
		subs.DetailData[0].GoldListLatihan = header[0].GoldListLatihan
		subs.DetailData[0].GoldJumlahpertemuan = header[0].GoldJumlahpertemuan
		subs.DetailData[0].GoldDurasi = header[0].GoldDurasi
		subs.DetailData[0].GoldStatuslangganan = "Belum Berlangganan"
		err = s.goldgym.InsertSubscriptionDetail(ctx, subs.DetailData[0])
		if err != nil {
			result = "Detail - Gagal"
			return result, errors.Wrap(err, "[Service][InsertSubscriptionDetail]")
		}
		log.Println("masokDetail-1")
	}

	if len(subs.DetailData) == 1 && subs.DetailData[0].GoldMenuId > 1 {
		subs.DetailData[0].GoldNamaPaket = header[subs.DetailData[0].GoldMenuId-1].GoldNamaPaket
		subs.DetailData[0].GoldNamaLayanan = header[subs.DetailData[0].GoldMenuId-1].GoldNamaLayanan
		subs.DetailData[0].GoldHarga = header[subs.DetailData[0].GoldMenuId-1].GoldHarga
		subs.DetailData[0].GoldId = subs.HeaderData.GoldId
		subs.DetailData[0].GoldJadwal = header[subs.DetailData[0].GoldMenuId-1].GoldJadwal
		subs.DetailData[0].GoldListLatihan = header[subs.DetailData[0].GoldMenuId-1].GoldListLatihan
		subs.DetailData[0].GoldJumlahpertemuan = header[subs.DetailData[0].GoldMenuId-1].GoldJumlahpertemuan
		subs.DetailData[0].GoldDurasi = header[subs.DetailData[0].GoldMenuId-1].GoldDurasi
		subs.DetailData[0].GoldStatuslangganan = "Belum Berlangganan"
		err = s.goldgym.InsertSubscriptionDetail(ctx, subs.DetailData[0])
		if err != nil {
			result = "Detail - Gagal"
			return result, errors.Wrap(err, "[Service][InsertSubscriptionDetail]")
		}
		log.Println("masokDetail-2")
	}

	log.Println("testMASOK", len(subs.DetailData))

	if len(subs.DetailData) > 1 {
		for x := range subs.DetailData {
			detailData = goldEntity.SubscriptionDetail{
				GoldMenuId:          subs.DetailData[x].GoldMenuId,
				GoldNamaPaket:       header[x].GoldNamaPaket,
				GoldNamaLayanan:     header[x].GoldNamaLayanan,
				GoldHarga:           header[x].GoldHarga,
				GoldId:              subs.HeaderData.GoldId,
				GoldJadwal:          header[x].GoldJadwal,
				GoldListLatihan:     header[x].GoldListLatihan,
				GoldJumlahpertemuan: header[x].GoldJumlahpertemuan,
				GoldDurasi:          header[x].GoldDurasi,
				GoldStatuslangganan: "Belum Berlangganan",
				// GoldStatuslangganan: subs.DetailData[x].GoldStatuslangganan,
			}
			totalHarga += header[x].GoldHarga
			insertDetailData = append(insertDetailData, detailData)

			// subs.DetailData[x].GoldNamaPaket = header[x].GoldNamaPaket
			// subs.DetailData[x].GoldNamaLayanan = header[x].GoldNamaLayanan
			// subs.DetailData[x].GoldHarga = header[x].GoldHarga
			// subs.DetailData[x].GoldId = subs.HeaderData.GoldId
			// subs.DetailData[x].GoldJadwal = header[x].GoldJadwal
			// subs.DetailData[x].GoldListLatihan = header[x].GoldListLatihan
			// subs.DetailData[x].GoldJumlahpertemuan = header[x].GoldJumlahpertemuan
			// subs.DetailData[x].GoldDurasi = header[x].GoldDurasi
		}
		log.Println("insertDetailData", insertDetailData)
		limitzI := 50
		totalzI := len(insertDetailData)
		countzI := int(math.Ceil(float64(totalzI) / float64(limitzI)))
		for i := 0; i < countzI; i++ {
			startzI := limitzI * i
			endzI := limitzI * (i + 1)
			if endzI > totalzI {
				endzI = totalzI
			}
			tempUpdatez := insertDetailData[startzI:endzI]
			err = s.goldgym.BulkInsertSubscriptionDetail(ctx, tempUpdatez)
			if err != nil {
				log.Println(err, "[Service][BulkInsertSubscriptionDetail]")
				// return result, errors.Wrap(err, "[Service][UpdateDataProcodFromTempSelisih]")
			}
		}
		log.Println("masokDetail-3")
	}

	subs.HeaderData.GoldTotalharga = totalHarga

	err = s.goldgym.InsertSubscription(ctx, subs.HeaderData)
	if err != nil {
		result = "Header - Gagal"
		return result, errors.Wrap(err, "[Service][InsertSubscription]")
	}

	result = "Berhasil"
	return result, err
}

func (s Service) DeleteSubscriptionHeader(ctx context.Context, subs goldEntity.DeleteSubs) (string, error) {
	var (
		result string
		err    error
	)
	// err = s.goldgym.DeleteSubscriptionHeader(ctx, subs)
	// if err != nil {
	// 	result = "Header - Gagal"
	// 	return result, errors.Wrap(err, "[Service][InsertSubscription]")
	// }

	err = s.goldgym.DeleteSubscriptionDetail(ctx, subs)
	if err != nil {
		result = "Detail - Gagal"
		return result, errors.Wrap(err, "[Service][InsertSubscriptionDetail]")
	}

	result = "Berhasil"
	return result, err
}

func (s Service) UpdateSubscriptionDetail(ctx context.Context, subs goldEntity.UpdateSubs) (string, error) {
	var (
		result string
		err    error
	)
	err = s.goldgym.UpdateSubscriptionDetail(ctx, subs)
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][InsertSubscription]")
	}

	result = "Berhasil"
	return result, err
}

// func (s Service) UpdateValidation(ctx context.Context)

func (s Service) UpdateDataPeserta(ctx context.Context, subs goldEntity.UpdatePassword) (string, error) {
	var (
		result string
		err    error
	)

	header, err := s.goldgym.GetValidationGoldOTP(ctx, subs.GoldOTP)

	if header.GoldOTP == "" {
		result = "Please Validation OTP First"
		return result, errors.Wrap(err, "[Service][InsertSubscription]")
	}

	if subs.GoldEmail == "" && subs.GoldOTP == "" {
		result = "Please Field the Email and OTP"
		return result, errors.Wrap(err, "[Service][InsertSubscription]")
	}

	if subs.GoldEmail == "" {
		result = "Please Field the Email"
		return result, errors.Wrap(err, "[Service][InsertSubscription]")
	}

	if subs.GoldOTP == "" {
		result = "Please Field the OTP"
		return result, errors.Wrap(err, "[Service][InsertSubscription]")
	}

	if subs.GoldOTP != header.GoldOTP {
		result = "OTP is incorrect (validation otp)"
		return result, errors.Wrap(err, "[Service][InsertSubscription]")
	}

	// OTP berlaku 10 menit sejak dibuat. Tanpa batas ini, OTP 6 digit
	// (1 juta kemungkinan) tidak pernah kedaluwarsa dan rawan di-brute-force
	// (temuan security review 2026-08-18).
	const otpTTL = 10 * time.Minute
	if header.GoldOtpCreatedAt == nil || time.Since(*header.GoldOtpCreatedAt) > otpTTL {
		result = "OTP sudah kedaluwarsa, silakan minta OTP baru"
		return result, errors.New("otp expired")
	}

	if subs.GoldOTP == header.GoldOTP {
		err = s.goldgym.UpdateDataPeserta(ctx, subs)
		err = s.goldgym.UpdateOtpIsNull(ctx, subs.GoldEmail)
		if err != nil {
			result = "Gagal"
			return result, errors.Wrap(err, "[Service][InsertSubscription]")
		}
	}
	result = "Berhasil"
	return result, err
}

// RegisterAsBuyer menandai akun user login sudah mendaftar sebagai pembeli
// (konfirmasi dari menu Daftar Pembeli) — setelah ini menu Mode Pembeli
// muncul di aplikasi.
func (s Service) RegisterAsBuyer(ctx context.Context, goldid int) (string, error) {
	if goldid <= 0 {
		return "Gagal", errors.New("user tidak valid")
	}
	if err := s.goldgym.UpdateGoldBuyerYN(ctx, goldid); err != nil {
		return "Gagal", errors.Wrap(err, "[Service][RegisterAsBuyer]")
	}
	return "Berhasil", nil
}

// SetToko menyimpan nama toko milik user login. Dipakai fitur "penjual jadi
// pembeli": saat user belanja di outlet lain, nota menampilkan nama toko ini
// (via sale_cust_id → gold_toko).
func (s Service) SetToko(ctx context.Context, goldid int, toko string) (string, error) {
	toko = strings.TrimSpace(toko)
	if goldid <= 0 {
		return "Gagal", errors.New("user tidak valid")
	}
	if toko == "" {
		return "Gagal", errors.New("nama toko wajib diisi")
	}
	if len(toko) > 100 {
		return "Gagal", errors.New("nama toko maksimal 100 karakter")
	}
	if err := s.goldgym.UpdateGoldToko(ctx, goldid, toko); err != nil {
		return "Gagal", errors.Wrap(err, "[Service][SetToko]")
	}
	return "Berhasil", nil
}

func (s Service) UpdateNama(ctx context.Context, subs goldEntity.UpdateNama) (string, error) {
	var (
		result string
		err    error
	)
	err = s.goldgym.UpdateNama(ctx, subs)
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][InsertSubscription]")
	}

	result = "Berhasil"
	return result, err
}

// hashedNomorKartu, err := argon2pw.GenerateSaltedHash(user.GoldNomorKartu)
// 		if err != nil {
// 			return result, errors.Wrap(err, "[SERVICE][CreateUser]")
// 		}

func (s Service) UpdateKartu(ctx context.Context, subs goldEntity.UpdateKartu) (string, error) {
	var (
		result string
		err    error
	)
	hashedNomorKartu, err := hashPassword(subs.GoldNomorKartu)
	if err != nil {
		return result, errors.Wrap(err, "[SERVICE][CreateUser]")
	}
	hashedCvv, err := hashPassword(subs.GoldCvv)
	if err != nil {
		return result, errors.Wrap(err, "[SERVICE][CreateUser]")
	}
	subs.GoldNomorKartu = hashedNomorKartu
	subs.GoldCvv = hashedCvv
	err = s.goldgym.UpdateKartu(ctx, subs)
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][InsertSubscription]")
	}

	result = "Berhasil"
	return result, err
}

func (s Service) Logout(ctx context.Context, subs goldEntity.Logout) (string, error) {
	var (
		result string
		err    error
	)
	err = s.goldgym.Logout(ctx, subs)
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][InsertSubscription]")
	}

	result = "Berhasil"
	return result, err
}

func (s Service) GetSubsWithUser(ctx context.Context) ([]goldEntity.GetSubsWithUser, error) {
	log.Println("service GetGoldUser object")

	users, err := s.goldgym.GetSubsWithUser(ctx)
	log.Println("servicegolduser", users)
	if err != nil {
		return users, errors.Wrap(err, "[Service][GetGoldUser]")
	}
	return users, nil
}

func (s Service) UpdateValidationOTP(ctx context.Context, otp string, email string) (string, error) {
	var (
		result string
		err    error
	)

	log.Println("params-service", otp, email)

	header, err := s.goldgym.GetValidationGoldOTP(ctx, otp)

	if header.GoldOTP == otp {
		err = s.goldgym.UpdateValidationOTP(ctx, email)
		err = s.goldgym.UpdateOtpIsNull(ctx, email)
		// if err != nil {
		result = "Berhasil"
		// 	return result, errors.Wrap(err, "[Service][InsertSubscription]")
		// }
	}

	if header.GoldOTP != otp {
		// err = s.goldgym.UpdateValidationOTP(ctx, email)
		// if err != nil {
		result = "OTP is incorrect"
		return result, errors.Wrap(err, "[Service][InsertSubscription]")
		// }
	}

	return result, err
}

func (s Service) UpdateOTP(ctx context.Context, email string) (string, error) {
	var (
		result string
		err    error
		otp    string
	)

	otp = generateOTP()

	err = sendOTP(email, otp)
	if err != nil {
		result = "Error"
		return result, errors.Wrap(err, "[Service][sendOTP]")
	}
	log.Println("params-service", otp, email)

	err = s.goldgym.UpdateOTP(ctx, otp, email)

	if err != nil {
		result = "OTP is incorrect"
		return result, errors.Wrap(err, "[Service][InsertSubscription]")
	}
	result = "Berhasil"

	// header, err := s.goldgym.GetValidationGoldOTP(ctx, otp)

	// if header.GoldOTP == otp {
	// 	err = s.goldgym.UpdateValidationOTP(ctx, email)
	// 	err = s.goldgym.UpdateOtpIsNull(ctx, email)
	// 	// if err != nil {
	// 	result = "Berhasil"
	// 	// 	return result, errors.Wrap(err, "[Service][InsertSubscription]")
	// 	// }
	// }

	// if header.GoldOTP != otp {
	// 	// err = s.goldgym.UpdateValidationOTP(ctx, email)
	// 	// if err != nil {
	// 	result = "OTP is incorrect"
	// 	return result, errors.Wrap(err, "[Service][InsertSubscription]")
	// 	// }
	// }

	return result, err
}

func (s Service) PaymentValidation(ctx context.Context, id int, menuid int, email string) (string, error) {
	var (
		result string
		err    error
	)

	otp := generateOTP()

	err = sendOTP(email, otp)
	if err != nil {
		result = "Error"
		return result, errors.Wrap(err, "[Service][sendOTP]")
	}
	log.Println("params-service", otp, email)

	// err = s.goldgym.UpdateOTPSubscription(ctx, otp, email)

	// if err != nil {
	// 	result = "OTP is incorrect"
	// 	return result, errors.Wrap(err, "[Service][InsertSubscription]")
	// }
	result = "Berhasil"

	return result, err
}

func (s Service) InsertSubscriptionDetail(ctx context.Context, user goldEntity.SubscriptionDetail) (string, error, response.Response) {
	var (
		result string
		err    error
		resp   response.Response
	)

	header, err := s.goldgym.GetOneSubscription(ctx, user.GoldMenuId)
	log.Println("testUser", user)
	headers, err := s.goldgym.GetSubscriptionHeader(ctx, user.GoldId)
	log.Println("tesHeaders", headers)
	if headers == (goldEntity.SubscriptionHeader{}) {
		result = "Subscription Header Empty"
		resp.StatusCode = 501
		resp.Error.Status = true
		return result, errors.Wrap(err, "[Service][InsertSubscription]"), resp
	}

	user.GoldNamaPaket = header.GoldNamaPaket
	user.GoldNamaLayanan = header.GoldNamaLayanan
	user.GoldHarga = header.GoldHarga
	user.GoldJadwal = header.GoldJadwal
	user.GoldListLatihan = header.GoldListLatihan
	user.GoldJumlahpertemuan = header.GoldJumlahpertemuan
	user.GoldDurasi = header.GoldDurasi
	user.GoldStatuslangganan = "Belum Berlangganan"

	err = s.goldgym.InsertSubscriptionDetail(ctx, user)
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][InsertSubscription]"), resp
	}

	result = "Berhasil"
	return result, err, resp
}

// func (s Service) UpdateOTPSubscription(ctx context.Context, id string) (string, time.Time, error) {
func (s Service) UpdateOTPSubscription(ctx context.Context, id string) (string, error) {
	var (
		result string
		err    error
		otp    string
		// expiration time.Time
		ids int
		// idss       int
	)
	// err = errors.New("404 Not Found")
	// err.response.Error
	log.Println("id", id)
	header, err := s.goldgym.GetGoldUserByEmail(ctx, id)

	ids = header.GoldId

	log.Println("header.GoldId", header.GoldId)
	log.Println("ids", ids)

	otp = generateOTP()

	err = sendOTP(id, otp)
	if err != nil {
		result = "Error"
		// return result, expiration, errors.Wrap(err, "[Service][sendOTP]")
		return result, errors.Wrap(err, "[Service][sendOTP]")
	}
	log.Println("params-service", otp, id)
	// idss = strconv.Itoa(ids)
	err = s.goldgym.UpdateOTPSubscription(ctx, otp, ids)
	// idss, _ = strconv.Atoi(ids)

	// // Calculate the expiration time (e.g., 5 minutes from now)
	// expiration = time.Now().Add(5 * time.Minute)

	if err != nil {
		result = "OTP is incorrect"
		// return result, expiration, errors.Wrap(err, "[Service][UpdateOTPSubscription]")
		return result, errors.Wrap(err, "[Service][UpdateOTPSubscription]")
	}
	result = "Berhasil"

	// return result, expiration, err
	return result, err
}

func (s Service) UpdatePayment(ctx context.Context, otp string, email string) (string, error, response.Response) {
	var (
		result            string
		otpHourMinuteConv float64
		nowHourMinuteConv float64
		convMinute        string
		convMinuteDate    string
		updatePayment     goldEntity.UpdatePayment
		resp              response.Response
	)
	header, err := s.goldgym.GetGoldUserByEmail(ctx, email)
	if header == (goldEntity.GetGoldUserss{}) {
		result = "Email Not Available"
		resp.StatusCode = 501
		resp.Error.Status = true
		// return result, expiration, errors.Wrap(err, "[Service][sendOTP]")
		return result, errors.Wrap(err, "[Service][GetGoldUserByEmail]"), resp
	}
	log.Println("errsssssssssssss", err)
	if err != nil {
		result = "Error"
		resp.StatusCode = 501
		resp.Error.Status = true
		// return result, expiration, errors.Wrap(err, "[Service][sendOTP]")
		return result, errors.Wrap(err, "[Service][GetGoldUserByEmail]"), resp
	}
	log.Println("header", header)

	subs, err := s.goldgym.GetSubscriptionHeader(ctx, header.GoldId)

	if subs.GoldOTP.IsZero() {
		result = "Please do OTP Subscription First"
		resp.StatusCode = 501
		resp.Error.Status = true
		// err != nil
		// return result, expiration, errors.Wrap(err, "[Service][sendOTP]")
		return result, errors.Wrap(err, "[Service][OTP-Subscription]"), resp
	}

	if otp != subs.GoldOTP.String {
		result = "OTP Incorrect"
		resp.StatusCode = 501
		resp.Error.Status = true
		// return result, expiration, errors.Wrap(err, "[Service][sendOTP]")
		return result, errors.Wrap(err, "[Service][OTP]"), resp
	}

	if otp == subs.GoldOTP.String {

		log.Println("subs", subs)

		log.Println("testNow", time.Now())

		stringToDate, err := time.Parse("2006-01-02 15:04:05", subs.GoldLastupdate.String)
		log.Println("stringToDate", stringToDate)
		if stringToDate.Minute() == 1 || stringToDate.Minute() == 2 || stringToDate.Minute() == 3 || stringToDate.Minute() == 4 || stringToDate.Minute() == 5 || stringToDate.Minute() == 6 || stringToDate.Minute() == 7 || stringToDate.Minute() == 8 || stringToDate.Minute() == 9 {
			convMinuteDate = "0" + strconv.Itoa(stringToDate.Minute())
			log.Println("masoooook1")
		} else {
			convMinuteDate = strconv.Itoa(stringToDate.Minute())
			log.Println("masoooook2")
		}

		otpHourMinute := strconv.Itoa(stringToDate.Hour()) + convMinuteDate
		otpHourMinuteConv, _ = strconv.ParseFloat(otpHourMinute, 64)

		now := time.Now()
		log.Println("now", now)
		if now.Minute() == 1 || now.Minute() == 2 || now.Minute() == 3 || now.Minute() == 4 || now.Minute() == 5 || now.Minute() == 6 || now.Minute() == 7 || now.Minute() == 8 || now.Minute() == 9 {
			convMinute = "0" + strconv.Itoa(now.Minute())
			log.Println("masoooook1-Now")
		} else {
			convMinute = strconv.Itoa(now.Minute())
			log.Println("masoooook2-Now")
		}

		log.Println("testNow", now.Minute())
		nowHourMinute := strconv.Itoa(now.Hour()) + convMinute
		log.Println("nowHourMinute", nowHourMinute)
		nowHourMinuteConv, _ = strconv.ParseFloat(nowHourMinute, 64)

		log.Println("testHourMinute", nowHourMinuteConv)

		log.Println("otpHourMinuteConv-before", otpHourMinuteConv)

		if stringToDate.Minute() == 59 {
			conv := strconv.Itoa(stringToDate.Hour()+1) + "01"
			otpHourMinuteConv, _ = strconv.ParseFloat(conv, 64)
			nowHourMinuteConv += 1
		}

		log.Println("otpHourMinuteConv", otpHourMinuteConv)
		log.Println("nowHourMinuteConv", nowHourMinuteConv)

		if nowHourMinuteConv >= otpHourMinuteConv+5.0 {
			log.Println("true-Time")
			result = "OTP expired"
			// return result, expiration, errors.Wrap(err, "[Service][UpdateOTPSubscription]")
			resp.StatusCode = 501
			resp.Error.Status = true
			return result, errors.Wrap(err, "[Service][UpdatePayment]"), resp
		}

		if nowHourMinuteConv <= otpHourMinuteConv+5.0 {
			log.Println("false-Time")
			updatePayment.GoldID = header.GoldId
			err = s.goldgym.UpdateValidasiPaymentHeader(ctx, updatePayment)
			err = s.goldgym.UpdateValidasiPaymentDetail(ctx, updatePayment)
			result = "OTP true"
		}

	}

	return result, err, resp
}

func (s Service) GetSubscriptionHeaderTotalHarga(ctx context.Context, email string) (goldEntity.SubscriptionHeaderPayment, error) {
	var (
		users goldEntity.SubscriptionHeaderPayment
	)
	log.Println("service GetGoldUser object")
	header, err := s.goldgym.GetGoldUserByEmail(ctx, email)
	if err != nil {
		return users, errors.Wrap(err, "[Service][GetGoldUser]")
	}
	users, err = s.goldgym.GetSubscriptionHeaderTotalHarga(ctx, header.GoldId)
	log.Println("servicegolduser", users)
	if err != nil {
		return users, errors.Wrap(err, "[Service][GetGoldUser]")
	}
	return users, nil
}

func (s Service) UploadTestingImages(ctx context.Context, testing goldEntity.Testings) (string, error) {
	result, err := s.goldgym.UploadTestingImages(ctx, testing)
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][UploadTestingImages]")
	}
	result = "Sukses"

	return result, err
}

func (s Service) GetTestingImage(ctx context.Context, id int) ([]byte, error) {
	image, err := s.goldgym.GetTestingImages(ctx, id)
	if err != nil {
		return image, errors.Wrap(err, "[SERVICE][GetTestingImage]")
	}

	return image, err
}
