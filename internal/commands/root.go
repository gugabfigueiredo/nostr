/*
Copyright © 2024 Guga Figueiredo gugafigueiredo@primal.net
*/

package nostr

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/gugabfigueiredo/nostr-cli/internal/commands"
)

var rootCmd = &cobra.Command{
	Use:   "nostr <command>",
	Short: "All your data are belong to you!",
	Long: `Nostr is a cli tool to manage your nostr data.
Generate a new identity, sign and verify messages, encrypt and decrypt data, and more.
All your data are belong to you!`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(commands.Keygen)
}
