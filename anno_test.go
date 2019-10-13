package anno_test

import (
	"sort"
	"testing"

	"github.com/matryer/anno"
	"github.com/matryer/is"
)

func TestFindMany(t *testing.T) {
	is := is.New(t)
	src := []byte("This is a #long string written by @mat containing links to https://downlist.io/.")

	notes, err := anno.FindMany(src, anno.URLs, anno.Mentions, anno.Hashtags)
	is.NoErr(err)
	is.Equal(len(notes), 3)

}

func TestFindManyString(t *testing.T) {
	is := is.New(t)
	src := "This is a #long string written by @mat containing links to https://downlist.io/."

	notes, err := anno.FindManyString(src, anno.URLs, anno.Mentions, anno.Hashtags)
	is.NoErr(err)
	is.Equal(len(notes), 3)

}

func TestFindManyStringDuplicates(t *testing.T) {
	is := is.New(t)
	src := `If I @mention and then I @mention again, I hope to get different start values`

	notes, err := anno.FindManyString(src, anno.Mentions)
	is.NoErr(err)
	is.Equal(len(notes), 2)

	is.True(notes[0].Start < notes[1].Start)

}

func TestFindManyURLsWithPunctuation(t *testing.T) {
	is := is.New(t)
	src := "What do you think about Facebook.com and Yahoo.com and Google.com?"

	notes, err := anno.FindManyString(src, anno.URLs)
	is.NoErr(err)
	is.Equal(len(notes), 3)

}

func TestFindString(t *testing.T) {
	is := is.New(t)
	s := "This is a #long string written by @mat containing links to https://downlist.io/."
	notes, err := anno.FindString(anno.URLs, s)
	is.NoErr(err)
	is.Equal(len(notes), 1)
	is.Equal(notes[0].String(), "https://downlist.io/")
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
	is.Equal(string(notes[0].Val), "field")
	is.Equal(notes[0].Start, 2)
	is.Equal(notes[0].End(), 2+len(notes[0].Val))
	is.Equal(notes[0].Kind, "thiskind")

	is.Equal(string(notes[1].Val), "able")
	is.Equal(notes[1].Start, 23)
	is.Equal(notes[1].End(), 23+len(notes[1].Val))
	is.Equal(notes[1].Kind, "thiskind")

	is.Equal(string(notes[2].Val), "find")
	is.Equal(notes[2].Start, 31)
	is.Equal(notes[2].End(), 31+len(notes[2].Val))
	is.Equal(notes[2].Kind, "thiskind")

	// sort the notes
	sort.Sort(notes)

	is.Equal(len(notes), 3)
	is.Equal(string(notes[0].Val), "field")
	is.Equal(notes[0].Start, 2)
	is.Equal(notes[0].End(), 2+len(notes[0].Val))
	is.Equal(notes[0].Kind, "thiskind")

	is.Equal(string(notes[1].Val), "able")
	is.Equal(notes[1].Start, 23)
	is.Equal(notes[1].End(), 23+len(notes[1].Val))
	is.Equal(notes[1].Kind, "thiskind")

	is.Equal(string(notes[2].Val), "find")
	is.Equal(notes[2].Start, 31)
	is.Equal(notes[2].End(), 31+len(notes[2].Val))
	is.Equal(notes[2].Kind, "thiskind")

}
