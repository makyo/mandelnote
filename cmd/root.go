package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "mandelnote <note file>",
	Short: "Run mandelnote",
	Long: `Run mandelnote

Mandelnote is a fractal note-taking app based loosely on the snowflake method.
It allows you to take notes - one item per 'card', and then expand cards into
additional levels of notes. Notes can be shifted into higher or lower levels of
detail. Projects can be exported to markdown format for use in turning notes
into a paper or story.`,
	Args: cobra.MaximumArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
	Version: "0.0.1",
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Yike: %v", err)
		os.Exit(1)
	}
}
