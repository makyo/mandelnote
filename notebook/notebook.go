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
	root         *card
	currentCard  *card
	currentIndex []int
	dirty        bool
}

type Revision struct {
	Message   string
	Timestamp time.Time
}

type card struct {
	title      string
	body       string
	parent     *card
	next       *card
	prev       *card
	firstChild *card
}

type Card struct {
	Title    string
	Body     string
	Current  bool
	Children []Card
}

// SetMetadata sets the metadata for the notebook.
func (nb *Notebook) SetMetadata(title, author, description string) {
	nb.Title = title
	nb.Author = author
	nb.Description = description
	nb.dirty = true
}

// AddRevision adds a timestamped revision message to the notebook.
func (nb *Notebook) AddRevision(text string) {
	nb.Revisions = append([]Revision{Revision{
		Message:   text,
		Timestamp: time.Now(),
	}}, nb.Revisions...)
	nb.Modified = time.Now()
	nb.dirty = true
}

// AddCard adds a card to the notebook
func (nb *Notebook) AddCard(title, body string, child bool) {
	c := &card{
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
	nb.dirty = true
}

// GetCard returns the contents of the current card
func (nb *Notebook) GetCard() (string, string) {
	return nb.currentCard.title, nb.currentCard.body
}

func (c *card) getTree(nb *Notebook) []Card {
	result := []Card{}
	currChild := c.firstChild
	for currChild != nil {
		result = append(result, Card{
			Title:    currChild.title,
			Body:     currChild.body,
			Current:  currChild == nb.currentCard,
			Children: currChild.getTree(nb),
		})
		currChild = currChild.next
	}
	return result
}

// GetTree returns a tree-representation of the contents of all cards
func (nb *Notebook) GetTree() []Card {
	return nb.root.getTree(nb)
}

// EditCard changes the contents of the current card.
func (nb *Notebook) EditCard(title, body string) {
	if nb.currentCard == nb.root {
		return
	}
	nb.currentCard.title = title
	nb.currentCard.body = body
	nb.dirty = true
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
		var to *card
		if diff == 1 {
			to = nb.currentCard.prev
		} else {
			to = nb.currentCard.next
		}
		if to != nil {
			nb.currentCard = to
		} else {
			if diff == 1 {
				for nb.currentCard.next != nil {
					nb.currentCard = nb.currentCard.next
				}
			} else {
				nb.currentCard = nb.currentCard.parent.firstChild
			}
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
	nb.dirty = true
	return nil
}

// Move will move the current card the number of steps in the given direction. It does not cycle.
func (nb *Notebook) Move(amount int) {
	diff := 0
	if amount > 0 {
		diff = -1
	} else {
		diff = 1
	}
	for amount != 0 {
		amount += diff
		if diff == 1 {
			if nb.currentCard.prev == nil {
				break
			}
			nb.currentCard.prev.next = nb.currentCard.next
			if nb.currentCard.next != nil {
				nb.currentCard.next.prev = nb.currentCard.prev
			}
			nb.currentCard.next = nb.currentCard.prev
			nb.currentCard.prev = nb.currentCard.prev.prev
			nb.currentCard.next.prev = nb.currentCard
			if nb.currentCard.prev != nil {
				nb.currentCard.prev.next = nb.currentCard
			} else {
				nb.currentCard.parent.firstChild = nb.currentCard
			}
		} else {
			if nb.currentCard.next == nil {
				break
			}
			nb.currentCard.next.prev = nb.currentCard.prev
			if nb.currentCard.prev != nil {
				nb.currentCard.prev.next = nb.currentCard.next
			}
			nb.currentCard.prev = nb.currentCard.next
			nb.currentCard.next = nb.currentCard.next.next
			nb.currentCard.prev.next = nb.currentCard
			if nb.currentCard.next != nil {
				nb.currentCard.next.prev = nb.currentCard
			}
			if nb.currentCard.parent.firstChild == nb.currentCard {
				nb.currentCard.parent.firstChild = nb.currentCard.prev
			}
		}
	}
	nb.dirty = true
}

// Merge will merge the bodies and children of cards from the number of cards specified into the current card.
func (nb *Notebook) Merge(amount int) {
	diff := 0
	if amount > 0 {
		diff = -1
	} else {
		diff = 1
	}
	for amount != 0 {
		amount += diff
		if diff == 1 {
			prev := nb.currentCard.prev

			// Can't merge up with no previous card.
			if prev == nil {
				break
			}

			// Append the body.
			prev.body = fmt.Sprintf("%s\n\n%s", prev.body, nb.currentCard.body)

			// Manage next/prev.
			prev.next = nb.currentCard.next
			if prev.next != nil {
				prev.next.prev = prev
			}

			// Manage children.
			if nb.currentCard.firstChild != nil {

				// Set parent for current children.
				child := nb.currentCard.firstChild
				for child != nil {
					child.parent = prev
					child = child.next
				}

				// Append children
				if prev.firstChild != nil {
					oldFirst := nb.currentCard.firstChild
					child = prev.firstChild
					for child.next != nil {
						child = child.next
					}
					child.next = oldFirst
					oldFirst.prev = child
				} else {
					prev.firstChild = nb.currentCard.firstChild
				}
			}
			nb.currentCard = prev
		} else {
			next := nb.currentCard.next

			// Can't merge down with no next card
			if next == nil {
				break
			}

			// Prepend body
			next.body = fmt.Sprintf("%s\n\n%s", nb.currentCard.body, next.body)
			next.title = nb.currentCard.title

			// Manage next/prev
			next.prev = nb.currentCard.prev
			if next.prev != nil {
				next.prev.next = next
			}

			// Manage children
			if nb.currentCard.firstChild != nil {

				// Set parent for current children
				child := nb.currentCard.firstChild
				for child.next != nil {
					child.parent = next
					child = child.next
				}

				// The above will skip the last child
				child.parent = next
				child.next = next.firstChild

				if next.firstChild != nil {
					next.firstChild.prev = child
				}
				next.firstChild = nb.currentCard.firstChild
			}

			// If merging the first child down, set the parent's first child.
			if nb.currentCard.parent.firstChild == nb.currentCard {
				nb.currentCard.parent.firstChild = next
			}
			nb.currentCard = next
		}
	}
	nb.dirty = true
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
	nb.dirty = true
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
	nb.dirty = true
	return nil
}

// Dirty returns whether or not the notebook has changes that have not been saved.
func (nb *Notebook) Dirty() bool {
	return nb.dirty
}

// New creates a new notebook
func New(filename, title, author, description string) *Notebook {
	root := &card{}
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
