package storage

import (
	"testing"

	"github.com/hanzki/moviebox-server/core"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSave(t *testing.T) {
	Convey("Given new storage mock", t, func() {
		mockStorage := NewMock()

		Convey("Should update search result ID", func() {
			r := &core.SearchResult{
				Title: "Testinator II",
			}

			r, err := mockStorage.Save(r)
			check(err)

			So(r.ID, ShouldNotEqual, "")
		})

		Convey("Should update CreatedAt", func() {
			r := &core.SearchResult{
				Title: "Testinator II",
			}

			r, err := mockStorage.Save(r)
			check(err)

			So(r.CreatedAt.IsZero(), ShouldBeFalse)
		})
	})

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func TestLoad(t *testing.T) {

	Convey("Given mock storage with existing results", t, func() {
		mockStorage := NewMock()
		r1 := &core.SearchResult{
			Title: "Testinator II",
		}
		r1, err := mockStorage.Save(r1)
		check(err)

		Convey("Should load existing result", func() {
			r, err := mockStorage.Load(r1.ID)
			check(err)

			So(r.Title, ShouldEqual, r1.Title)
		})

		Convey("Should return error if matching result doesn't exist", func() {
			_, err := mockStorage.Load("Not a real ID")
			So(err, ShouldNotBeNil)
		})
	})
}

func TestUpdate(t *testing.T) {
	Convey("Given mock storage with existing results", t, func() {
		mockStorage := NewMock()
		r1 := &core.SearchResult{
			Title: "Testinator II",
		}
		r1, err := mockStorage.Save(r1)
		check(err)

		Convey("Should allow updating existing result", func() {
			r := &core.SearchResult{}
			*r = *r1
			r.Title = "Mocky V"
			So(r.Title, ShouldNotEqual, r1.Title)

			err := mockStorage.Update(r)
			check(err)

			r, err = mockStorage.Load(r.ID)
			check(err)

			So(r.Title, ShouldEqual, "Mocky V")
		})
	})
}
