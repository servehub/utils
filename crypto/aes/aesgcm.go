package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	version    = 1
	saltSize   = 16
	ivSize     = 12
	iterations = 200000
	keySize    = 32
)

func generateSecretKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, iterations, keySize, sha256.New)
}

func Encrypt(password, data string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key := generateSecretKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := make([]byte, ivSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	encrypted := aesGCM.Seal(nil, iv, []byte(data), nil)

	buf := make([]byte, 1+saltSize+ivSize+len(encrypted))
	buf[0] = version
	copy(buf[1:], salt)
	copy(buf[1+saltSize:], iv)
	copy(buf[1+saltSize+ivSize:], encrypted)

	return base64.StdEncoding.EncodeToString(buf), nil
}

func Decrypt(password, encryptedData string) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	if encrypted[0] != version {
		return "", errors.New("unknown encryption version")
	}

	salt := encrypted[1 : 1+saltSize]
	key := generateSecretKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	iv := encrypted[1+saltSize : 1+saltSize+ivSize]
	ciphertext := encrypted[1+saltSize+ivSize:]

	decrypted, err := aesGCM.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
