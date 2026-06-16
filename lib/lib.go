package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"golang.org/x/crypto/argon2"
)

const (
	// defaultSecretKeyPrefix is an application-specific secret used as part of the key derivation.
	// Override it at runtime via ENCDEC_SECRET_PREFIX for production deployments.
	defaultSecretKeyPrefix = "jY-1"

	// derivedKeyLen is the target length (in bytes) for the AES key.
	// AES accepts 16/24/32 bytes. We use 32 bytes (AES-256).
	derivedKeyLen = 32

	// limits to mitigate trivial local DoS via huge inputs
	maxCiphertextHexLen = 1 << 20 // 1 MiB of hex (~512KiB bytes)
	maxPlaintextLen     = 1 << 20 // 1 MiB

	// passphrase format marker
	passphrasePrefix = "p1:"
)

// GenerateKey returns a machine-bound AES-256 key encoded in hex (64 hex chars).
// It is derived from:
//   - secret prefix (defaultSecretKeyPrefix or ENCDEC_SECRET_PREFIX)
//   - machine id
//   - OS
//
// NOTE: This is a deterministic, machine-bound derivation; it is not a password-based KDF.
// If machineID changes (or you move the ciphertext to another machine), decryption will fail.
//
// It returns an error (instead of panicking) when the machine id cannot be read,
// so the CLI can honor its contract: print the error on stderr and exit 1.
func GenerateKey() (string, error) {
	key, err := generateKeyBytes()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

// generateKeyBytes derives a 32-byte key using SHA-256 and returns raw bytes.
// This avoids hex roundtrips and ensures the effective AES key size is 32 bytes.
func generateKeyBytes() ([]byte, error) {
	mID, err := machineid.ID()
	if err != nil {
		return nil, fmt.Errorf("machine id: %w", err)
	}

	secretPrefix := os.Getenv("ENCDEC_SECRET_PREFIX")
	if secretPrefix == "" {
		secretPrefix = defaultSecretKeyPrefix
	}

	uniqueID := secretPrefix + mID + runtime.GOOS
	sum := sha256.Sum256([]byte(uniqueID))
	return sum[:], nil
}

// EncryptString encrypts plaintext using AES-GCM.
// keyString must be a hex string representing 16/24/32 bytes.
func EncryptString(plaintext, keyString string) (string, error) {
	if len(plaintext) > maxPlaintextLen {
		return "", fmt.Errorf("plaintext too large (max %d bytes)", maxPlaintextLen)
	}

	key, err := decodeHexKey(keyString)
	if err != nil {
		return "", err
	}

	return encryptWithKeyBytes([]byte(plaintext), key)
}

// DecryptString decrypts ciphertext (hex-encoded) using AES-GCM.
// keyString must be a hex string representing 16/24/32 bytes.
func DecryptString(cipherHex, keyString string) (string, error) {
	if len(cipherHex) > maxCiphertextHexLen {
		return "", fmt.Errorf("ciphertext too large (max %d hex chars)", maxCiphertextHexLen)
	}

	key, err := decodeHexKey(keyString)
	if err != nil {
		return "", err
	}

	plain, err := decryptWithKeyBytes(cipherHex, key)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

// EncryptStringWithPassphrase encrypts plaintext using a passphrase-based key (Argon2id).
// It returns a self-contained string prefixed with "p1:".
// Format: p1:<salt_b64>:<cipher_hex>
func EncryptStringWithPassphrase(plaintext, passphrase string) (string, error) {
	if passphrase == "" {
		return "", errors.New("passphrase is empty")
	}
	if len(plaintext) > maxPlaintextLen {
		return "", fmt.Errorf("plaintext too large (max %d bytes)", maxPlaintextLen)
	}

	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("salt: %w", err)
	}

	key := argon2.IDKey([]byte(passphrase), salt, 1, 64*1024, 4, derivedKeyLen) // ~64MiB memory
	cipherHex, err := encryptWithKeyBytes([]byte(plaintext), key)
	if err != nil {
		return "", err
	}

	return passphrasePrefix + base64.RawStdEncoding.EncodeToString(salt) + ":" + cipherHex, nil
}

// DecryptStringWithPassphrase decrypts a ciphertext produced by EncryptStringWithPassphrase.
// Accepted formats:
//   - p1:<salt_b64>:<cipher_hex>
//   - or if cipherText doesn't start with "p1:", it is treated as plain hex (compat) and salt is not available -> error.
func DecryptStringWithPassphrase(cipherText, passphrase string) (string, error) {
	if passphrase == "" {
		return "", errors.New("passphrase is empty")
	}
	if len(cipherText) > maxCiphertextHexLen+64 { // small cushion for prefix/salt
		return "", fmt.Errorf("ciphertext too large")
	}

	if !strings.HasPrefix(cipherText, passphrasePrefix) {
		return "", errors.New("invalid passphrase ciphertext format (missing p1: prefix)")
	}

	parts := strings.SplitN(strings.TrimPrefix(cipherText, passphrasePrefix), ":", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid passphrase ciphertext format")
	}

	saltB64 := parts[0]
	cipherHex := parts[1]

	salt, err := base64.RawStdEncoding.DecodeString(saltB64)
	if err != nil {
		return "", fmt.Errorf("salt base64: %w", err)
	}
	if len(salt) < 8 {
		return "", errors.New("salt too short")
	}

	key := argon2.IDKey([]byte(passphrase), salt, 1, 64*1024, 4, derivedKeyLen)
	plain, err := decryptWithKeyBytes(cipherHex, key)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func encryptWithKeyBytes(plaintext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("gcm: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce: %w", err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

func decryptWithKeyBytes(cipherHex string, key []byte) ([]byte, error) {
	enc, err := hex.DecodeString(cipherHex)
	if err != nil {
		return nil, fmt.Errorf("ciphertext is not valid hex: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("gcm: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(enc) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	return plaintext, nil
}

func decodeHexKey(keyString string) ([]byte, error) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return nil, fmt.Errorf("key is not valid hex: %w", err)
	}
	switch len(key) {
	case 16, 24, 32:
		return key, nil
	default:
		return nil, fmt.Errorf("invalid key length: got %d bytes, want 16/24/32", len(key))
	}
}
