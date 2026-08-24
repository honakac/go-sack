package cli

import (
	"log"

	"github.com/honakac/go-sack/internal/cli/cmds"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sack",
	Short: "Sack - Just a small sack for all your files. A minimal tar alternative written in Go.",
}

func Run() {
	rootCmd.AddCommand(
		cmds.VersionCmd,
		cmds.CreateCmd,
		cmds.ListCmd,
		cmds.ExtractCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
