package main

import (
	"log"

	"github.com/honakac/go-sack/cli/internal/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sack",
	Short: "Sack - Just a small sack for all your files. A minimal tar alternative written in Go.",
}

func main() {
	rootCmd.AddCommand(
		cmd.VersionCmd,
		cmd.CreateCmd,
		cmd.ListCmd,
		cmd.ExtractCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
