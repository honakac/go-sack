package cmds

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/honakac/go-sack/lib"
	"github.com/spf13/cobra"
)

type extractOptions struct {
	dir *string
}

var extractOpts = new(extractOptions)
var ExtractCmd = &cobra.Command{
	Use:   "extract ARCHIVE",
	Short: "Create sack archive",
	Args:  cobra.ExactArgs(1),
	Run:   ExtractRun,
}

func init() {
	extractOpts.dir = ExtractCmd.Flags().StringP("directory", "C", ".", "change to directory DIR")
}

func ExtractRun(cmd *cobra.Command, args []string) {
	a, err := os.Open(args[0])
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	r, err := lib.NewReader(a)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range r.Files {
		f.Name = lib.CleanPath(f.Name, true)

		reader, err := r.OpenFile(f)
		if err != nil {
			log.Fatal(err)
		}

		origDir := filepath.Dir(f.Name)
		dir := fmt.Sprintf("%s/%s", *extractOpts.dir, origDir)
		if origDir != "." && origDir != "/" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Fatal(err)
			}
		}

		newName := fmt.Sprintf("%s/%s", *extractOpts.dir, f.Name)
		writer, err := os.OpenFile(newName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(f.Info.Mode))
		if err != nil {
			log.Fatal(err)
		}
		if _, err := io.Copy(writer, reader); err != nil {
			log.Fatal(err)
		}
	}
}
