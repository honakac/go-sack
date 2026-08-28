package cmds

import (
	"fmt"
	"log"
	"os"

	"github.com/honakac/go-sack/lib"
	"github.com/spf13/cobra"
)

// type createOptions struct {
// }

// var createOpts = new(createOptions)
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

	r, err := lib.NewReader(a)
	if err != nil {
		log.Fatal(err)
	}

	for k, f := range r.Files {
		fmt.Printf("%s '%s' %d bytes\n", os.FileMode(f.Mode), k, f.Size)
	}
}
