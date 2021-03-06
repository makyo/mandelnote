package notebook_test

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	yaml "gopkg.in/yaml.v2"

	"github.com/makyo/mandelnote/notebook"
)

func TestNotebook(t *testing.T) {
	Convey("When working with a notebook", t, func() {

		nb := notebook.New(
			"test.md",
			"Test notebook",
			"Test author",
			"This is a notebook for testing")

		Convey("It can have metadata", func() {

			Convey("One can set the metadata", func() {
				nb.SetMetadata("New title", "New author", "New description")
				out, err := nb.MarshalHeader()
				So(err, ShouldBeNil)

				// Bit of a round-about way to do it, but I'm not testing all that text.
				nb2 := &notebook.Notebook{}
				err = yaml.Unmarshal(out, nb2)
				So(err, ShouldBeNil)
				So(nb2.Title, ShouldEqual, "New title")
				So(nb2.Author, ShouldEqual, "New author")
				So(nb2.Description, ShouldEqual, "New description")
				So(nb2.Modified.After(nb.Created), ShouldBeTrue)
			})

			Convey("One can add revision info", func() {
				nb.AddRevision("Made some changes")
				out, err := nb.MarshalHeader()
				So(err, ShouldBeNil)

				nb2 := &notebook.Notebook{}
				err = yaml.Unmarshal(out, nb2)
				So(err, ShouldBeNil)
				So(len(nb2.Revisions), ShouldEqual, 1)
				So(nb2.Revisions[0].Message, ShouldEqual, "Made some changes")
				So(nb2.Revisions[0].Timestamp.After(nb.Created), ShouldBeTrue)
				So(nb2.Modified.After(nb.Created), ShouldBeTrue)
			})
		})

		Convey("It can have cards", func() {

			nb.AddCard("Card 1 Title", "Card 1 body", false)
			title, body := nb.GetCard()
			So(title, ShouldEqual, "Card 1 Title")
			So(body, ShouldEqual, "Card 1 body")
			So(nb.GetTree(), ShouldHaveLength, 1)

			Convey("Adding cards makes the notebook dirty", func() {
				So(nb.Dirty(), ShouldBeTrue)
			})

			nb.AddCard("Card 2 Title", "Card 2 body", false)
			title, body = nb.GetCard()
			So(title, ShouldEqual, "Card 2 Title")
			So(body, ShouldEqual, "Card 2 body")
			So(nb.GetTree(), ShouldHaveLength, 2)

			Convey("It can move them", func() {

				nb = notebook.New("", "Moving", "", "")
				nb.AddCard("Test Up", "up", false)
				nb.AddCard("Test Down", "down", false)
				nb.AddCard("Bottom text", "bottom text", false)
				nb.Cycle(-1)

				Convey("Up", func() {
					tree := nb.GetTree()
					So(tree[0].Title, ShouldEqual, "Test Up")
					So(tree[1].Title, ShouldEqual, "Test Down")
					So(tree[1].Current, ShouldEqual, true)
					So(tree[2].Title, ShouldEqual, "Bottom text")
					nb.Move(-1)
					tree = nb.GetTree()
					So(tree[0].Title, ShouldEqual, "Test Down")
					So(tree[0].Current, ShouldEqual, true)
					So(tree[1].Title, ShouldEqual, "Test Up")
					So(tree[2].Title, ShouldEqual, "Bottom text")
					nb.Move(-1)
					tree = nb.GetTree()
					So(tree[0].Title, ShouldEqual, "Test Down")
					So(tree[0].Current, ShouldEqual, true)
					So(tree[1].Title, ShouldEqual, "Test Up")
					So(tree[2].Title, ShouldEqual, "Bottom text")
					nb.Cycle(2)
					nb.Move(-1)
					tree = nb.GetTree()
					So(tree[0].Title, ShouldEqual, "Test Down")
					So(tree[1].Title, ShouldEqual, "Bottom text")
					So(tree[1].Current, ShouldEqual, true)
					So(tree[2].Title, ShouldEqual, "Test Up")
				})

				Convey("Down", func() {
					tree := nb.GetTree()
					So(tree[0].Title, ShouldEqual, "Test Up")
					So(tree[1].Title, ShouldEqual, "Test Down")
					So(tree[1].Current, ShouldEqual, true)
					So(tree[2].Title, ShouldEqual, "Bottom text")
					nb.Move(1)
					tree = nb.GetTree()
					So(tree[0].Title, ShouldEqual, "Test Up")
					So(tree[1].Title, ShouldEqual, "Bottom text")
					So(tree[2].Title, ShouldEqual, "Test Down")
					So(tree[2].Current, ShouldEqual, true)
					nb.Move(1)
					tree = nb.GetTree()
					So(tree[0].Title, ShouldEqual, "Test Up")
					So(tree[1].Title, ShouldEqual, "Bottom text")
					So(tree[2].Title, ShouldEqual, "Test Down")
					So(tree[2].Current, ShouldEqual, true)
					nb.Cycle(-2)
					nb.Move(1)
					tree = nb.GetTree()
					So(tree[0].Title, ShouldEqual, "Bottom text")
					So(tree[1].Title, ShouldEqual, "Test Up")
					So(tree[1].Current, ShouldEqual, true)
					So(tree[2].Title, ShouldEqual, "Test Down")
				})
			})

			Convey("It can merge them", func() {

				nb = notebook.New("", "Merging", "", "")
				nb.AddCard("Test Up", "up", false)
				nb.AddCard("Up 1", "1", true)
				nb.AddCard("Up 2", "2", false)
				nb.Exit()
				nb.AddCard("Test Down", "down", false)
				nb.AddCard("Bottom text", "bottom text", false)
				nb.AddCard("Bottom 1", "1", true)
				nb.AddCard("Bottom 2", "2", false)
				nb.Exit()
				nb.Cycle(-1)

				Convey("Up", func() {
					nb.Merge(-1)
					title, body = nb.GetCard()
					So(title, ShouldEqual, "Test Up")
					So(body, ShouldEqual, "up\n\ndown")

					nb.Cycle(1)
					nb.Merge(-1)
					title, body = nb.GetCard()
					So(title, ShouldEqual, "Test Up")
					So(body, ShouldEqual, "up\n\ndown\n\nbottom text")

					Convey("Including children", func() {
						nb.Enter()
						title, _ = nb.GetCard()
						So(title, ShouldEqual, "Up 1")
						nb.Cycle(2)
						title, _ = nb.GetCard()
						So(title, ShouldEqual, "Bottom 1")
					})
				})

				Convey("And down", func() {
					nb.Merge(1)
					title, body = nb.GetCard()
					So(title, ShouldEqual, "Test Down")
					So(body, ShouldEqual, "down\n\nbottom text")

					nb.Cycle(-1)
					nb.Merge(1)
					title, body = nb.GetCard()
					So(title, ShouldEqual, "Test Up")
					So(body, ShouldEqual, "up\n\ndown\n\nbottom text")

					Convey("Including children", func() {
						nb.Enter()
						title, _ = nb.GetCard()
						So(title, ShouldEqual, "Up 1")
						nb.Cycle(2)
						title, _ = nb.GetCard()
						So(title, ShouldEqual, "Bottom 1")
					})
				})

				Convey("More than one at a time", func() {
					nb.Cycle(1)
					nb.Merge(-2)
					title, body = nb.GetCard()
					So(title, ShouldEqual, "Test Up")
					So(body, ShouldEqual, "up\n\ndown\n\nbottom text")
				})

				Convey("No-ops", func() {
					nb.Cycle(1)
					before := nb.GetTree()
					nb.Merge(1)
					So(nb.GetTree(), ShouldResemble, before)
					nb.Cycle(-2)
					nb.Merge(-1)
					// Needed for Current to be the same.
					nb.Cycle(2)
					So(nb.GetTree(), ShouldResemble, before)
				})
			})

			Convey("And edit them", func() {
				nb.EditCard("Rose", "Companion")
				title, body = nb.GetCard()
				So(title, ShouldEqual, "Rose")
				So(body, ShouldEqual, "Companion")

				Convey("Editing the root is a no-op", func() {
					nb2 := notebook.New("empty.md", "Empty", "Empty", "Empty")
					nb2.EditCard("bad", "wolf")
					So(nb2.GetTree(), ShouldHaveLength, 0)
				})
			})

			Convey("And delete them", func() {
				err := nb.Delete(false)
				So(err, ShouldBeNil)

				body = nb.MarshalBody()
				So(body, ShouldEqual, "\n# Card 1 Title\n\nCard 1 body\n")

				nb.AddCard("bad", "wolf", true)
				err = nb.Delete(false)
				So(err, ShouldBeNil)

				body = nb.MarshalBody()
				So(body, ShouldEqual, "\n# Card 1 Title\n\nCard 1 body\n")

				nb.AddCard("bad", "wolf", true)
				nb.AddCard("good", "wolf", false)
				nb.Cycle(-1)
				err = nb.Delete(false)
				So(err, ShouldBeNil)

				body = nb.MarshalBody()
				So(body, ShouldEqual, "\n# Card 1 Title\n\nCard 1 body\n\n## good\n\nwolf\n")

				nb.Exit()
				err = nb.Delete(false)
				So(err.Error(), ShouldEqual, "card still has children, delete requires force")

				nb2 := notebook.New("empty.md", "Empty", "Empty", "Empty")
				err = nb2.Delete(false)
				So(err.Error(), ShouldEqual, "nothing to delete")
			})

			Convey("Cards can have children", func() {
				nb.AddCard("Card 2.1 Title", "Card 2.1 body", true)

				body = nb.MarshalBody()
				So(body, ShouldEqual, "\n# Card 1 Title\n\nCard 1 body\n\n# Card 2 Title\n\nCard 2 body\n\n## Card 2.1 Title\n\nCard 2.1 body\n")

				nb.Exit()
				nb.AddCard("Card 2.1-a Title", "Card 2.1-a body", true)
				nb.Exit()
				nb.AddCard("Card 2.1-b Title", "Card 2.1-b body", true)

				body = nb.MarshalBody()
				So(body, ShouldEqual, "\n# Card 1 Title\n\nCard 1 body\n\n# Card 2 Title\n\nCard 2 body\n\n## Card 2.1 Title\n\nCard 2.1 body\n\n## Card 2.1-a Title\n\nCard 2.1-a body\n\n## Card 2.1-b Title\n\nCard 2.1-b body\n")
				nb.Delete(true)
				nb.Delete(true)

				Convey("One can move between cards", func() {
					title, body := nb.GetCard()
					So(title, ShouldEqual, "Card 2.1 Title")
					So(body, ShouldEqual, "Card 2.1 body")

					Convey("One can enter and exit", func() {
						nb.Exit()
						title, body = nb.GetCard()
						So(title, ShouldEqual, "Card 2 Title")
						So(body, ShouldEqual, "Card 2 body")

						nb.Enter()
						title, body = nb.GetCard()
						So(title, ShouldEqual, "Card 2.1 Title")
						So(body, ShouldEqual, "Card 2.1 body")
					})
					nb.Exit()

					Convey("One can move to siblings", func() {
						nb.Cycle(-1)
						title, body = nb.GetCard()
						So(title, ShouldEqual, "Card 1 Title")
						So(body, ShouldEqual, "Card 1 body")

						nb.Cycle(1)
						title, body = nb.GetCard()
						So(title, ShouldEqual, "Card 2 Title")
						So(body, ShouldEqual, "Card 2 body")

						Convey("One moves in a cycle", func() {
							title, body = nb.GetCard()
							So(title, ShouldEqual, "Card 2 Title")
							So(body, ShouldEqual, "Card 2 body")
							nb.AddCard("Card 3 Title", "Card 3 body", false)
							title, body = nb.GetCard()
							So(title, ShouldEqual, "Card 3 Title")
							So(body, ShouldEqual, "Card 3 body")
							nb.Cycle(1)
							title, body = nb.GetCard()
							So(title, ShouldEqual, "Card 1 Title")
							So(body, ShouldEqual, "Card 1 body")
							nb.Cycle(1)
							title, body = nb.GetCard()
							So(title, ShouldEqual, "Card 2 Title")
							So(body, ShouldEqual, "Card 2 body")
							nb.Cycle(-1)
							title, body = nb.GetCard()
							So(title, ShouldEqual, "Card 1 Title")
							So(body, ShouldEqual, "Card 1 body")
							nb.Cycle(-1)
							title, body = nb.GetCard()
							So(title, ShouldEqual, "Card 3 Title")
							So(body, ShouldEqual, "Card 3 body")
						})
					})
				})

				Convey("Children can be promoted", func() {
					nb.AddCard("Card 2.2 Title", "Card 2.2 body", false)
					nb.AddCard("Card 2.3 Title", "Card 2.3 body", false)
					body = nb.MarshalBody()
					So(body, ShouldContainSubstring, "\n## Card 2.3 Title")
					title, body := nb.GetCard()
					So(title, ShouldContainSubstring, "Card 2.3 Title")
					err := nb.Promote()
					So(err, ShouldBeNil)
					body = nb.MarshalBody()
					So(body, ShouldContainSubstring, "\n# Card 2.3 Title")
					So(body, ShouldNotContainSubstring, "\n## Card 2.3 Title")

					nb.Cycle(-1)
					nb.Enter()
					title, body = nb.GetCard()
					So(title, ShouldEqual, "Card 2.1 Title")
					nb.AddCard("Card 2.1.1 Title", "Card 2.1.1 body", true)
					title, body = nb.GetCard()
					So(title, ShouldEqual, "Card 2.1.1 Title")
					body = nb.MarshalBody()
					So(body, ShouldContainSubstring, "\n### Card 2.1.1 Title")
					err = nb.Promote()
					So(err, ShouldBeNil)
					body = nb.MarshalBody()
					So(body, ShouldNotContainSubstring, "\n### Card 2.1.1 Title")
					So(body, ShouldContainSubstring, "\n## Card 2.1.1 Title")

					Convey("But not past root level", func() {
						nb.Exit()
						title, body = nb.GetCard()
						So(title, ShouldEqual, "Card 2 Title")
						err := nb.Promote()
						So(err.Error(), ShouldEqual, "unable to promote any further")
						err = nb.PromoteAll(false)
						So(err.Error(), ShouldEqual, "unable to promote any further")

						nb2 := notebook.New("empty.md", "Empty", "Empty", "Empty")
						err = nb2.Promote()
						So(err.Error(), ShouldEqual, "unable to promote any further")
					})

					Convey("All current level children can be promoted", func() {
						nb.AddCard("2.1.1.1", "2.1.1.1", true)
						nb.AddCard("2.1.1.2", "2.1.1.2", false)
						body = nb.MarshalBody()
						So(body, ShouldContainSubstring, "\n### 2.1.1.1")
						So(body, ShouldContainSubstring, "\n### 2.1.1.2")
						err = nb.PromoteAll(false)
						So(err, ShouldBeNil)
						body = nb.MarshalBody()
						So(body, ShouldNotContainSubstring, "\n### 2.1.1.1")
						So(body, ShouldNotContainSubstring, "\n### 2.1.1.2")
						So(body, ShouldContainSubstring, "\n## 2.1.1.1")
						So(body, ShouldContainSubstring, "\n## 2.1.1.2")
					})

					Convey("Children can replace their parent", func() {
						nb.AddCard("2.1.1.1", "2.1.1.1", true)
						nb.AddCard("2.1.1.2", "2.1.1.2", false)
						body = nb.MarshalBody()
						So(body, ShouldContainSubstring, "\n## Card 2.1.1 Title")
						So(body, ShouldContainSubstring, "\n### 2.1.1.1")
						So(body, ShouldContainSubstring, "\n### 2.1.1.2")
						err = nb.PromoteAll(true)
						So(err, ShouldBeNil)
						body = nb.MarshalBody()
						So(body, ShouldNotContainSubstring, "\n## Card 2.1.1 Title")
						So(body, ShouldNotContainSubstring, "\n### 2.1.1.1")
						So(body, ShouldNotContainSubstring, "\n### 2.1.1.2")
						So(body, ShouldContainSubstring, "\n## 2.1.1.1")
						So(body, ShouldContainSubstring, "\n## 2.1.1.2")
					})
				})
			})
		})

		Convey("It can be marshalled and unmarshalled", func() {
			nb.AddRevision("Test revision")
			nb.AddCard("Test card", "test\n again", false)
			nb.AddCard("Test 2", "2", false)
			nb.AddCard("Child", "child", true)
			nb.Exit()
			nb.AddCard("Test 3", "3", false)
			marshalled := nb.Marshal()

			nb2, err := notebook.Unmarshal(marshalled)
			So(nb2, ShouldNotBeNil)
			So(err, ShouldBeNil)
			marshalled2 := nb2.Marshal()
			So(marshalled2, ShouldEqual, marshalled)

			_, err = notebook.Unmarshal("bad-wolf")
			So(err.Error(), ShouldEqual, "malformed notebook; must contain metadata block and body")
			_, err = notebook.Unmarshal("---\nbad---\nwolf")
			So(err.Error(), ShouldContainSubstring, "yaml: unmarshal errors")
			_, err = notebook.Unmarshal("---\n---\n\nbad-wolf")
			So(err.Error(), ShouldEqual, "malformed notebook; cannot have body without header")
			_, err = notebook.Unmarshal("---\n---\n\n#bad-wolf")
			So(err.Error(), ShouldEqual, "malformed notebook; title must contain depth marker and text")
			_, err = notebook.Unmarshal("---\n---\n\n## bad-wolf")
			So(err.Error(), ShouldEqual, "malformed notebook; must start at header depth 1, found ## bad-wolf")
			_, err = notebook.Unmarshal("---\n---\n\n# bad\n\n### wolf")
			So(err.Error(), ShouldEqual, "malformed notebook; header depths must increase by 1")

			Convey("It can do so with a file", func() {
				nb, err := notebook.Open("../README.md")
				So(err, ShouldBeNil)
				So(nb, ShouldNotBeNil)
				title, _ := nb.GetCard()
				So(title, ShouldEqual, "Mandelnote")

				Convey("Including a new file", func() {
					nb, err = notebook.Open("new-notebook.md")
					So(err, ShouldBeNil)
					So(nb, ShouldNotBeNil)
					So(nb.Title, ShouldEqual, "New notebook")
					f, _ := ioutil.TempFile("", "notebook.md")
					defer os.Remove(f.Name())
					nb.SetFile(f.Name())
					nb.Save()
					nb, err := notebook.Open(f.Name())
					So(err, ShouldBeNil)
					So(nb, ShouldNotBeNil)
					So(nb.Title, ShouldEqual, "New notebook")
				})

				Convey("But not a bad file", func() {
					nb, err = notebook.Open("..")
					So(err, ShouldNotBeNil)
					So(nb, ShouldBeNil)
					nb = notebook.New("..", "bad-wolf", "", "")
					err := nb.Save()
					So(err, ShouldNotBeNil)
				})

				Convey("Or an invalid file", func() {
					nb, err = notebook.Open("notebook.go")
					So(err, ShouldNotBeNil)
					So(nb, ShouldBeNil)
				})
			})
		})
	})
}
