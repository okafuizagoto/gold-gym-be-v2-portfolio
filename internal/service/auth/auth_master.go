package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	goldAuthEntity "gold-gym-be/internal/entity/auth/v2"
	goldTokenEntity "gold-gym-be/internal/entity/token"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s Service) RefreshToken(ctx context.Context, refreshToken string) (goldAuthEntity.Token, error) {
	var (
		tokenRedis goldTokenEntity.TokenRedis
		token      goldAuthEntity.Token
	)
	// ----- hash refresh token ------
	hash := sha256.Sum256([]byte(refreshToken))
	refreshHash := hex.EncodeToString(hash[:])
	// ----- hash refresh token ------

	// ----- key refresh ------
	key := "refresh:" + refreshHash
	// ----- key refresh ------

	err := s.redis.GetFromRedis(ctx, key, &tokenRedis)
	if err != nil {
		// return token, errors.Wrap(err, "[SERVICE][RefreshToken][GetFromRedis]")
		return token, errors.New("invalid refresh token")
	}

	if tokenRedis.ExpiresAt < time.Now().Unix() {
		return token, errors.New("refresh token expired")
	}

	t := time.Now()
	d := 15 * time.Minute
	e := t.Add(d)

	role := tokenRedis.Role
	if role == "" {
		role = "SELLER"
	}

	sign := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  goldAuthEntity.JwtApplicationName,
		"sub":  tokenRedis.UserID,
		"jti":  tokenRedis.SessionId,
		"nbf":  t.Unix(),
		"iat":  t.Unix(),
		"exp":  e.Unix(),
		"role": role,
	})

	newAccessToken, err := sign.SignedString(goldAuthEntity.JwtSecret)

	token = goldAuthEntity.Token{
		AccessToken: newAccessToken,
		ExpiresIn:   e.Unix() - t.Unix(),
		ExpiresAt:   e.Unix(),
		TokenType:   "Bearer",
	}

	return token, err
}

func (s Service) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := HashToken(refreshToken)
	key := "refresh:" + tokenHash

	err := s.redis.DeleteFromRedis(ctx, key)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][Logout]")
	}

	return nil
}
