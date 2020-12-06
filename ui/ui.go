package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/makyo/ansigo"
	"github.com/makyo/gotui"

	"github.com/makyo/mandelnote/notebook"
)

var (
	columns int = 12
)

type tui struct {
	g             *gotui.Gui
	nb            *notebook.Notebook
	colWidth      int
	cardWidth     int
	modalOpen     bool
	cards         []*card
	cardNameIndex int
	currentDepth  int
}

type card struct {
	card    notebook.Card
	view    *gotui.View
	name    string
	depth   int
	current bool
}

func (t *tui) onResize(g *gotui.Gui, x, y int) error {
	t.colWidth = x / columns
	t.cardWidth = 1
	for t.cardWidth*t.colWidth < 90 {
		t.cardWidth++
	}
	if err := t.setTitle(g); err != nil {
		return err
	}
	return t.drawCards(g, x)
}

func (t *tui) setTitle(g *gotui.Gui) error {
	maxX, _ := g.Size()
	if v, err := g.SetView("title", -1, 0, maxX+1, 1); err != nil {
		return err
	} else {
		title := fmt.Sprintf("   %s ── %s  ", t.nb.Title, t.nb.Author)
		helpMsg := "  Help: ctrl+H  "
		v.Clear()
		fmt.Fprintf(v, ansigo.MaybeApplyWithReset("underline+8", fmt.Sprintf("%s%s%s",
			title,
			strings.Repeat(" ", maxX-len(title)-len(helpMsg)+3),
			helpMsg)))
	}
	return nil
}

func (t *tui) quit(g *gotui.Gui, v *gotui.View) error {
	return gotui.ErrQuit
}

func (t *tui) keybindings(g *gotui.Gui) error {
	if err := g.SetKeybinding("", gotui.KeyCtrlQ, gotui.ModNone, t.quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyCtrlH, gotui.ModNone, t.showHelp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyArrowUp, gotui.ModAlt, t.cycleUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyArrowDown, gotui.ModAlt, t.cycleUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyArrowRight, gotui.ModAlt, t.enter); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyArrowLeft, gotui.ModAlt, t.exit); err != nil {
		return err
	}
	if err := g.SetKeybinding("modal", gotui.KeyEnter, gotui.ModNone, t.closeModal); err != nil {
		return err
	}
	if err := g.SetKeybinding("modal", gotui.KeyArrowUp, gotui.ModNone, t.scrollModalUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("modal", gotui.KeyArrowDown, gotui.ModNone, t.scrollModalDown); err != nil {
		return err
	}
	return nil
}

func (t *tui) layout(g *gotui.Gui) error {
	maxX, _ := g.Size()
	if v, err := g.SetView("title", 0, -1, maxX, 1); err != nil {
		if err != gotui.ErrUnknownView {
			return err
		} else {
			v.Frame = false
		}
	}
	g.Update(func(gg *gotui.Gui) error {
		maxX, maxY := g.Size()
		return t.onResize(gg, maxX, maxY)
	})
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
