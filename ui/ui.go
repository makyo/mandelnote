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
	columns    int = 12
	confirmMsg     = fmt.Sprintf(" Confirm: %ses/%so ", ansigo.MaybeApplyWithReset("underline", "Y"), ansigo.MaybeApplyWithReset("underline", "N"))
)

type tui struct {
	g             *gotui.Gui
	nb            *notebook.Notebook
	colWidth      int
	cardWidth     int
	modalOpen     bool
	editorOpen    bool
	focusMode     bool
	cards         []*card
	cardNameIndex int
	currentDepth  int
	currentY      int
	currentHeight int
	inputs        map[string]string
	confirmYesFn  func(*gotui.Gui) error
	confirmNoFn   func(*gotui.Gui) error
}

func (t *tui) onResize(g *gotui.Gui, x, y int) error {
	t.colWidth = x / columns
	if t.colWidth < 6 {
		t.colWidth = 6
	}
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
		helpMsg := "  Hit ? for help  "
		v.Clear()
		fmt.Fprintf(v, ansigo.MaybeApplyWithReset("underline+8", fmt.Sprintf("%s%s%s",
			title,
			strings.Repeat(" ", maxX-len(title)-len(helpMsg)+3),
			helpMsg)))
	}
	return nil
}

func (t *tui) editMetadata(g *gotui.Gui, v *gotui.View) error {
	if t.modalOpen {
		if t.editorOpen {
			g.CurrentView().EditWrite('e')
		}
		return nil
	}
	return nil
}

func (t *tui) save(g *gotui.Gui, v *gotui.View) error {
	if t.modalOpen {
		if t.editorOpen {
			g.CurrentView().EditWrite('s')
		}
		return nil
	}
	return t.nb.Save()
}

func (t *tui) saveAs(g *gotui.Gui, v *gotui.View) error {
	if t.modalOpen {
		return nil
	}
	return nil
}

func (t *tui) quit(g *gotui.Gui, v *gotui.View) error {
	if t.modalOpen {
		return nil
	}
	if t.nb.Dirty() {
		maxX, maxY := g.Size()
		if v, err := g.SetView("confirm", t.colWidth*2, maxY/2-2, maxX-t.colWidth*2, maxY/2+2); err != nil {
			if err != gotui.ErrUnknownView {
				return err
			}
			t.modalOpen = true
			v.Frame = true
			v.FrameFgColor = gotui.ColorCyan | gotui.AttrBold
			v.TitleFgColor = gotui.AttrBold
			v.Title = " Quit "
			v.Wrap = true
			v.WordWrap = true
			fmt.Fprint(v, "You have unsaved changes. Would you like to save before quitting?")
			g.SetCurrentView("confirm")
		}
		if v, err := g.SetView("confirmActions", maxX-t.colWidth*2-20, maxY/2+1, maxX-t.colWidth*2-2, maxY/2+3); err != nil {
			if err != gotui.ErrUnknownView {
				return err
			}
			v.Frame = false
			fmt.Fprintf(v, confirmMsg)
		}
		g.SetCurrentView("confirm")
		t.confirmYesFn = func(gg *gotui.Gui) error {
			err := t.nb.Save()
			if err != nil {
				return err
			}
			return gotui.ErrQuit
		}
		t.confirmNoFn = func(gg *gotui.Gui) error {
			return gotui.ErrQuit
		}
		return nil
	} else {
		return gotui.ErrQuit
	}
}

func (t *tui) keybindings(g *gotui.Gui) error {
	// Notebook-wide keybindings
	if err := g.SetKeybinding("", '?', gotui.ModNone, t.showHelp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'e', gotui.ModNone, t.editMetadata); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 's', gotui.ModNone, t.save); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyCtrlS, gotui.ModNone, t.saveAs); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyCtrlQ, gotui.ModNone, t.quit); err != nil {
		return err
	}

	// Card tasks
	if err := g.SetKeybinding("", 'n', gotui.ModNone, t.newCard); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'N', gotui.ModNone, t.newChild); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyEnter, gotui.ModNone, t.edit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'f', gotui.ModNone, t.focus); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyTab, gotui.ModNone, t.toggleEdit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyCtrlW, gotui.ModNone, t.closeEditor); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'p', gotui.ModNone, t.promote); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'P', gotui.ModNone, t.promoteAll); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'm', gotui.ModNone, t.mergeDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'M', gotui.ModNone, t.mergeUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'd', gotui.ModNone, t.moveDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'u', gotui.ModNone, t.moveUp); err != nil {
		return err
	}

	// Card movement
	if err := g.SetKeybinding("", gotui.KeyArrowUp, gotui.ModNone, t.cycleUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyArrowDown, gotui.ModNone, t.cycleDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyArrowRight, gotui.ModNone, t.enter); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gotui.KeyArrowLeft, gotui.ModNone, t.exit); err != nil {
		return err
	}

	// Modal tasks
	if err := g.SetKeybinding("modal", 'q', gotui.ModNone, t.closeModal); err != nil {
		return err
	}
	if err := g.SetKeybinding("modal", gotui.KeyArrowUp, gotui.ModNone, t.scrollModalUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("modal", gotui.KeyArrowDown, gotui.ModNone, t.scrollModalDown); err != nil {
		return err
	}

	// Confirm tasks
	if err := g.SetKeybinding("confirm", 'y', gotui.ModNone, t.confirmYes); err != nil {
		return err
	}
	if err := g.SetKeybinding("confirm", 'n', gotui.ModNone, t.confirmNo); err != nil {
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

	t.g.Cursor = false
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
