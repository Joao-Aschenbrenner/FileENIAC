package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

type Vault struct {
	key []byte
}

func New(keyHex string) (*Vault, error) {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid vault key hex: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("vault key must be 32 bytes (got %d)", len(key))
	}
	return &Vault{key: key}, nil
}

func GenerateKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("failed to generate vault key: %w", err)
	}
	return hex.EncodeToString(key), nil
}

func (v *Vault) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(v.key)
	if err != nil {
		return "", fmt.Errorf("vault encrypt: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("vault gcm: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("vault nonce: %w", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (v *Vault) Decrypt(cipherB64 string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cipherB64)
	if err != nil {
		return "", fmt.Errorf("vault decode: %w", err)
	}
	block, err := aes.NewCipher(v.key)
	if err != nil {
		return "", fmt.Errorf("vault decrypt cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("vault decrypt gcm: %w", err)
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("vault ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("vault decrypt failed: %w", err)
	}
	return string(plaintext), nil
}

func (v *Vault) IsZero() bool {
	for _, b := range v.key {
		if b != 0 {
			return false
		}
	}
	return true
}
