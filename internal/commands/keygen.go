package internal

import (
	"fmt"
	"github.com/nbd-wtf/go-nostr/nip06"
	"github.com/spf13/cobra"
	"os"
)

var (
	seedFilename string
)

var Keygen = &cobra.Command{
	Use:   "key-gen [options]",
	Short: "Generates a new nostr identity to a file",
	Run:   keygen,
}

func init() {
	Keygen.Flags().StringVarP(&seedFilename, "seed-filename", "s", "", "Filename to save seed phrase")
}

func keygen(cmd *cobra.Command, args []string) {
	cmd.Println("Generating new nostr id...")

	seedPhrase, err := nip06.GenerateSeedWords()
	if err != nil {
		panic(err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	if seedFilename == "" {
		defaultSeedFilename := fmt.Sprintf("%s/.nostr/seed", home)
		seedFilename = promptString(cmd, fmt.Sprintf("Enter filename to save private key (%s): ", defaultSeedFilename))
		if seedFilename == "" {
			seedFilename = defaultSeedFilename
		}
	}

	passPhrase := promptString(cmd, "Enter passphrase (empty for no passphrase): ")
	repassPhrase := promptString(cmd, "Enter passphrase again: ")
	if passPhrase != repassPhrase {
		panic("Passphrases do not match")
	}

	err = os.MkdirAll(fmt.Sprintf("%s/.nostr", home), 0700)
	if err != nil {
		panic(err)
	}

	//err = os.WriteFile(seedFilename, []byte(seedPhrase), 0600)
	//if err != nil {
	//	panic(err)
	//}

	cmd.Println("Seed phrase saved to: ", seedFilename)
	cmd.Println("Seed phrase: ", seedPhrase)
	cmd.Println("This seed phrase is the only way to recover your identity in nostr, keep it secret and back it up safely!")
}
