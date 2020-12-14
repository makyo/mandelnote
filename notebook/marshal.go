package notebook

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Marshal generates a Markdown string of a card and all its children, with its title in a header.
func (c *card) Marshal(depth int) string {
	body := ""
	curr := c
	for curr != nil {
		body += fmt.Sprintf("\n%s %s\n\n%s\n", strings.Repeat("#", depth), curr.title, curr.body)
		if curr.firstChild != nil {
			body += curr.firstChild.Marshal(depth + 1)
		}
		curr = curr.next
	}
	return body
}

// MarshalHeader generates a yaml block of the notebook's metadata
func (nb *Notebook) MarshalHeader() ([]byte, error) {
	return yaml.Marshal(nb)
}

func (nb *Notebook) MarshalBody() string {
	return nb.root.firstChild.Marshal(1)
}

func (nb *Notebook) Marshal() string {
	header, _ := nb.MarshalHeader()
	return fmt.Sprintf("---\n%s\n---\n%s", header, nb.MarshalBody())
}

func Unmarshal(contents string) (*Notebook, error) {
	parts := strings.SplitN(contents, "---\n", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("malformed notebook; must contain metadata block and body")
	}
	nb := New("", "", "", "")
	err := yaml.Unmarshal([]byte(parts[1]), nb)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(parts[2], "\n")
	haveValidFirst := false
	currDepth := 1
	for _, line := range lines {
		currTitle, currBody := nb.GetCard()
		if line == "" && len(currBody) == 0 {
			continue
		}
		if len(line) > 0 && line[0] == '#' {
			if nb.currentCard != nil {
				nb.EditCard(nb.currentCard.title, strings.TrimRight(nb.currentCard.body, "\n"))
			}
			lineParts := strings.SplitN(line, " ", 2)
			if len(lineParts) != 2 {
				return nil, fmt.Errorf("malformed notebook; title must contain depth marker and text")
			}
			depth, title := lineParts[0], lineParts[1]
			if !haveValidFirst && len(depth) > 1 {
				return nil, fmt.Errorf("malformed notebook; must start at header depth 1, found %s", line)
			}
			haveValidFirst = true
			if len(depth) == currDepth {
				nb.AddCard(title, "", false)
			} else if len(depth) == currDepth+1 {
				nb.AddCard(title, "", true)
				currDepth++
			} else if len(depth) < currDepth {
				for currDepth != len(depth) {
					currDepth--
					nb.Exit()
				}
				nb.AddCard(title, "", false)
			} else {
				return nil, fmt.Errorf("malformed notebook; header depths must increase by 1")
			}
		} else {
			if !haveValidFirst {
				return nil, fmt.Errorf("malformed notebook; cannot have body without header")
			}
			if len(currBody) > 0 {
				if len(line) == 0 {
					nb.EditCard(currTitle, currBody+"\n")
				} else {
					nb.EditCard(currTitle, strings.TrimRight(fmt.Sprintf("%s\n%s", currBody, line), "\n"))
				}
			} else {
				nb.EditCard(currTitle, line)
			}
		}
	}
	if nb.currentCard != nil {
		nb.EditCard(nb.currentCard.title, strings.TrimRight(nb.currentCard.body, "\n"))
	}
	nb.currentCard = nb.root.firstChild
	return nb, nil
}

// SetFile changes the file to which the notebook is saved.
func (nb *Notebook) SetFile(filename string) {
	nb.filename = filename
}

// Save saves the notebook's contents to disk.
func (nb *Notebook) Save() error {
	err := ioutil.WriteFile(nb.filename, []byte(nb.Marshal()), 0644)
	if err != nil {
		return err
	}
	nb.dirty = false
	return nil
}

// Open opens a notebook from a file
func Open(filename string) (*Notebook, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		if err.Error() == fmt.Sprintf("open %s: no such file or directory", filename) {
			return New(filename, "New notebook", "", ""), nil
		}
		return nil, err
	}
	nb, err := Unmarshal(string(contents))
	if nb != nil {
		nb.filename = filename
		nb.dirty = false
	}
	return nb, err
}
