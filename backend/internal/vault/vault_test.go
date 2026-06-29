// SPDX-License-Identifier: MIT
package vault

import (
	"testing"
)

func TestGenerateKey(t *testing.T) {
	k1, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}
	k2, _ := GenerateKey()
	if k1 == k2 {
		t.Error("expected unique keys")
	}
}

func TestNew_ValidKey(t *testing.T) {
	key, _ := GenerateKey()
	v, err := New(key)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if v.IsZero() {
		t.Error("vault should not be zero")
	}
}

func TestNew_InvalidKey(t *testing.T) {
	_, err := New("nothex")
	if err == nil {
		t.Error("expected error for invalid hex")
	}
}

func TestNew_WrongKeyLength(t *testing.T) {
	_, err := New("aa") // 1 byte
	if err == nil {
		t.Error("expected error for short key")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key, _ := GenerateKey()
	v, _ := New(key)

	original := "my-secret-password-123!@#"
	ciphertext, err := v.Encrypt(original)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if ciphertext == original {
		t.Error("ciphertext should differ from plaintext")
	}

	decrypted, err := v.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decrypted != original {
		t.Errorf("expected %q, got %q", original, decrypted)
	}
}

func TestEncrypt_ProducesDifferentOutput(t *testing.T) {
	key, _ := GenerateKey()
	v, _ := New(key)

	c1, _ := v.Encrypt("same")
	c2, _ := v.Encrypt("same")
	if c1 == c2 {
		t.Error("encrypt should produce different output each time (nonce)")
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	k1, _ := GenerateKey()
	k2, _ := GenerateKey()
	v1, _ := New(k1)
	v2, _ := New(k2)

	ciphertext, _ := v1.Encrypt("secret")
	_, err := v2.Decrypt(ciphertext)
	if err == nil {
		t.Error("expected error decrypting with wrong key")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	key, _ := GenerateKey()
	v, _ := New(key)

	_, err := v.Decrypt("!!!invalid-base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64")
	}
}

func TestVault_IsZero(t *testing.T) {
	v := &Vault{key: make([]byte, 32)}
	if !v.IsZero() {
		t.Error("expected zero vault")
	}
}

func TestVault_NonZero(t *testing.T) {
	key, _ := GenerateKey()
	v, _ := New(key)
	if v.IsZero() {
		t.Error("expected non-zero vault")
	}
}
