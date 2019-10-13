package anno_test

import (
	"fmt"
	"testing"

	"github.com/matryer/anno"
	"github.com/matryer/is"
)

func TestExpander(t *testing.T) {
	is := is.New(t)

	expander := anno.Expander{
		"url": func(b string) string {
			return fmt.Sprintf(`<a href="%[1]s">%[1]s</a>`, b)
		},
		"mention": func(b string) string {
			return fmt.Sprintf(`<a href="https://downlist.io/%[1]s">%[1]s</a>`, b)
		},
	}
	src := "This is a #long string written by @mat containing links to https://downlist.io/."
	notes, err := anno.FindManyString(src, anno.Mentions, anno.URLs, anno.Hashtags)
	is.NoErr(err)

	out := expander.Expand(src, notes)
	is.Equal(out, `This is a #long string written by <a href="https://downlist.io/@mat">@mat</a> containing links to <a href="https://downlist.io/">https://downlist.io/</a>.`)

}

func TestExpanderDuplicates(t *testing.T) {
	is := is.New(t)

	expander := anno.Expander{
		"mention": func(b string) string {
			return fmt.Sprintf(`<a href="/profile/%[1]s">%[1]s</a>`, b)
		},
	}
	src := "Hello from @mat and @david and @david and @mat :)"
	notes, err := anno.FindManyString(src, anno.Mentions, anno.URLs, anno.Hashtags)
	is.NoErr(err)

	out := expander.Expand(src, notes)
	is.Equal(out, `Hello from <a href="/profile/@mat">@mat</a> and <a href="/profile/@david">@david</a> and <a href="/profile/@david">@david</a> and <a href="/profile/@mat">@mat</a> :)`)

}
