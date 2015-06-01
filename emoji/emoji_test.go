package emoji_test

import (
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/anno"
	"github.com/matryer/anno/emoji"
)

func TestEmoji(t *testing.T) {
	is := is.New(t)
	src := []byte("You make me want to :smile: you know.")
	notes, err := emoji.Find(src)
	is.NoErr(err)
	is.OK(notes)
	is.Equal(len(notes), 1)
	is.Equal(notes[0].Val, []byte(":smile:"))
	is.Equal(notes[0].Start, 20)
	is.Equal(notes[0].End, 20+len(notes[0].Val))
	is.Equal(notes[0].Kind, "emoji")

	e := anno.Expander{
		"emoji": emoji.Expand,
	}
	is.Equal(e.Expand(string(src), notes), "You make me want to 😄 you know.")

}
