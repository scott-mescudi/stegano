package pkg

import (
	"bytes"
	"fmt"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		password  string
		plaintext []byte
	}{
		{
			password:  "correcthorsebatterystaple",
			plaintext: []byte("This is a test message."),
		},
		{
			password:  "password123",
			plaintext: []byte("Another test message!"),
		},
		{
			password:  "password123",
			plaintext: []byte("Short"),
		},
		{
			password:  "password123",
			plaintext: []byte(""),
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("password=%s", tt.password), func(t *testing.T) {
			// Encrypt
			ciphertext, err := Encrypt(tt.password, tt.plaintext)
			if err != nil {
				t.Fatalf("Encryption failed: %v", err)
			}

			// Decrypt
			decryptedText, err := Decrypt(tt.password, ciphertext)
			if err != nil {
				t.Fatalf("Decryption failed: %v", err)
			}

			// Check if decrypted text matches original plaintext
			if !bytes.Equal(decryptedText, tt.plaintext) {
				t.Errorf("Decrypted text does not match original plaintext. Expected: %s, got: %s", tt.plaintext, decryptedText)
			}
		})
	}
}

func TestDecryptWithIncorrectPassword(t *testing.T) {
	plaintext := []byte("Test for incorrect password")
	password := "correctpassword"
	incorrectPassword := "wrongpassword"

	// Encrypt with the correct password
	ciphertext, err := Encrypt(password, plaintext)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Decrypt with incorrect password
	_, err = Decrypt(incorrectPassword, ciphertext)
	if err == nil {
		t.Errorf("Expected decryption to fail with incorrect password")
	}
}

func TestDecryptWithCorruptedCiphertext(t *testing.T) {
	plaintext := []byte("Test for corrupted ciphertext")
	password := "correctpassword"

	// Encrypt with the correct password
	ciphertext, err := Encrypt(password, plaintext)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Corrupt the ciphertext by changing some bytes
	corruptedCiphertext := append([]byte{}, ciphertext...)
	corruptedCiphertext[10] = 0xFF

	// Decrypt the corrupted ciphertext
	_, err = Decrypt(password, corruptedCiphertext)
	if err == nil {
		t.Errorf("Expected decryption to fail with corrupted ciphertext")
	}
}

func TestEmptyPlaintext(t *testing.T) {
	// Test with empty plaintext
	plaintext := []byte("")
	password := "password123"

	// Encrypt with the correct password
	ciphertext, err := Encrypt(password, plaintext)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Decrypt the ciphertext
	decryptedText, err := Decrypt(password, ciphertext)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// Check if decrypted text is also empty
	if len(decryptedText) != 0 {
		t.Errorf("Expected empty decrypted text, got: %s", decryptedText)
	}
}
