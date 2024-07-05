package commands

import (
	"fmt"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip06"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	seedPhrase string
	privateKey string
	seedSource string
	keySource  string
)

var SeedGen = &cobra.Command{
	Use:   "seed-gen [-s seed | -f seed-filename]",
	Short: "Generates and/or saves a new nostr seed id file",
	Run:   seedGen,
}

var KeyGen = &cobra.Command{
	Use:   "key-gen [-s seed | -f seed-filename]",
	Short: "Generates and/or saves new key pair",
	Run:   keyGen,
}

func init() {
	SeedGen.Flags().StringVarP(&seedPhrase, "seed", "s", "", "user provided seed phrase")
	SeedGen.Flags().StringVarP(&seedSource, "seed-source", "f", "", "user provided plain-text seed phrase source filename")

	KeyGen.Flags().StringVarP(&seedPhrase, "seed", "s", "", "user provided seed phrase")
	KeyGen.Flags().StringVarP(&privateKey, "private-key", "p", "", "user provided private key")
	KeyGen.Flags().StringVarP(&seedSource, "seed-source", "f", "", "user provided plain-text seed phrase source filename")
	KeyGen.Flags().StringVarP(&keySource, "key-source", "k", "", "user provided private key source filename")
}

func seedGen(cmd *cobra.Command, _ []string) {
	cmd.Println("Generating new nostr id...")
	generateSeed(cmd)
}

func resolveSeed(cmd *cobra.Command) string {
	if seedPhrase != "" {
		return seedPhrase
	}

	if seedSource != "" {
		seedSourcePassphrase := promptSecret(cmd, fmt.Sprintf("Enter passphrase to read seed phrase file (%s): ", seedSource))
		sourceSeedBytes, err := readNostrFile(seedSource, seedSourcePassphrase)
		if err != nil {
			cmd.Printf("failed to read seed phrase file: %v", err)
			os.Exit(1)
		}
		seedFromSource := strings.TrimSuffix(string(sourceSeedBytes), "\n")
		if !nip06.ValidateWords(seedFromSource) {
			cmd.Println("seed phrase is invalid")
			os.Exit(1)
		}
		return seedFromSource
	}

	newSeedPhrase, err := nip06.GenerateSeedWords()
	if err != nil {
		cmd.Printf("failed to generate seed phrase: %v", err)
		os.Exit(1)
	}
	return newSeedPhrase
}

func generateSeed(cmd *cobra.Command) string {

	home, err := os.UserHomeDir()
	if err != nil {
		cmd.Printf("failed to get home directory: %v", err)
		os.Exit(1)
	}

	defaultSeedFilename := fmt.Sprintf("%s/.nostr/seed", home)
	seedFilename := promptString(cmd, fmt.Sprintf("Enter filename to save private seed (%s): ", defaultSeedFilename))
	if seedFilename == "" {
		seedFilename = defaultSeedFilename
	}

	passPhrase := promptSecret(cmd, "Enter passphrase (empty for no passphrase): ")
	repassPhrase := promptSecret(cmd, "Enter passphrase again: ")
	if passPhrase != repassPhrase {
		cmd.Println("Passphrases do not match")
		os.Exit(1)
	}

	//check if file already exists, ask for passphrase to overwrite
	if _, err := os.Stat(seedFilename); err == nil {
		if overwrite := promptBool(cmd, "Seed file already exists. Overwrite?", false); overwrite {
			overwritePassphrase := promptSecret(cmd, fmt.Sprintf("Enter passphrase to overwrite existing seed file (%s): ", seedFilename))
			_, err = readNostrFile(seedFilename, overwritePassphrase)
			if err != nil {
				cmd.Printf("failed to read seed phrase file: %v", err)
				os.Exit(1)
			}
		}
	}

	err = os.MkdirAll(fmt.Sprintf("%s/.nostr", home), 0700)
	if err != nil {
		cmd.Printf("failed to create .nostr directory: %v", err)
		os.Exit(1)
	}

	seedPhrase := resolveSeed(cmd)

	err = writeNostrFile(seedFilename, []byte(seedPhrase), passPhrase)
	if err != nil {
		cmd.Printf("failed to write seed phrase file: %v", err)
		os.Exit(1)
	}

	cmd.Println("Seed phrase saved to:", seedFilename)
	cmd.Println("Seed phrase:", seedPhrase)
	cmd.Println(`=================================================================================
This seed phrase is your identity in nostr! Keep it secret and back it up safely!
=================================================================================`)

	return seedFilename
}

func keyGen(cmd *cobra.Command, _ []string) {
	cmd.Println("Generating new key pair...")
	generateKey(cmd)
}

func resolveKey(cmd *cobra.Command) string {
	if privateKey != "" {
		return privateKey
	}

	if keySource != "" {
		keySourcePassphrase := promptSecret(cmd, fmt.Sprintf("Enter passphrase to read key file (%s): ", keySource))
		keyBytes, err := readNostrFile(keySource, keySourcePassphrase)
		if err != nil {
			cmd.Printf("failed to read key file: %v", err)
			os.Exit(1)
		}
		return strings.TrimSuffix(string(keyBytes), "\n")
	}

	seedPhrase := resolveSeed(cmd)
	pvtKey, err := nip06.PrivateKeyFromSeed(nip06.SeedFromWords(seedPhrase))
	if err != nil {
		cmd.Printf("failed to generate private key: %v", err)
		os.Exit(1)
	}
	return pvtKey
}

func generateKey(cmd *cobra.Command) {

	home, err := os.UserHomeDir()
	if err != nil {
		cmd.Printf("failed to get home directory: %v", err)
		os.Exit(1)
	}

	defaultKeyFilename := fmt.Sprintf("%s/.nostr/key", home)
	keyFilename := promptString(cmd, fmt.Sprintf("Enter filename to save key pair (%s): ", defaultKeyFilename))
	if keyFilename == "" {
		keyFilename = defaultKeyFilename
	}

	passPhrase := promptSecret(cmd, "Enter passphrase (empty for no passphrase): ")
	repassPhrase := promptSecret(cmd, "Enter passphrase again: ")
	if passPhrase != repassPhrase {
		cmd.Println("Passphrases do not match")
		os.Exit(1)
	}

	//check if file already exists, ask for passphrase to overwrite
	if _, err := os.Stat(keyFilename); err == nil {
		if overwrite := promptBool(cmd, "Key file already exists. Overwrite?", false); overwrite {
			overwritePassphrase := promptSecret(cmd, fmt.Sprintf("Enter passphrase to overwrite existing key file (%s): ", keyFilename))
			// try to read the file with the passphrase
			_, err = readNostrFile(keyFilename, overwritePassphrase)
			if err != nil {
				cmd.Printf("failed to read key file: %v", err)
				os.Exit(1)
			}
		}
	}

	err = os.MkdirAll(fmt.Sprintf("%s/.nostr", home), 0700)
	if err != nil {
		cmd.Printf("failed to create .nostr directory: %v", err)
		os.Exit(1)
	}

	pvtKey := resolveKey(cmd)
	pubKey, err := nostr.GetPublicKey(pvtKey)
	if err != nil {
		cmd.Printf("failed to generate public key: %v", err)
		os.Exit(1)
	}

	err = writeNostrFile(keyFilename, []byte(pvtKey), passPhrase)
	if err != nil {
		cmd.Printf("failed to write key file: %v", err)
		os.Exit(1)
	}

	pubKeyFilename := fmt.Sprintf("%s.pub", keyFilename)
	err = writeNostrFile(pubKeyFilename, []byte(pubKey), passPhrase)
	if err != nil {
		cmd.Printf("failed to write public key file: %v", err)
		os.Remove(keyFilename)
		os.Exit(1)
	}

	npub, err := nip19.EncodePublicKey(pubKey)
	if err != nil {
		cmd.Printf("failed to encode public key: %v", err)
		os.Remove(keyFilename)
		os.Remove(pubKeyFilename)
		os.Exit(1)
	}

	cmd.Printf("Key saved to: %s, %s\n", keyFilename, pubKeyFilename)
	cmd.Printf("Public key: %s\n", npub)
	cmd.Println(`==============================================================================
This key pair is your identity in nostr! Keep it secret and back it up safely!
==============================================================================`)
}
