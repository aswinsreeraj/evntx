package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

func getAESKey() []byte {
	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		key = "default_unsafe_development_key"
	}
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}


func EncryptAES(plainText string) (string, error) {
	if plainText == "" {
		return "", errors.New("cannot encrypt empty string")
	}

	key := getAESKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	dst := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return hex.EncodeToString(dst), nil
}


func DecryptAES(cryptoHex string) (string, error) {
	if cryptoHex == "" {
		return "", errors.New("cannot decrypt empty string")
	}

	dst, err := hex.DecodeString(cryptoHex)
	if err != nil {
		return "", err
	}

	key := getAESKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(dst) < gcm.NonceSize() {
		return "", errors.New("cipher text too short")
	}

	nonce, cipherText := dst[:gcm.NonceSize()], dst[gcm.NonceSize():]

	plainTextBytes, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainTextBytes), nil
}
