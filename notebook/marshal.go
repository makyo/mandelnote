package notebook

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Marshal generates a Markdown string of a card and all its children, with its title in a header.
func (c *Card) Marshal(depth int) string {
	body := fmt.Sprintf("\n%s %s\n\n%s\n", strings.Repeat("#", depth), c.title, c.body)

	for _, child := range c.children {
		body = body + child.Marshal(depth+1)
	}

	return body
}

// MarshalHeader generates a yaml block of the notebook's metadata
func (nb *Notebook) MarshalHeader() ([]byte, error) {
	return yaml.Marshal(nb)
}

func (nb *Notebook) MarshalBody() string {
	body := ""

	for _, card := range nb.cards {
		body = body + card.Marshal(1)
	}

	return body
}

func (nb *Notebook) Marshal() ([]byte, error) {
	return []byte{}, nil
}
