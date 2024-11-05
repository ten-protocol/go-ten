package encryption

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestNewEncryptor(t *testing.T) {
	key := make([]byte, 32) // 256-bit key
	_, err := rand.Read(key)
	if err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}

	encryptor, err := NewEncryptor(key)
	if err != nil {
		t.Fatalf("NewEncryptor failed: %v", err)
	}

	if encryptor == nil {
		t.Fatal("NewEncryptor returned nil")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := make([]byte, 32) // 256-bit key
	_, err := rand.Read(key)
	if err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}

	encryptor, err := NewEncryptor(key)
	if err != nil {
		t.Fatalf("NewEncryptor failed: %v", err)
	}

	plaintext := []byte("Hello, World!")

	ciphertext, err := encryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := encryptor.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatalf("Decrypted text does not match original plaintext")
	}
}

func TestEncryptDecryptEmptyString(t *testing.T) {
	key := make([]byte, 32) // 256-bit key
	_, err := rand.Read(key)
	if err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}

	encryptor, err := NewEncryptor(key)
	if err != nil {
		t.Fatalf("NewEncryptor failed: %v", err)
	}

	plaintext := []byte("")

	ciphertext, err := encryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encryption of empty string failed: %v", err)
	}

	decrypted, err := encryptor.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decryption of empty string failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatalf("Decrypted empty string does not match original")
	}
}

func TestDecryptInvalidCiphertext(t *testing.T) {
	key := make([]byte, 32) // 256-bit key
	_, err := rand.Read(key)
	if err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}

	encryptor, err := NewEncryptor(key)
	if err != nil {
		t.Fatalf("NewEncryptor failed: %v", err)
	}

	invalidCiphertext := []byte("This is not a valid ciphertext")

	_, err = encryptor.Decrypt(invalidCiphertext)
	if err == nil {
		t.Fatal("Decryption of invalid ciphertext should have failed, but didn't")
	}
}

func TestNewEncryptorInvalidKeySize(t *testing.T) {
	invalidKey := make([]byte, 31) // Invalid key size (not 16, 24, or 32 bytes)
	_, err := rand.Read(invalidKey)
	if err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}

	_, err = NewEncryptor(invalidKey)
	if err == nil {
		t.Fatal("NewEncryptor should have failed with invalid key size, but didn't")
	}
}
