package ui

import (
	"fmt"
	"os"

	"github.com/makyo/gotui"

	"github.com/makyo/mandelnote/notebook"
)

type tui struct {
	g  *gotui.Gui
	nb *notebook.Notebook
}

func (t *tui) keybindings(g *gotui.Gui) error {
	return nil
}

func (t *tui) onResize(g *gotui.Gui, x, y int) error {
	return nil
}

func (t *tui) layout(g *gotui.Gui) error {
	return nil
}

func (t *tui) Run() {
	var err error
	t.g, err = gotui.NewGui(gotui.Output256)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to create ui: %v", err)
		os.Exit(2)
	}
	defer t.g.Close()

	t.g.Cursor = true
	t.g.Mouse = true

	t.g.SetManagerFunc(t.layout)
	t.g.SetResizeFunc(t.onResize)

	if err := t.keybindings(t.g); err != nil {
		fmt.Fprintf(os.Stderr, "unable to create keybindings: %v", err)
		os.Exit(2)
	}

}

func New(n *notebook.Notebook) *tui {
	return &tui{
		nb: n,
	}
}
