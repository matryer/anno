package anno_test

import (
	"testing"

	"github.com/cheekybits/is"
	"github.com/downlist/anno"
)

func TestMany(t *testing.T) {
	is := is.New(t)
	f1 := anno.FinderFunc(func(s []byte) ([]*anno.Note, error) {
		return []*anno.Note{{
			Val: []byte("one"),
		}}, nil
	})
	f2 := anno.FinderFunc(func(s []byte) ([]*anno.Note, error) {
		return []*anno.Note{{
			Val: []byte("two"),
		}}, nil
	})
	notes, err := anno.Many(f1, f2).Find([]byte(""))
	is.NoErr(err)
	is.Equal(len(notes), 2)
}

func TestFindString(t *testing.T) {
	is := is.New(t)

	s := "This is a #long string written by @mat containing links to https://downlist.io/."
	finder := anno.Many(anno.Hashtags, anno.Mentions, anno.URLs)
	notes, err := anno.FindString(finder, s)
	is.NoErr(err)
	is.Equal(len(notes), 3)

	is.Equal(notes[0].String(), "#long")
	is.Equal(notes[1].String(), "@mat")
	is.Equal(notes[2].String(), "https://downlist.io/")

}

func TestFieldFinder(t *testing.T) {
	is := is.New(t)

	s := []byte("A field finder will be able to find fields.")
	var fn anno.Finder
	fn = anno.FieldFunc("thiskind", func(f []byte) (bool, []byte) {
		if string(f) == "field" {
			return true, f
		}
		if string(f) == "able" {
			return true, f
		}
		if string(f) == "find" {
			return true, f
		}
		return false, f
	})
	notes, err := fn.Find(s)
	is.NoErr(err)

	is.Equal(len(notes), 3)
	is.Equal(notes[0].Val, []byte("field"))
	is.Equal(notes[0].Start, 2)
	is.Equal(notes[0].End, 2+len(notes[0].Val))
	is.Equal(notes[0].Kind, "thiskind")

	is.Equal(notes[1].Val, []byte("able"))
	is.Equal(notes[1].Start, 23)
	is.Equal(notes[1].End, 23+len(notes[1].Val))
	is.Equal(notes[1].Kind, "thiskind")

	is.Equal(notes[2].Val, []byte("find"))
	is.Equal(notes[2].Start, 8)
	is.Equal(notes[2].End, 8+len(notes[2].Val))
	is.Equal(notes[2].Kind, "thiskind")

}
