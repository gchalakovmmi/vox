package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/url"
	"os"
	"strings"
)

// encryptEmailToken encrypts a token with AES-CFB and returns a URL-safe string.
func encryptEmailToken(token string) (string, error) {
	// TODO: pass the EMAIL_LOGIN_ENCRYPT_KEY directly from the config
	key := []byte(os.Getenv("EMAIL_LOGIN_ENCRYPT_KEY"))
	if len(key) != 32 {
		return "", errors.New("EMAIL_LOGIN_ENCRYPT_KEY must be 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(token))
	cfb.XORKeyStream(cipherText, []byte(token))

	return url.QueryEscape(base64.StdEncoding.EncodeToString(append(iv, cipherText...))), nil
}

// decryptEmailToken reverses encryptEmailToken.
func decryptEmailToken(enc string) (string, error) {
	raw, err := url.QueryUnescape(enc)
	if err != nil {
		return "", err
	}
	raw = strings.ReplaceAll(raw, " ", "+")

	cipherBytes, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return "", err
	}
	if len(cipherBytes) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	// TODO: pass the EMAIL_LOGIN_ENCRYPT_KEY directly from the config
	key := []byte(os.Getenv("EMAIL_LOGIN_ENCRYPT_KEY"))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv, cipherText := cipherBytes[:aes.BlockSize], cipherBytes[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	plain := make([]byte, len(cipherText))
	cfb.XORKeyStream(plain, cipherText)
	return string(plain), nil
}
