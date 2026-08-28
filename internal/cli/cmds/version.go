package cmds

import (
	"fmt"

	"github.com/honakac/go-sack/lib"
	"github.com/spf13/cobra"
)

const (
	Version = "1.0.4"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("go-sack cli v%s\n", Version)
		fmt.Printf("go-sack library v%s\n", lib.Version)
	},
}
