package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"
)

// RSAVerifySignature verifies an RSA SHA256 signature
// publicKeyBase64: base64 encoded public key (can be PEM or DER format)
// payload: the string that was signed
// signatureBase64: base64 encoded signature
func RSAVerifySignature(publicKeyBase64, payload, signatureBase64 string) (bool, error) {
	// Decode the public key from base64
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return false, errors.New("failed to decode public key from base64")
	}

	// Try to parse as PEM first
	var pubKey *rsa.PublicKey
	block, _ := pem.Decode(publicKeyBytes)
	if block != nil {
		// It's PEM format
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			// Try parsing as PKCS1
			pubKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
			if err != nil {
				return false, errors.New("failed to parse public key from PEM")
			}
		} else {
			var ok bool
			pubKey, ok = pub.(*rsa.PublicKey)
			if !ok {
				return false, errors.New("not an RSA public key")
			}
		}
	} else {
		// Try as raw DER format
		pub, err := x509.ParsePKIXPublicKey(publicKeyBytes)
		if err != nil {
			// Try parsing as PKCS1
			pubKey, err = x509.ParsePKCS1PublicKey(publicKeyBytes)
			if err != nil {
				return false, errors.New("failed to parse public key")
			}
		} else {
			var ok bool
			pubKey, ok = pub.(*rsa.PublicKey)
			if !ok {
				return false, errors.New("not an RSA public key")
			}
		}
	}

	// Decode the signature from base64
	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return false, errors.New("failed to decode signature from base64")
	}

	// Hash the payload
	hashed := sha256.Sum256([]byte(payload))

	// Verify the signature
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return false, nil // Invalid signature, but not an error
	}

	return true, nil
}

// RSASign signs a payload with RSA SHA256
// privateKeyInput: can be:
//   - base64 encoded PEM private key
//   - raw PEM private key string
//   - base64 encoded DER private key
//
// payload: the string to sign
// Returns: base64 encoded signature
func RSASign(privateKeyInput, payload string) (string, error) {
	var privateKeyBytes []byte
	var err error

	// Clean up the input - remove any whitespace/newlines that might break base64
	cleanInput := strings.TrimSpace(privateKeyInput)

	// Check if it's already a PEM string (not base64 encoded)
	if strings.Contains(cleanInput, "-----BEGIN") {
		privateKeyBytes = []byte(cleanInput)
	} else {
		// It's base64 encoded - decode it first
		// Try standard base64 first
		privateKeyBytes, err = base64.StdEncoding.DecodeString(cleanInput)
		if err != nil {
			// Try URL-safe base64
			urlSafe := strings.ReplaceAll(cleanInput, "-", "+")
			urlSafe = strings.ReplaceAll(urlSafe, "_", "/")

			// Add padding if needed
			switch len(urlSafe) % 4 {
			case 2:
				urlSafe += "=="
			case 3:
				urlSafe += "="
			}

			privateKeyBytes, err = base64.StdEncoding.DecodeString(urlSafe)
			if err != nil {
				return "", errors.New("failed to decode private key from base64: " + err.Error())
			}
		}
	}

	// Parse the private key
	privKey, err := parsePrivateKey(privateKeyBytes)
	if err != nil {
		return "", err
	}

	// Hash the payload
	hashed := sha256.Sum256([]byte(payload))

	// Sign
	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// parsePrivateKey tries to parse a private key from various formats
func parsePrivateKey(keyBytes []byte) (*rsa.PrivateKey, error) {
	// Try to parse as PEM first
	block, _ := pem.Decode(keyBytes)
	if block != nil {
		// It's PEM format
		return parsePrivateKeyDER(block.Bytes)
	}

	// Check if the decoded bytes might be a PEM string
	if strings.Contains(string(keyBytes), "-----BEGIN") {
		block, _ := pem.Decode(keyBytes)
		if block != nil {
			return parsePrivateKeyDER(block.Bytes)
		}
	}

	// Try as raw DER format
	return parsePrivateKeyDER(keyBytes)
}

// parsePrivateKeyDER parses DER-encoded private key
func parsePrivateKeyDER(derBytes []byte) (*rsa.PrivateKey, error) {
	// Try PKCS8 first
	key, err := x509.ParsePKCS8PrivateKey(derBytes)
	if err == nil {
		privKey, ok := key.(*rsa.PrivateKey)
		if ok {
			return privKey, nil
		}
		return nil, errors.New("not an RSA private key (PKCS8)")
	}

	// Try PKCS1
	privKey, err := x509.ParsePKCS1PrivateKey(derBytes)
	if err == nil {
		return privKey, nil
	}

	// Try EC and convert error message
	_, ecErr := x509.ParseECPrivateKey(derBytes)
	if ecErr == nil {
		return nil, errors.New("this is an EC private key, not RSA")
	}

	return nil, errors.New("failed to parse private key: not PKCS8 or PKCS1 format")
}

// ParsePublicKeyFromPEM parses a PEM-encoded public key string
func ParsePublicKeyFromPEM(pemStr string) (*rsa.PublicKey, error) {
	// Handle both with and without PEM headers
	if !strings.Contains(pemStr, "-----BEGIN") {
		// Try decoding from base64 first
		decoded, err := base64.StdEncoding.DecodeString(pemStr)
		if err == nil {
			pemStr = string(decoded)
		}
	}

	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		// Try parsing directly
		decoded, err := base64.StdEncoding.DecodeString(pemStr)
		if err != nil {
			return nil, errors.New("failed to decode public key")
		}

		pub, err := x509.ParsePKIXPublicKey(decoded)
		if err != nil {
			pubKey, err := x509.ParsePKCS1PublicKey(decoded)
			if err != nil {
				return nil, errors.New("failed to parse public key")
			}
			return pubKey, nil
		}

		pubKey, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("not an RSA public key")
		}
		return pubKey, nil
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, errors.New("failed to parse public key from PEM block")
		}
		return pubKey, nil
	}

	pubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return pubKey, nil
}
