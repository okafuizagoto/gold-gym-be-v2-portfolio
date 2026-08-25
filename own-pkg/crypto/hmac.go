package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

// HMACSHA256 generates HMAC-SHA256 signature
// Returns hex-encoded hash (lowercase)
func HMACSHA256(payload, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

// HMACSHA256Base64 generates HMAC-SHA256 signature
// Returns base64-encoded hash
func HMACSHA256Base64(payload, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// HMACSHA512 generates HMAC-SHA512 signature
// Returns hex-encoded hash (lowercase)
func HMACSHA512(payload, secret string) string {
	h := hmac.New(sha512.New, []byte(secret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

// HMACSHA512Base64 generates HMAC-SHA512 signature
// Returns base64-encoded hash
func HMACSHA512Base64(payload, secret string) string {
	h := hmac.New(sha512.New, []byte(secret))
	h.Write([]byte(payload))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// SHA256Hash generates SHA256 hash
// Returns hex-encoded hash (lowercase)
func SHA256Hash(data string) string {
	h := sha256.Sum256([]byte(data))
	return strings.ToLower(hex.EncodeToString(h[:]))
}

// VerifyHMACSHA512 verifies HMAC-SHA512 signature
// signatureBase64: base64-encoded signature to verify
func VerifyHMACSHA512(payload, secret, signatureBase64 string) bool {
	expected := HMACSHA512Base64(payload, secret)

	// Decode both signatures for comparison
	expectedBytes, err1 := base64.StdEncoding.DecodeString(expected)
	signatureBytes, err2 := base64.StdEncoding.DecodeString(signatureBase64)

	if err1 != nil || err2 != nil {
		return false
	}

	return hmac.Equal(expectedBytes, signatureBytes)
}

// VerifyHMACSHA256 verifies HMAC-SHA256 signature (hex-encoded)
func VerifyHMACSHA256(payload, secret, signatureHex string) bool {
	expected := HMACSHA256(payload, secret)
	return hmac.Equal([]byte(expected), []byte(strings.ToLower(signatureHex)))
}

// BuildServiceSignaturePayload builds the payload string for service signature
// Format: {METHOD}:{PATH}:{ACCESS_TOKEN}:{SHA256(BODY)}:{TIMESTAMP}
func BuildServiceSignaturePayload(method, path, accessToken, body, timestamp string) string {
	bodyHash := SHA256Hash(body)
	return strings.Join([]string{method, path, accessToken, bodyHash, timestamp}, ":")
}
