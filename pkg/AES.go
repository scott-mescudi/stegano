package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"golang.org/x/crypto/scrypt"
	"io"
)

func deriveKey(password, salt []byte) ([]byte, error) {
	const keyLen = 32 // AES-256 key size
	return scrypt.Key(password, salt, 32768, 8, 1, keyLen)
}

func Encrypt(password string, plaintext []byte) (ciphertext []byte, err error) {
	// Generate a random salt
	salt := make([]byte, 16)
	if _, err = io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	// Derive the AES key from the password and salt
	key, err := deriveKey([]byte(password), salt)
	if err != nil {
		return nil, err
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a random nonce for AES-GCM
	nonce := make([]byte, 12) // GCM nonce size
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Create the GCM cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Encrypt the plaintext
	ct := aesGCM.Seal(nil, nonce, plaintext, nil)
	salt = append(salt, nonce...)
	salt = append(salt, ct...)
	return salt, nil
}

func Decrypt(password string, ciphertext []byte) ([]byte, error) {
	key, err := deriveKey([]byte(password), ciphertext[:16])
	if err != nil {
		return nil, err
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create the GCM cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Decrypt the ciphertext
	plaintext, err := aesGCM.Open(nil, ciphertext[16:28], ciphertext[28:], nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
