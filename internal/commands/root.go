/*
Copyright © 2024 Guga Figueiredo gugafigueiredo@primal.net
*/

package commands

import (
	"os"

	"github.com/spf13/cobra"
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
	rootCmd.PersistentFlags().BoolP("help", "h", false, "help for nostr")

	rootCmd.AddCommand(Keygen)
}
