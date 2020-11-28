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

			nb.AddCard("Card 1 Title", "Card 1 body")
			nb.AddCard("Card 2 Title", "Card 2 body")

			body := nb.MarshalBody()
			So(body, ShouldEqual, "\n# Card 1 Title\n\nCard 1 body\n\n# Card 2 Title\n\nCard 2 body\n")

			Convey("And delete them", func() {
				err := nb.Delete(false)
				So(err, ShouldBeNil)

				body = nb.MarshalBody()
				So(body, ShouldEqual, "\n# Card 1 Title\n\nCard 1 body\n")
			})
		})
	})
}
