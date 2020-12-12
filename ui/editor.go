package ui

import (
	"fmt"
	"strings"

	"github.com/makyo/gotui"
	tb "github.com/nsf/termbox-go"
)

func (t *tui) edit(g *gotui.Gui, v *gotui.View) error {
	if t.modalOpen {
		if t.editorOpen {
			g.CurrentView().EditNewLine()
		}
		return nil
	}
	maxX, maxY := g.Size()
	title, body := t.nb.GetCard()
	if eb, err := g.SetView("editor", t.colWidth, t.colWidth/2, maxX-t.colWidth, maxY-t.colWidth/2); err != nil {
		if err != gotui.ErrUnknownView {
			return err
		}
		g.SetViewOnTop("editor")
		if _, err = g.SetCurrentView("editor"); err != nil {
			return err
		}
		eb.Editable = true
		eb.Wrap = true
		eb.WordWrap = true

		fmt.Fprint(eb, body)
	}
	if et, err := g.SetView("editorTitle", t.colWidth, t.colWidth/2-2, maxX-t.colWidth, t.colWidth/2); err != nil {
		if err != gotui.ErrUnknownView {
			return err
		}
		et.Editable = true
		et.FrameFgColor = gotui.Attribute(tb.AttrDim)
		et.FgColor = gotui.Attribute(tb.AttrDim)

		fmt.Fprint(et, title)
	}
	g.Cursor = true
	t.editorOpen = true
	t.modalOpen = true
	return nil
}

func (t *tui) closeEditor(g *gotui.Gui, v *gotui.View) error {
	if !t.editorOpen {
		return nil
	}
	eb, err := g.View("editor")
	if err != nil {
		return err
	}
	et, err := g.View("editorTitle")
	if err != nil {
		return err
	}
	t.nb.EditCard(strings.Join(et.BufferLines(), "\n"), strings.Join(eb.BufferLines(), "\n"))
	if err = g.DeleteView("editor"); err != nil {
		return err
	}
	if err = g.DeleteView("editorTitle"); err != nil {
		return err
	}
	g.Cursor = false
	t.modalOpen = false
	t.editorOpen = false
	t.focusMode = false
	maxX, _ := g.Size()
	t.drawCards(g, maxX)
	return nil
}

func (t *tui) focus(g *gotui.Gui, v *gotui.View) error {
	if t.modalOpen {
		if t.editorOpen {
			g.CurrentView().EditWrite('f')
		}
		return nil
	}
	maxX, maxY := g.Size()
	margin := (maxX - t.cardWidth*t.colWidth) / 2
	title, body := t.nb.GetCard()
	if eb, err := g.SetView("editor", margin, 4, maxX-margin, maxY); err != nil {
		if err != gotui.ErrUnknownView {
			return err
		}
		g.SetViewOnTop("editor")
		if _, err = g.SetCurrentView("editor"); err != nil {
			return err
		}
		eb.Frame = false
		eb.Editable = true
		eb.Wrap = true
		eb.WordWrap = true

		fmt.Fprint(eb, body)
	}
	if et, err := g.SetView("editorTitle", margin, 1, maxX-margin, 3); err != nil {
		if err != gotui.ErrUnknownView {
			return err
		}
		et.Editable = true
		et.FrameFgColor = gotui.Attribute(tb.AttrDim | tb.ColorDarkGray)
		et.FgColor = gotui.Attribute(tb.AttrDim)

		fmt.Fprint(et, title)
	}
	g.Cursor = true
	t.editorOpen = true
	t.modalOpen = true
	t.focusMode = true
	go g.Update(func(gg *gotui.Gui) error {
		return t.clearCards(gg)
	})
	return nil
}

func (t *tui) toggleEdit(g *gotui.Gui, v *gotui.View) error {
	if !t.editorOpen {
		return nil
	}
	v.FrameFgColor = gotui.Attribute(tb.AttrDim)
	v.FgColor = gotui.Attribute(tb.AttrDim)
	if v.Name() == "editor" {
		v, err := g.SetCurrentView("editorTitle")
		if err != nil {
			return err
		}
		v.FrameFgColor = gotui.ColorDefault
		v.FgColor = gotui.ColorDefault
	} else {
		v, err := g.SetCurrentView("editor")
		if err != nil {
			return err
		}
		v.FrameFgColor = gotui.ColorDefault
		v.FgColor = gotui.ColorDefault
	}
	return nil
}
