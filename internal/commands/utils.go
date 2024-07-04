package commands

import (
	"bufio"
	"github.com/spf13/cobra"
	"strings"
)

func promptString(cmd *cobra.Command, prompt string) string {
	cmd.Print(prompt)
	reader := bufio.NewReader(cmd.InOrStdin())
	text, _ := reader.ReadString('\n')
	return strings.TrimSuffix(text, "\n")
}
