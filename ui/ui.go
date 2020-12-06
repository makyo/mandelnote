package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/makyo/ansigo"
	"github.com/makyo/gotui"
	tb "github.com/nsf/termbox-go"

	"github.com/makyo/mandelnote/notebook"
)

var (
	columns  int    = 12
	helpText string = fmt.Sprintf(
		`

		%s

		A snowflake-method writing tool.

		%s

		%s - edit notebook metadata
		%s - save
		%s - quit

		%s - toggle between editing card title and card body
		%s - new card
		%s - new child card
		%s - promote card
		%s - promote all cards at this level

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
		ansigo.MaybeApplyWithReset("cyan", "ctrl+M"),
		ansigo.MaybeApplyWithReset("cyan", "ctrl+S"),
		ansigo.MaybeApplyWithReset("cyan", "ctrl+Q"),
		ansigo.MaybeApplyWithReset("cyan", "shift+tab   "),
		ansigo.MaybeApplyWithReset("cyan", "ctrl+N      "),
		ansigo.MaybeApplyWithReset("cyan", "ctrl+shift+N"),
		ansigo.MaybeApplyWithReset("cyan", "ctrl+P      "),
		ansigo.MaybeApplyWithReset("cyan", "ctrl+shift+P"),
		ansigo.MaybeApplyWithReset("cyan", "ctrl+down "),
		ansigo.MaybeApplyWithReset("cyan", "ctrl+up   "),
		ansigo.MaybeApplyWithReset("cyan", "ctrl+right"),
		ansigo.MaybeApplyWithReset("cyan", "ctrl+left "),
		ansigo.MaybeApplyWithReset("underline", "More information"),
		ansigo.MaybeApplyWithReset("italic+6", "https://mandelnote.projects.makyo.io"),
		ansigo.MaybeApplyWithReset("italic+6", "https://github.com/makyo/mandelnote"),
		ansigo.MaybeApplyWithReset("underline", "Contributors"),
		ansigo.MaybeApplyWithReset("italic+6", "https://makyo.is"),
	)
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

func (t *tui) drawCard(currentCard notebook.Card, g *gotui.Gui, top, left, depth int) (int, error) {
	indent := left + (t.colWidth * depth)
	name := fmt.Sprintf("card-%d", t.cardNameIndex)
	t.cardNameIndex++
	if v, err := g.SetView(name, indent, top+2, indent+(t.cardWidth*t.colWidth), top+10); err != nil {
		if err != gotui.ErrUnknownView {
			return -1, fmt.Errorf("couldn't create card view: %v", err)
		}
		c := &card{
			card:    currentCard,
			view:    v,
			name:    name,
			depth:   depth,
			current: currentCard.Current,
		}
		t.cards = append(t.cards, c)
		v.Frame = true
		if !c.current {
			v.FrameFgColor = gotui.Attribute(tb.AttrDim | tb.ColorDarkGray)
			v.TitleFgColor = gotui.Attribute(tb.AttrDim | tb.ColorDarkGray)
			v.FgColor = gotui.Attribute(tb.AttrDim | tb.ColorDarkGray)
		} else {
			if _, err = g.SetCurrentView(name); err != nil {
				return -1, err
			}
			t.currentDepth = depth
		}
		v.Title = fmt.Sprintf(" %s ", c.card.Title)
		fmt.Fprint(v, c.card.Body)

		for _, child := range currentCard.Children {
			newTop, err := t.drawCard(child, g, top+10, indent, depth+1)
			if err != nil {
				return -1, err
			}
			top = newTop
		}
	}
	return top, nil
}

func (t *tui) drawCards(g *gotui.Gui, width int) error {
	for _, c := range t.cards {
		if err := g.DeleteView(c.name); err != nil {
			return err
		}
	}
	t.cardNameIndex = 0
	t.cards = []*card{}
	tree := t.nb.GetTree()
	left := (((columns - t.cardWidth) / 2) * t.colWidth)
	if (columns-t.cardWidth)%2 == 1 {
		left += t.colWidth / 2
	}
	top := 0
	for _, c := range tree {
		newTop, err := t.drawCard(c, g, top, left, 0)
		if err != nil {
			return err
		}
		top = newTop + 10
	}
	for _, c := range t.cards {
		if x1, y1, x2, y2, err := g.ViewPosition(c.name); err != nil {
			return err
		} else {
			if v, err := g.SetView(c.name, x1-t.currentDepth*t.colWidth, y1, x2-t.currentDepth*t.colWidth, y2); err != nil {
				return err
			} else {
				c.view = v
			}
		}
	}
	return nil
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

func (t *tui) createModal(title, content string) {
	if t.modalOpen {
		return
	}
	go t.g.Update(func(g *gotui.Gui) error {
		maxX, maxY := g.Size()
		if v, err := g.SetView("modal", 3, 3, maxX-4, maxY-4); err != nil {
			if err != gotui.ErrUnknownView {
				return err
			}
			t.modalOpen = true
			v.Frame = true
			v.FrameFgColor = gotui.ColorCyan | gotui.AttrBold
			v.Wrap = true
			v.WordWrap = true
			fmt.Fprint(v, content)
		}
		if v, err := g.SetView("modalTitle", 5, 2, len(title)+8, 4); err != nil {
			if err != gotui.ErrUnknownView {
				return err
			}
			v.Frame = false
			fmt.Fprintf(v, " %s ", ansigo.MaybeApplyWithReset("bold", title))
		}
		modalHelpText := " Scroll: ↑/↓ | Close: <Enter> "
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
