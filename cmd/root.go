package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/makyo/mandelnote/notebook"
	"github.com/makyo/mandelnote/ui"
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
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		nb, err := notebook.Open(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening notebook: %v", err)
			os.Exit(1)
			return
		}
		tui := ui.New(nb)
		tui.Run()
	},
	Version: "0.0.1",
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Yike: %v", err)
		os.Exit(1)
	}
}
