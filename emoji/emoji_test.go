package emoji_test

import (
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/anno"
	"github.com/matryer/anno/emoji"
)

func TestEmoji(t *testing.T) {
	is := is.New(t)
	src := []byte(":beer: makes me want to :smile: you know.")
	notes, err := emoji.Find(src)
	is.NoErr(err)
	is.OK(notes)
	is.Equal(len(notes), 2)
	is.Equal(notes[0].Val, []byte(":beer:"))
	is.Equal(notes[0].Start, 0)
	is.Equal(notes[0].End(), 0+len(notes[0].Val))
	is.Equal(notes[0].Kind, "emoji")
	is.Equal(notes[1].Val, []byte(":smile:"))
	is.Equal(notes[1].Start, 24)
	is.Equal(notes[1].End(), 24+len(notes[1].Val))
	is.Equal(notes[1].Kind, "emoji")

	e := anno.Expander{
		"emoji": emoji.Expand,
	}
	is.Equal(e.Expand(string(src), notes), "🍺 makes me want to 😄 you know.")

}

func TestEmojiWithPunctuation(t *testing.T) {
	is := is.New(t)
	src := []byte(`":beer:"?`)
	notes, err := emoji.Find(src)
	is.NoErr(err)
	is.OK(notes)
	is.Equal(len(notes), 1)
	is.Equal(notes[0].Val, []byte(`:beer:`))
	is.Equal(notes[0].Start, 1)
	is.Equal(notes[0].End(), 1+len(notes[0].Val))
	is.Equal(notes[0].Kind, "emoji")

	e := anno.Expander{
		"emoji": emoji.Expand,
	}
	is.Equal(e.Expand(string(src), notes), `"🍺"?`)

}
