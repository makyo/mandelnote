package notebook

import (
	"time"
)

type notebook struct {
	title           string
	author          string
	description     string
	revisionMessage []*message
	created         time.Time
	modified        time.Time
	cards           []*card
}

type message struct {
	message string
	created time.Time
}

type card struct {
	title    string
	body     string
	children []*card
}

// AddCard adds a card to the notebook
func (n *notebook) AddCard(title, body string) error {
}

// SetMetadata sets the metadata for the notebook.
func (n *notebook) SetMetadata(title, author, description string) error {
}

// Delete deletes a card. If it has child and force is not set, it refuses.
func (c *card) Delete(force bool) error {
}

// AddChild adds a child card to the current card.
func (c *card) AddChild(title, body string) error {
}

// Promote promotes a child card to the level of its parent.
func (c *card) Promote() error {
}

// PromoteAll promotes all child cards to the level of their parent. If replace is true, the child cards replace the parent.
func (c *card) PromoteAll(replace bool) error {
}

// New creates a new notebook
func New(title, author, description string) (*notebook, error) {
}
