package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type ClaimKey string

const (
	JwtApplicationName = "GOLD-GYM-BE"
)

var (
	JwtSigningMethod = jwt.SigningMethodHS256
	JwtSecret        = []byte("a7fecfed-14c8-4f54-84a7-e43fe9cf1823")
)

type Token struct {
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
	ExpiresIn           int64  `json:"expires_in"`
	ExpiresAt           int64  `json:"expires_at"`
	TokenType           string `json:"token_type"`
	ForceChangePassword int    `json:"force_change_password"`
}

type LoginRequest struct {
	NIP      string `json:"nip"`
	Password string `json:"password"`
}

type AccessListRequest struct {
	PTID  int    `json:"pt_id"`
	AppID int    `json:"app_id"`
	NIP   string `json:"nip"`
}
