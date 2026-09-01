package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/honakac/sack"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list ARCHIVE",
	Short: "Create sack archive",
	Args:  cobra.ExactArgs(1),
	Run:   listRun,
}

func listRun(cmd *cobra.Command, args []string) {
	a, err := os.Open(args[0])
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	r, err := sack.NewReader(a)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Count files: %d\n", len(r.Files))

	for k, f := range r.Files {
		fmt.Printf("%s '%s' %d bytes\n", os.FileMode(f.Mode), k, f.Size)
	}
}
