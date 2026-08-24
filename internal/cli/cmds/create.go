package cmds

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/honakac/go-sack/lib"
	"github.com/spf13/cobra"
)

type createOptions struct {
	verbose *bool
}

var createOpts = new(createOptions)
var CreateCmd = &cobra.Command{
	Use:   "create ARCHIVE FILE...",
	Short: "Create sack archive",
	Args:  cobra.ExactArgs(2),
	Run:   createRun,
}
var outputFile string

func init() {
	createOpts.verbose = CreateCmd.Flags().BoolP("verbose", "v", false, "Show debug messages")
}

func addFile(w *lib.Writer, filepath string, name string, s os.FileInfo) error {
	if lib.CleanPath(filepath, false) == outputFile {
		return nil
	}

	r, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer r.Close()

	w.AddStream(lib.CleanPath(name, true), s.Size(), uint32(s.Mode()), r)

	return nil
}

func scanFolder(w *lib.Writer, dir string, name string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		filepath := fmt.Sprintf("%s/%s", dir, file.Name())
		namepath := fmt.Sprintf("%s/%s", name, file.Name())

		if file.Type().IsRegular() {
			s, err := os.Stat(filepath)
			if err != nil {
				return err
			}

			if *createOpts.verbose {
				fmt.Printf("Packing %s...\n", namepath)
			}

			if err := addFile(w, filepath, namepath, s); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := scanFolder(w, filepath, namepath); err != nil {
				return err
			}
		}
	}
	return nil
}

func createRun(cmd *cobra.Command, args []string) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	outputFile = lib.CleanPath(fmt.Sprintf("%s/%s", currentDir, args[0]), false)

	a, err := os.OpenFile(args[0], os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	w, err := lib.NewWriter(a)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := w.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	for _, arg := range args[1:] {
		s, err := os.Stat(arg)
		if err != nil {
			log.Fatal(err)
		}

		if s.IsDir() {
			if err := scanFolder(w, arg, filepath.Base(arg)); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := addFile(w, arg, filepath.Base(arg), s); err != nil {
				log.Fatal(err)
			}
		}
	}
}
