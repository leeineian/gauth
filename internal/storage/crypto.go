package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	saltLen  = 16
	nonceLen = 12
	keyLen   = 32
	timeCost = 3
	memCost  = 64 * 1024
	threads  = 4
)

// encrypt generates a derived key and encrypts data using AES-GCM-256
func encrypt(data []byte, password string) ([]byte, error) {
	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	key := argon2.IDKey([]byte(password), salt, timeCost, memCost, threads, keyLen)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, nonceLen)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Layout: [salt][nonce][ciphertext]
	cipherText := gcm.Seal(nil, nonce, data, nil)

	out := make([]byte, 0, len(salt)+len(nonce)+len(cipherText))
	out = append(out, salt...)
	out = append(out, nonce...)
	out = append(out, cipherText...)

	return out, nil
}

// decrypt derives the key from the password and decrypts the provided data
func decrypt(data []byte, password string) ([]byte, error) {
	if len(data) < saltLen+nonceLen {
		return nil, fmt.Errorf("ciphertext too short")
	}

	salt := data[:saltLen]
	nonce := data[saltLen : saltLen+nonceLen]
	cipherText := data[saltLen+nonceLen:]

	key := argon2.IDKey([]byte(password), salt, timeCost, memCost, threads, keyLen)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plainText, nil
}
