package notebook

import (
	"time"
)

type Notebook struct {
	filename        string
	title           string
	author          string
	description     string
	revisionMessages []message
	created         time.Time
	modified        time.Time
	cards           []*Card
	currentCard     *Card
	currentIndex    []int
}

type message struct {
	message string
	created time.Time
}

type Card struct {
	title    string
	body     string
	children []*Card
}

// Addcard adds a card to the notebook
func (nb *Notebook) AddCard(title, body string) {
	c := &card{
		title: title,
		body: body,
		children: []*card{},
	}
	if !nb.currIndex {
		if !nb.cards {
			nb.cards = []*card{c}
		} else {
			nb.currIndex = []int{0}
			nb.AddCard(title, body)
		}
	} else if len(nb.currIndex) == 1 {
		nb.cards = append(nb.cards, &card{
			title: title,
			body: body,
			children: []*card{}
		})
	} else {
		curr = nb.traverse(nb.currIndex[:len(nb.currIndex) -2])
		curr.children = append(curr.children, c)
	}
}

// SetMetadata sets the metadata for the notebook.
func (nb *Notebook) SetMetadata(title, author, description string) error {
	nb.title = title
	nb.author = author
	nb.description = description
}

// AddRevision adds a timestamped revision message to the notebook.
func (nb *Notebook) AddRevision(text string) {
	nb.revisionMessages = append([]message{message{
		message: text,
		created: time.Now(),
	}}, nb.revisionMessages...)
}

// Cycle set the current card by moving through the current card stack, looping around on overflow.
func (nb *Notebook) Cycle(amount int) {
	if len(nb.currIndex) == 0 {
		if len(nb.cards) == 0 {
			return
		}
		nb.currIndex = []int{0}
	}
	stack := nb.cards
	if len(nb.currIndex) > 1 {
		for _, i := range nb.currIndex {
			stack = stack[i].children
		}
	}
	for nb.currIndex[0] + amount >= len(stack) {
		amount = amount - len(stack)
	}
	nb.currIndex[len(nb.currIndex)-1] += amount
	nb.currCard = nb.cards[nb.currIndex[len(nb.currIndex)-1]]
}

// Enter sets the new current card to the first of the children of the previous current card.
func (nb *Notebook) Enter() {
	if nb.currCard.children {
		nb.currCard = nb.currCard.children[0]
		nb.currIndex = append(nb.currIndex, 0)
	}
}

// Exit goes back to the old current card.
func (nb *Notebook) Exit() {
	if len(nb.currIndex) > 1 {
		nb.currIndex = nb.currIndex[:len(nb.currIndex)-2]
		nb.currCard = nb.traverse(nb.currIndex)
	}
}

func (nb *Notebook) traverse(to []int) *Card {
	if len(to) == 1 {
		return nb.cards[to[0]]
	}
	var c Card
	for _, i := to {
		if !card {
			c = nb.cards[i]
		} else {
			c = c.children[i]
		}
	}
	return c
}

// Delete deletes a card. If it has child and force is not set, it refuses.
func (nb *Notebook) Delete(force bool) error {
	if len(nb.currCard.children) > 0 && !force {
		return fmt.Errorf("card still has children, delete requires force")
	}
	if len(nb.currIndex) == 1 {
		nb.cards = append(nb.cards[:nb.currIndex[0]], nb.cards[nb.currIndex[0]+1:]...)
	}
	index := nb.currIndex[len(nb.currIndex) - 2]
	c = nb.traverse(index)
	c.children = append(c.children[:index], c.children[index+1:]...)
	return nil
}

// Promote promotes a child card to the level of its parent.
func (nb *Notebook) Promote() error {
	if len(nb.currIndex) <= 1 {
		return fmt.Errorf("unable to promote any further")
	} else if len(nb.currIndex) == 2 {
		nb.cards = append(nb.cards[:nb.currIndex[0]], nb.currCard, nb.cards[nb.currIndex:]...)
		nb.Delete(true)
	} else {
		index := nb.currIndex[len(nb.currIndex)-2]
		c = nb.traverse(nb.currIndex[len(nb.currIndex)-3])
		c.children = append(c.children[:index], nb.currCard, c.children[index:]...)
		nb.Delete(true)
	}
}

// PromoteAll promotes all child cards to the level of their parent. If replace is true, the child cards replace the parent.
func (nb *Notebook) PromoteAll(replace bool) error {
	plus := 0
	if replace {
		plus = 1
	}
	if len(nb.currIndex) <= 1 {
		return fmt.Errorf("unable to promote any further")
	} else if len(nb.currIndex) == 2 {
		nb.cards = append(nb.cards[:nb.currIndex[0]], nb.cards[nb.currIndex[0]].children..., nb.cards[nb.currIndex[0]+plus:]...)
	} else {
		index := nb.currIndex[len(nb.currIndex)-2]
		c = nb.traverse(nb.currIndex[len(nb.currIndex)-3])
		c.children = append(c.children[:index], nb.currCard, c.children[index:]...)
		nb.Delete(true)
	}
}

// New creates a new notebook
func New(filename, title, author, description string) (*Notebook, error) {
	return &Notebook{
		filename:    filename,
		title:       title,
		author:      author,
		description: description,
		created:     time.Now(),
		cards:       []*Card{},
	}
}
