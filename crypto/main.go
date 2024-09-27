package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

/**
 * GenerateEdDSAKeyPair generates a new EdDSA key pair.
 * It returns the public key and the private key in hex-encoded strings.
 */
func GenerateEdDSAKeyPair() (string, string, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate EdDSA key pair: %v", err)
	}

	return hex.EncodeToString(pub), hex.EncodeToString(priv.Seed()), nil
}
