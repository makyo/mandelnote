package ui

import (
	"fmt"

	"github.com/makyo/gotui"
	tb "github.com/nsf/termbox-go"

	"github.com/makyo/mandelnote/notebook"
)

type card struct {
	card    notebook.Card
	view    *gotui.View
	name    string
	depth   int
	top     int
	current bool
}

func (t *tui) newCard(g *gotui.Gui, v *gotui.View) error {
	if t.modalOpen {
		return nil
	}
	maxX, _ := g.Size()
	t.nb.AddCard("New Card", " ", false)
	g.Update(func(gg *gotui.Gui) error {
		return t.drawCards(gg, maxX)
	})
	return nil
}

func (t *tui) newChild(g *gotui.Gui, v *gotui.View) error {
	if t.modalOpen {
		return nil
	}
	maxX, _ := g.Size()
	t.nb.AddCard("New Card", " ", true)
	g.Update(func(gg *gotui.Gui) error {
		return t.drawCards(gg, maxX)
	})
	return nil
}

func (t *tui) cycleUp(g *gotui.Gui, v *gotui.View) error {
	if t.modalOpen {
		return nil
	}
	maxX, _ := g.Size()
	t.nb.Cycle(-1)
	g.Update(func(gg *gotui.Gui) error {
		return t.drawCards(gg, maxX)
	})
	return nil
}

func (t *tui) cycleDown(g *gotui.Gui, v *gotui.View) error {
	if t.modalOpen {
		return nil
	}
	maxX, _ := g.Size()
	t.nb.Cycle(1)
	g.Update(func(gg *gotui.Gui) error {
		return t.drawCards(gg, maxX)
	})
	return nil
}

func (t *tui) enter(g *gotui.Gui, v *gotui.View) error {
	if t.modalOpen {
		return nil
	}
	maxX, _ := g.Size()
	t.nb.Enter()
	g.Update(func(gg *gotui.Gui) error {
		return t.drawCards(gg, maxX)
	})
	return nil
}

func (t *tui) exit(g *gotui.Gui, v *gotui.View) error {
	if t.modalOpen {
		return nil
	}
	maxX, _ := g.Size()
	t.nb.Exit()
	g.Update(func(gg *gotui.Gui) error {
		return t.drawCards(gg, maxX)
	})
	return nil
}

func (t *tui) drawCard(currentCard notebook.Card, g *gotui.Gui, height, top, left, depth int) (int, error) {
	indent := left + (t.colWidth * depth)
	name := fmt.Sprintf("card-%d", t.cardNameIndex)
	t.cardNameIndex++
	x1 := indent
	y1 := top + 2
	x2 := indent + (t.cardWidth * t.colWidth)
	y2 := top + 10
	var rows int
	if v, err := g.SetView(name, x1, y1, x2, y2); err != nil {
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
		v.Wrap = true
		v.WordWrap = true
		if !c.current {
			v.FrameFgColor = gotui.Attribute(tb.AttrDim | tb.ColorDarkGray)
			v.TitleFgColor = gotui.Attribute(tb.AttrDim | tb.ColorDarkGray)
			v.FgColor = gotui.Attribute(tb.AttrDim | tb.ColorDarkGray)
		} else {
			if !t.modalOpen {
				if _, err = g.SetCurrentView(name); err != nil {
					return -1, err
				}
			}
			t.currentDepth = depth
			t.currentY = top
		}
		v.Title = fmt.Sprintf(" %s ", c.card.Title)
		fmt.Fprint(v, c.card.Body)

		rows = len(v.BufferLines())
		if rows > height/3 {
			rows = height / 3
		}
		if c.current {
			t.currentHeight = rows
		}
		if _, err := g.SetView(name, x1, y1, x2, y1+rows+1); err != nil {
			return -1, err
		}

		if _, err := g.SetViewOnBottom(name); err != nil {
			return -1, err
		}

		for _, child := range currentCard.Children {
			newTop, err := t.drawCard(child, g, height, top+rows+2, indent, depth+1)
			if err != nil {
				return -1, err
			}
			top = newTop + len(t.cards[len(t.cards)-1].view.ViewBufferLines())
		}
	}
	return top + rows, nil
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
	_, maxY := g.Size()
	for _, c := range tree {
		newTop, err := t.drawCard(c, g, maxY, top, left, 0)
		if err != nil {
			return err
		}
		top = newTop + len(t.cards[len(t.cards)-1].view.ViewBufferLines()) + 2
	}
	offset := 0
	if t.currentY > maxY/2-t.currentHeight-3 {
		offset = t.currentY + (t.currentHeight / 2)
	}
	for _, c := range t.cards {
		if x1, y1, x2, y2, err := g.ViewPosition(c.name); err != nil {
			return err
		} else {
			if v, err := g.SetView(
				c.name, x1-t.currentDepth*t.colWidth,
				y1-offset+maxY/2,
				x2-t.currentDepth*t.colWidth,
				y2-offset+maxY/2); err != nil {
				return err
			} else {
				c.view = v
			}
		}
	}
	return nil
}
