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
	root         *Card
	currentCard  *Card
	currentIndex []int
}

type Revision struct {
	Message   string
	Timestamp time.Time
}

type Card struct {
	title      string
	body       string
	parent     *Card
	next       *Card
	prev       *Card
	firstChild *Card
}

type cardContents struct {
	Title    string
	Body     string
	Children []cardContents
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
		title: title,
		body:  body,
	}
	if child || nb.currentCard == nb.root {
		c.parent = nb.currentCard
		currChild := nb.currentCard.firstChild
		if currChild == nil {
			nb.currentCard.firstChild = c
		} else {
			for currChild.next != nil {
				currChild = currChild.next
			}
			currChild.next = c
			c.prev = currChild
		}
	} else {
		c.parent = nb.currentCard.parent
		c.next = nb.currentCard.next
		nb.currentCard.next = c
		c.prev = nb.currentCard
	}
	nb.currentCard = c
}

// GetCard returns the contents of the current card
func (nb *Notebook) GetCard() (string, string) {
	return nb.currentCard.title, nb.currentCard.body
}

func (c *Card) getTree() []cardContents {
	result := []cardContents{}
	currChild := c.firstChild
	for currChild != nil {
		result = append(result, cardContents{
			Title:    currChild.title,
			Body:     currChild.body,
			Children: currChild.getTree(),
		})
		currChild = currChild.next
	}
	return result
}

// GetTree returns a tree-representation of the contents of all cards
func (nb *Notebook) GetTree() []cardContents {
	return nb.root.getTree()
}

// EditCard changes the contents of the current card.
func (nb *Notebook) EditCard(title, body string) {
	if nb.currentCard == nb.root {
		return
	}
	nb.currentCard.title = title
	nb.currentCard.body = body
}

// Cycle set the current card by moving through the current card stack, looping around on overflow.
func (nb *Notebook) Cycle(amount int) {
	diff := 0
	if amount > 0 {
		diff = -1
	} else {
		diff = 1
	}
	for amount != 0 {
		if nb.currentCard.next != nil {
			nb.currentCard = nb.currentCard.next
		} else {
			nb.currentCard = nb.currentCard.parent.firstChild
		}
		amount += diff
	}
}

// Enter sets the new current card to the first of the children of the previous current card.
func (nb *Notebook) Enter() {
	if nb.currentCard.firstChild != nil {
		nb.currentCard = nb.currentCard.firstChild
	}
}

// Exit goes back to the old current card.
func (nb *Notebook) Exit() {
	if nb.currentCard != nb.root && nb.currentCard.parent != nb.root {
		nb.currentCard = nb.currentCard.parent
	}
}

// Delete deletes a card. If it has child and force is not set, it refuses.
func (nb *Notebook) Delete(force bool) error {
	if nb.currentCard == nb.root {
		return fmt.Errorf("nothing to delete")
	}
	if nb.currentCard.firstChild != nil && !force {
		return fmt.Errorf("card still has children, delete requires force")
	}
	if nb.currentCard.prev == nil {
		nb.currentCard.parent.firstChild = nb.currentCard.next
		if nb.currentCard.next == nil {
			nb.currentCard = nb.currentCard.parent
		} else {
			nb.currentCard = nb.currentCard.next
			nb.currentCard.prev = nil
		}
	} else {
		nb.currentCard.prev.next = nb.currentCard.next
		nb.currentCard = nb.currentCard.prev
	}
	return nil
}

// Promote promotes a child card to the level of its parent.
func (nb *Notebook) Promote() error {
	if nb.currentCard == nb.root || nb.currentCard.parent == nb.root {
		return fmt.Errorf("unable to promote any further")
	}
	parent := nb.currentCard.parent
	if nb.currentCard.prev != nil {
		nb.currentCard.prev.next = nb.currentCard.next
	}
	if nb.currentCard == parent.firstChild {
		parent.firstChild = nb.currentCard.next
	}
	nb.currentCard.prev = parent
	nb.currentCard.parent = parent.parent
	nb.currentCard.next = parent.next
	parent.next = nb.currentCard
	return nil
}

// PromoteAll promotes all child cards to the level of their parent. If replace is true, the child cards replace the parent.
func (nb *Notebook) PromoteAll(replace bool) error {
	if nb.currentCard == nb.root || nb.currentCard.parent == nb.root {
		return fmt.Errorf("unable to promote any further")
	}
	parent := nb.currentCard.parent
	first := parent.firstChild
	first.prev = parent
	curr := first
	for {
		curr.parent = parent.parent
		if curr.next == nil {
			curr.next = parent.next
			break
		}
		curr = curr.next
	}
	parent.next = first
	parent.firstChild = nil
	nb.currentCard = first
	if replace {
		first.prev = parent.prev
		if first.prev != nil {
			first.prev.next = first
		}
	}
	return nil
}

// New creates a new notebook
func New(filename, title, author, description string) *Notebook {
	root := &Card{}
	return &Notebook{
		filename:    filename,
		Title:       title,
		Author:      author,
		Description: description,
		Created:     time.Now(),
		Modified:    time.Now(),
		Revisions:   []Revision{},
		root:        root,
		currentCard: root,
	}
}
