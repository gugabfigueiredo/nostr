package commands

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/pbkdf2"
	"os"
	"strings"
)

func promptString(cmd *cobra.Command, prompt string) string {
	cmd.Print(prompt)
	reader := bufio.NewReader(cmd.InOrStdin())
	text, _ := reader.ReadString('\n')
	return strings.TrimSuffix(text, "\n")
}

func writeNostrFile(filename string, data []byte, passphrase string) error {
	if passphrase == "" {
		return os.WriteFile(filename, data, 0600)
	}
	// Generate a salt for the key derivation
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}

	// Derive a key using PBKDF2
	key := pbkdf2.Key([]byte(passphrase), salt, 10000, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return err
	}

	encrypted := gcm.Seal(nonce, nonce, data, nil)
	encrypted = append(salt, encrypted...) // Prepend salt to the encrypted data

	return os.WriteFile(filename, encrypted, 0644)
}

func readNostrFile(inputFile string, passphrase string) ([]byte, error) {
	raw, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, err
	}

	if passphrase == "" {
		return raw, nil
	}

	salt := raw[:16]
	raw = raw[16:]

	key := pbkdf2.Key([]byte(passphrase), salt, 10000, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]
	decrypted, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
