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

func (t *tui) onResize(g *gotui.Gui, x, y int) error {
	return nil
}

func (t *tui) layout(g *gotui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("tree", 1, 1, maxX/4, maxY-1); err != nil {
		if err != gotui.ErrUnknownView {
			return err
		} else {
			v.Title = fmt.Sprintf(" %s ", t.nb.Title)
			fmt.Fprintln(v, "asdf")
		}
	}
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

	if err := t.g.MainLoop(); err != nil && err != gotui.ErrQuit {
		fmt.Fprintf(os.Stderr, "error running mainloop: %v", err)
		os.Exit(3)
	}
}

func New(n *notebook.Notebook) *tui {
	return &tui{
		nb: n,
	}
}
