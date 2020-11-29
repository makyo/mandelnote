package notebook

import (
	"fmt"
	"time"
)

type Notebook struct {
	filename     string
	Title        string
	Author       string
	Description  string
	Revisions    []Revision
	Created      time.Time
	Modified     time.Time
	cards        []*Card
	currentCard  *Card
	currentIndex []int
}

type Revision struct {
	Message   string
	Timestamp time.Time
}

type Card struct {
	title    string
	body     string
	children []*Card
}

// SetMetadata sets the metadata for the notebook.
func (nb *Notebook) SetMetadata(title, author, description string) {
	nb.Title = title
	nb.Author = author
	nb.Description = description
}

// AddRevision adds a timestamped revision message to the notebook.
func (nb *Notebook) AddRevision(text string) {
	nb.Revisions = append([]Revision{Revision{
		Message:   text,
		Timestamp: time.Now(),
	}}, nb.Revisions...)
	nb.Modified = time.Now()
}

// AddCard adds a card to the notebook
func (nb *Notebook) AddCard(title, body string, child bool) {
	c := &Card{
		title:    title,
		body:     body,
		children: []*Card{},
	}
	if len(nb.cards) == 0 {
		nb.cards = []*Card{}
	}
	if len(nb.currentIndex) == 0 {
		nb.currentIndex = []int{0}
		nb.cards = []*Card{c}
		return
	}
	if child {
		nb.currentCard.children = append(nb.currentCard.children, c)
		nb.currentIndex = append(nb.currentIndex, len(nb.currentCard.children)-1)
		nb.currentCard = nb.traverse(nb.currentIndex)
		return
	} else if len(nb.currentIndex) == 1 {
		after := nb.cards[:nb.currentIndex[0]]
		nb.cards = append(nb.cards[nb.currentIndex[0]:], &Card{
			title:    title,
			body:     body,
			children: []*Card{},
		})
		nb.cards = append(nb.cards, after...)
	} else {
		curr := nb.traverse(nb.currentIndex[:len(nb.currentIndex)-1])
		index := nb.currentIndex[len(nb.currentIndex)-1] + 1
		curr.children = append(
			curr.children[:index], append(
				[]*Card{c},
				curr.children[index:]...)...)
	}
	nb.currentIndex[len(nb.currentIndex)-1]++
	nb.currentCard = nb.traverse(nb.currentIndex)
}

func (nb *Notebook) GetCard() (string, string) {
	return nb.currentCard.title, nb.currentCard.body
}

func (nb *Notebook) EditCard(title, body string) {
	nb.currentCard.title = title
	nb.currentCard.body = body
}

// Cycle set the current card by moving through the current card stack, looping around on overflow.
func (nb *Notebook) Cycle(amount int) {
	if len(nb.currentIndex) == 0 {
		if len(nb.cards) == 0 {
			return
		}
		nb.currentIndex = []int{0}
	}
	var length int
	if len(nb.currentIndex) == 1 {
		length = len(nb.cards)
	} else {
		curr := nb.traverse(nb.currentIndex[:len(nb.currentIndex)-1])
		length = len(curr.children)
	}
	newIndex := (nb.currentIndex[len(nb.currentIndex)-1] + amount) % length
	if newIndex < 0 {
		newIndex = newIndex * -1
	}
	nb.currentIndex[len(nb.currentIndex)-1] = newIndex
	nb.currentCard = nb.traverse(nb.currentIndex)
}

// Enter sets the new current card to the first of the children of the previous current card.
func (nb *Notebook) Enter() {
	if len(nb.currentCard.children) != 0 {
		nb.currentCard = nb.currentCard.children[0]
		nb.currentIndex = append(nb.currentIndex, 0)
	}
}

// Exit goes back to the old current card.
func (nb *Notebook) Exit() {
	if len(nb.currentIndex) > 1 {
		nb.currentIndex = nb.currentIndex[:len(nb.currentIndex)-1]
		nb.currentCard = nb.traverse(nb.currentIndex)
	}
}

func (nb *Notebook) traverse(to []int) *Card {
	if len(to) == 1 {
		return nb.cards[to[0]]
	}
	var c *Card
	for _, i := range to {
		if c == nil {
			c = nb.cards[i]
		} else {
			c = c.children[i]
		}
	}
	return c
}

// Delete deletes a card. If it has child and force is not set, it refuses.
func (nb *Notebook) Delete(force bool) error {
	if len(nb.currentCard.children) > 0 && !force {
		return fmt.Errorf("card still has children, delete requires force")
	}
	if len(nb.currentIndex) == 1 {
		nb.cards = append(
			nb.cards[:nb.currentIndex[0]],
			nb.cards[nb.currentIndex[0]+1:]...)
	} else {
		c := nb.traverse(nb.currentIndex[:len(nb.currentIndex)-1])
		c.children = append(
			c.children[:nb.currentIndex[len(nb.currentIndex)-1]],
			c.children[nb.currentIndex[len(nb.currentIndex)-1]+1:]...)
	}
	return nil
}

// Promote promotes a child card to the level of its parent.
func (nb *Notebook) Promote() error {
	if len(nb.currentIndex) <= 1 {
		return fmt.Errorf("unable to promote any further")
	} else if len(nb.currentIndex) == 2 {
		cards := append(
			nb.cards[:nb.currentIndex[0]+1],
			nb.currentCard)
		nb.cards = append(
			cards,
			nb.cards[nb.currentIndex[0]+1:]...)
		nb.Delete(true)
	} else {
		index := nb.currentIndex[:len(nb.currentIndex)-2]
		c := nb.traverse(index)
		nb.Delete(true)
		children := append(
			c.children[:nb.currentIndex[len(nb.currentIndex)-1]],
			nb.currentCard)
		c.children = append(
			children,
			c.children[nb.currentIndex[len(nb.currentIndex)-1]:]...)
	}
	nb.currentIndex = nb.currentIndex[:len(nb.currentIndex)-1]
	nb.currentCard = nb.traverse(nb.currentIndex)
	return nil
}

// PromoteAll promotes all child cards to the level of their parent. If replace is true, the child cards replace the parent.
func (nb *Notebook) PromoteAll(replace bool) error {
	plus := 0
	if replace {
		plus = 1
	}
	if len(nb.currentIndex) <= 1 {
		return fmt.Errorf("unable to promote any further")
	} else if len(nb.currentIndex) == 2 {
		cards := append(
			nb.cards[:nb.currentIndex[0]+1],
			nb.currentCard)
		nb.cards = append(
			cards,
			nb.cards[nb.currentIndex[0]+1:]...)
		nb.Delete(true)
	} else {
		index := nb.currentIndex[:len(nb.currentIndex)-1]
		c := nb.traverse(nb.currentIndex)
		children := append(
			c.children[:index[len(index)-1]+1+plus],
			nb.currentCard)
		c.children = append(
			children,
			c.children[index[len(index)-1]+1+plus:]...)
		nb.Delete(true)
	}
	return nil
}

// New creates a new notebook
func New(filename, title, author, description string) *Notebook {
	return &Notebook{
		filename:    filename,
		Title:       title,
		Author:      author,
		Description: description,
		Created:     time.Now(),
		Modified:    time.Now(),
		Revisions:   []Revision{},
		cards:       []*Card{},
	}
}
