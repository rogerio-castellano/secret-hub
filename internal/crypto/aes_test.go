package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := make([]byte, 32) // AES-256 requires a 32-byte key
	for i := range key {
		key[i] = byte(i)
	}
	plaintext := []byte("this is a test secret")

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("decrypted text does not match original. Got %s, want %s", decrypted, plaintext)
	}
}

func TestEncryptWithWrongKeySize(t *testing.T) {
	badKey := []byte("shortkey")
	if _, err := Encrypt(badKey, []byte("test")); err == nil {
		t.Error("expected error due to invalid key size, got nil")
	}
}
