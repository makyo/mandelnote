package ui

import (
	"fmt"

	"github.com/makyo/ansigo"
	"github.com/makyo/gotui"
)

var (
	helpText string = fmt.Sprintf(
		`

		%s

		A snowflake-method writing tool.

		%s

		%s - edit notebook metadata
		%s - save/save as...
		%s - quit

		%s - new card
		%s - new child card
		%s - edit card
		%s - focus edit
		%s - stop editing
		%s - toggle between editing card title and card body
		%s - promote card
		%s - promote all cards at this level
		%s - merge card up
		%s - merge card down
		%s - move card up
		%s - move card down

		%s - move to next card
		%s - move to previous card
		%s - move to first child card
		%s - move to parent card

		%s

		Docs, examples, and reasoning at %s

		Released under an MIT license. Find the source at %s

		%s

		Madison Scott-Clary - %s
		`,
		ansigo.MaybeApplyWithReset("bold+underline", "Mandelnote"),
		ansigo.MaybeApplyWithReset("underline", "Keybindings"),
		ansigo.MaybeApplyWithReset("cyan", "e       "),
		ansigo.MaybeApplyWithReset("cyan", "s/ctrl+S"),
		ansigo.MaybeApplyWithReset("cyan", "ctrl+Q  "),

		ansigo.MaybeApplyWithReset("cyan", "n      "),
		ansigo.MaybeApplyWithReset("cyan", "N      "),
		ansigo.MaybeApplyWithReset("cyan", "enter  "),
		ansigo.MaybeApplyWithReset("cyan", "f      "),
		ansigo.MaybeApplyWithReset("cyan", "ctlr+W "),
		ansigo.MaybeApplyWithReset("cyan", "tab    "),
		ansigo.MaybeApplyWithReset("cyan", "p      "),
		ansigo.MaybeApplyWithReset("cyan", "P      "),
		ansigo.MaybeApplyWithReset("cyan", "m      "),
		ansigo.MaybeApplyWithReset("cyan", "M      "),
		ansigo.MaybeApplyWithReset("cyan", "u      "),
		ansigo.MaybeApplyWithReset("cyan", "d      "),

		ansigo.MaybeApplyWithReset("cyan", "down "),
		ansigo.MaybeApplyWithReset("cyan", "up   "),
		ansigo.MaybeApplyWithReset("cyan", "right"),
		ansigo.MaybeApplyWithReset("cyan", "left "),

		ansigo.MaybeApplyWithReset("underline", "More information"),
		ansigo.MaybeApplyWithReset("italic+6", "https://mandelnote.projects.makyo.io"),
		ansigo.MaybeApplyWithReset("italic+6", "https://github.com/makyo/mandelnote"),
		ansigo.MaybeApplyWithReset("underline", "Contributors"),
		ansigo.MaybeApplyWithReset("italic+6", "https://makyo.is"))
)

func (t *tui) showHelp(g *gotui.Gui, v *gotui.View) error {
	if t.editorOpen {
		g.CurrentView().EditWrite('?')
		return nil
	}
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
	if err := g.DeleteView("modalHelp"); err != nil {
		return err
	}
	g.Cursor = false
	t.modalOpen = false
	return nil
}

func (t *tui) createModal(title, content string) {
	if t.modalOpen {
		return
	}
	t.g.Update(func(g *gotui.Gui) error {
		maxX, maxY := g.Size()
		if v, err := g.SetView("modal", 3, 3, maxX-4, maxY-4); err != nil {
			if err != gotui.ErrUnknownView {
				return err
			}
			t.modalOpen = true
			v.Frame = true
			v.FrameFgColor = gotui.ColorCyan | gotui.AttrBold
			v.TitleFgColor = gotui.AttrBold
			v.Title = fmt.Sprintf(" %s ", title)
			v.Wrap = true
			v.WordWrap = true
			fmt.Fprint(v, content)
			if _, err := g.SetViewOnTop("modal"); err != nil {
				return err
			}
		}
		modalHelpText := " Scroll: ↑/↓ | Close: Q "
		if v, err := g.SetView("modalHelp", maxX-3-len(modalHelpText), maxY-5, maxX-6, maxY-3); err != nil {
			if err != gotui.ErrUnknownView {
				return err
			}
			v.Frame = false
			fmt.Fprint(v, ansigo.MaybeApplyWithReset("bold", modalHelpText))
		}
		if _, err := g.SetCurrentView("modal"); err != nil {
			return err
		}
		return nil
	})
}

func confirmNoop(g *gotui.Gui) error {
	return nil
}

func (t *tui) confirmYes(g *gotui.Gui, v *gotui.View) error {
	if t.editorOpen {
		g.CurrentView().EditWrite('y')
	}
	if err := g.DeleteView("confirm"); err != nil {
		return err
	}
	if err := g.DeleteView("confirmActions"); err != nil {
		return err
	}
	if err := t.confirmYesFn(g); err != nil {
		return err
	}
	t.confirmYesFn = confirmNoop
	t.confirmNoFn = confirmNoop
	return nil
}

func (t *tui) confirmNo(g *gotui.Gui, v *gotui.View) error {
	if t.editorOpen {
		g.CurrentView().EditWrite('n')
	}
	if err := g.DeleteView("confirm"); err != nil {
		return err
	}
	if err := g.DeleteView("confirmActions"); err != nil {
		return err
	}
	if err := t.confirmNoFn(g); err != nil {
		return err
	}
	t.confirmYesFn = confirmNoop
	t.confirmNoFn = confirmNoop
	return nil
}
