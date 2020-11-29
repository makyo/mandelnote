package notebook_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	yaml "gopkg.in/yaml.v2"

	"github.com/makyo/mandelnote/notebook"
)

func TestConfig(t *testing.T) {
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
			nb.AddCard("Card 2 Title", "Card 2 body", false)

			body := nb.MarshalBody()
			So(body, ShouldEqual, "\n# Card 1 Title\n\nCard 1 body\n\n# Card 2 Title\n\nCard 2 body\n")

			Convey("And edit them", func() {
				nb.EditCard("Rose", "Companion")

				body = nb.MarshalBody()
				So(body, ShouldEqual, "\n# Card 1 Title\n\nCard 1 body\n\n# Rose\n\nCompanion\n")
			})

			Convey("And delete them", func() {
				err := nb.Delete(false)
				So(err, ShouldBeNil)

				body = nb.MarshalBody()
				So(body, ShouldEqual, "\n# Card 1 Title\n\nCard 1 body\n")
			})

			Convey("Cards can have children", func() {
				nb.AddCard("Card 2.1 Title", "Card 2.1 body", true)

				body = nb.MarshalBody()
				So(body, ShouldEqual, "\n# Card 1 Title\n\nCard 1 body\n\n# Card 2 Title\n\nCard 2 body\n\n## Card 2.1 Title\n\nCard 2.1 body\n")

				Convey("One can move between cards", func() {
					title, body := nb.GetCard()
					So(title, ShouldEqual, "Card 2.1 Title")
					So(body, ShouldEqual, "Card 2.1 body")

					nb.Exit()
					title, body = nb.GetCard()
					So(title, ShouldEqual, "Card 2 Title")
					So(body, ShouldEqual, "Card 2 body")

					nb.Cycle(-1)
					title, body = nb.GetCard()
					So(title, ShouldEqual, "Card 1 Title")
					So(body, ShouldEqual, "Card 1 body")

					nb.Cycle(1)
					title, body = nb.GetCard()
					So(title, ShouldEqual, "Card 2 Title")
					So(body, ShouldEqual, "Card 2 body")

					nb.Enter()
					title, body = nb.GetCard()
					So(title, ShouldEqual, "Card 2.1 Title")
					So(body, ShouldEqual, "Card 2.1 body")
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

					nb.Cycle(-2)
					nb.Enter()
					nb.AddCard("Card 2.1.1 Title", "Card 2.1.1 body", true)
					title, body = nb.GetCard()
					So(title, ShouldEqual, "Card 2.1.1 Title")
					body = nb.MarshalBody()
					So(body, ShouldEqual, "\n### Card 2.1.1 Title")
					err = nb.Promote()
					So(err, ShouldBeNil)
					body = nb.MarshalBody()
					So(body, ShouldNotContainSubstring, "\n### Card 2.1.1 Title")
					So(body, ShouldEqual, "\n## Card 2.1.1 Title")
				})
			})
		})
	})
}
