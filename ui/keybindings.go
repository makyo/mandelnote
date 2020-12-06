package ui

import (
	"github.com/makyo/gotui"
)

// Notes:
// Vim mode - hit Esc first, then leave off ctrl
// Tab - moves between title and body of card
// Ctrl-RArr - moves to first child card
// Ctrl-LArr - moves to parent card
// Ctrl-DArr - moves to next card
// Ctrl-UArr - moves to prev card
// Ctrl-p - Promotes
// Ctrl-P - Promotes all
// Ctrl-n - Adds card at current level and goes to it
// Ctrl-N - Adds child card and goes to it
// Ctrl-e - Edit metadata
// Ctrl-s - Save
// Ctrl-q - Quit

func (t *tui) quit(g *gotui.Gui, v *gotui.View) error {
	return gotui.ErrQuit
}

func (t *tui) cycleUp(g *gotui.Gui, v *gotui.View) error {
	maxX, _ := g.Size()
	t.nb.Cycle(-1)
	g.Update(func(gg *gotui.Gui) error {
		return t.drawCards(gg, maxX)
	})
	return nil
}

func (t *tui) cycleDown(g *gotui.Gui, v *gotui.View) error {
	maxX, _ := g.Size()
	t.nb.Cycle(1)
	g.Update(func(gg *gotui.Gui) error {
		return t.drawCards(gg, maxX)
	})
	return nil
}

func (t *tui) enter(g *gotui.Gui, v *gotui.View) error {
	maxX, _ := g.Size()
	t.nb.Enter()
	g.Update(func(gg *gotui.Gui) error {
		return t.drawCards(gg, maxX)
	})
	return nil
}

func (t *tui) exit(g *gotui.Gui, v *gotui.View) error {
	maxX, _ := g.Size()
	t.nb.Exit()
	g.Update(func(gg *gotui.Gui) error {
		return t.drawCards(gg, maxX)
	})
	return nil
}

func (t *tui) showHelp(g *gotui.Gui, v *gotui.View) error {
	t.createModal("Help", helpText)
	return nil
}

func (t *tui) scrollModalUp(g *gotui.Gui, v *gotui.View) error {
	if !t.modalOpen {
		return nil
	}
	_, y := v.Origin()
	if y != 0 {
		v.SetOrigin(0, y-1)
	}
	return nil
}

func (t *tui) scrollModalDown(g *gotui.Gui, v *gotui.View) error {
	_, y := v.Origin()
	_, maxY := v.Size()
	lines := len(v.ViewBufferLines())
	if y+maxY <= lines {
		v.SetOrigin(0, y+1)
	}
	return nil
}

func (t *tui) closeModal(g *gotui.Gui, v *gotui.View) error {
	if !t.modalOpen {
		return nil
	}
	if err := g.DeleteView("modal"); err != nil {
		return err
	}
	if err := g.DeleteView("modalTitle"); err != nil {
		return err
	}
	if err := g.DeleteView("modalHelp"); err != nil {
		return err
	}
	t.modalOpen = false
	return nil
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
