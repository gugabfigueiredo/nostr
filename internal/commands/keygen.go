package commands

import (
	"fmt"
	"github.com/gugabfigueiredo/nostr-cli/internal/io"
	"github.com/nbd-wtf/go-nostr/nip06"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	sourceSeedPhrase   string
	sourceSeedFilename string
)

var Keygen = &cobra.Command{
	Use:   "key-gen [-s seed | -f seed-filename]",
	Short: "Generates or saves",
	Run:   keygen,
}

func init() {
	Keygen.Flags().StringVarP(&sourceSeedPhrase, "seed", "s", "", "user provided seed phrase")
	Keygen.Flags().StringVarP(&sourceSeedFilename, "seed-filename", "f", "", "user provided seed phrase source filename")
}

func keygen(cmd *cobra.Command, args []string) {
	cmd.Println("Generating new nostr id...")

	var seedPhrase string
	if sourceSeedPhrase != "" {
		seedPhrase = sourceSeedPhrase
	} else if sourceSeedFilename != "" {
		seedPhraseBytes, err := os.ReadFile(sourceSeedFilename)
		if err != nil {
			cmd.Printf("failed to read seed phrase file: %v", err)
			os.Exit(1)
		}
		seedPhrase = strings.TrimSuffix(string(seedPhraseBytes), "\n")
	} else {
		newSeedPhrase, err := nip06.GenerateSeedWords()
		if err != nil {
			cmd.Printf("failed to generate seed phrase: %v", err)
			os.Exit(1)
		}
		seedPhrase = newSeedPhrase
	}

	home, err := os.UserHomeDir()
	if err != nil {
		cmd.Printf("failed to get home directory: %v", err)
		os.Exit(1)
	}

	defaultSeedFilename := fmt.Sprintf("%s/.nostr/seed", home)
	seedFilename := promptString(cmd, fmt.Sprintf("Enter filename to save private key (%s): ", defaultSeedFilename))
	if seedFilename == "" {
		seedFilename = defaultSeedFilename
	}

	passPhrase := promptString(cmd, "Enter passphrase (empty for no passphrase): ")
	repassPhrase := promptString(cmd, "Enter passphrase again: ")
	if passPhrase != repassPhrase {
		cmd.Println("Passphrases do not match")
		os.Exit(1)
	}

	err = os.MkdirAll(fmt.Sprintf("%s/.nostr", home), 0700)
	if err != nil {
		cmd.Printf("failed to create .nostr directory: %v", err)
		os.Exit(1)
	}

	err = io.WriteNostrFile(seedFilename, []byte(seedPhrase), passPhrase)
	if err != nil {
		cmd.Printf("failed to write seed phrase file: %v", err)
		os.Exit(1)
	}

	cmd.Println("Seed phrase saved to:", seedFilename)
	cmd.Println("Seed phrase:", seedPhrase)
	cmd.Println(`============================================================================================
This seed phrase is the key to your identity in nostr! Keep it secret and back it up safely!
============================================================================================`)
}
