package aes_test

import (
	"testing"

	"github.com/servehub/utils/crypto/aes"
)

func TestEncryptDecrypt(t *testing.T) {
	password := "my-password"
	data := "Hello, World!"

	encrypted, err := aes.Encrypt(password, data)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := aes.Decrypt(password, encrypted)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if decrypted != data {
		t.Errorf("Decryption result mismatch: got %v, want %v", decrypted, data)
	}
}

func TestEncryptDecryptWithPreGenerated(t *testing.T) {
	password := "LokRXbLjcWqM0ClMes4Hdaqfe3kiMLjP"
	data := " test data 444 "

	encrypted := "ARxPQU90KCm0b1J6PYqAPL9FdPeATYr6GsZaaapQoTPsxJAXtr1wXq3e+l1RznLE9wFHLzWxCPXEtW+f"

	decrypted, err := aes.Decrypt(password, encrypted)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if decrypted != data {
		t.Errorf("Decryption result mismatch: got %v, want %v", decrypted, data)
	}
}
