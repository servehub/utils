package crypto_test

import (
	goCrypto "crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/servehub/utils/crypto"
)

func TestGenerateEdDSAKeyPair(t *testing.T) {
	pub, priv, err := crypto.GenerateEdDSAKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate EdDSA key pair: %v", err)
	}

	pubBytes, err := hex.DecodeString(pub)
	if err != nil {
		t.Fatalf("Failed to decode public key: %v", err)
	}

	privBytes, err := hex.DecodeString(priv)
	if err != nil {
		t.Fatalf("Failed to decode private key: %v", err)
	}

	if len(pubBytes) != ed25519.PublicKeySize {
		t.Errorf("Invalid public key size: got %d, want %d", len(pubBytes), ed25519.PublicKeySize)
	}

	if len(privBytes) != ed25519.SeedSize {
		t.Errorf("Invalid private key size: got %d, want %d", len(privBytes), ed25519.SeedSize)
	}

	privKey := ed25519.NewKeyFromSeed(privBytes)

	sig, err := privKey.Sign(rand.Reader, []byte("test message"), goCrypto.Hash(0))
	if err != nil {
		t.Fatalf("Failed to sign message: %v", err)
	}

	if !ed25519.Verify(pubBytes, []byte("test message"), sig) {
		t.Errorf("Generated key pair is not valid")
	}
}
