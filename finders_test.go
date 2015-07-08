package anno_test

import (
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/anno"
)

func TestURL(t *testing.T) {
	is := is.New(t)
	src := []byte("My website is https://downlist.io/ come and check it out - or go to http://www.codeandthat.com/ instead.")
	matches, err := anno.URLs(src)
	is.NoErr(err)
	is.OK(matches)
	is.Equal(len(matches), 2)
	is.Equal(matches[0].Val, []byte("https://downlist.io/"))
	is.Equal(matches[0].Start, 14)
	is.Equal(matches[0].End(), 14+len(matches[0].Val))
	is.Equal(matches[0].Kind, "url")
	is.Equal(matches[1].Val, []byte("http://www.codeandthat.com/"))
	is.Equal(matches[1].Start, 68)
	is.Equal(matches[1].End(), 68+len(matches[1].Val))
	is.Equal(matches[1].Kind, "url")
}

func TestEmail(t *testing.T) {
	is := is.New(t)
	src := []byte("Send me an email to please-reply@downlist.io if you like.")
	matches, err := anno.Emails(src)
	is.NoErr(err)
	is.OK(matches)
	is.Equal(len(matches), 1)
	is.Equal(matches[0].Val, []byte("please-reply@downlist.io"))
	is.Equal(matches[0].Start, 20)
	is.Equal(matches[0].End(), 20+len(matches[0].Val))
	is.Equal(matches[0].Kind, "email")
}

func TestMention(t *testing.T) {
	is := is.New(t)
	src := []byte("Call me @matryer on Twitter, or follow @downlistapp instead.")
	matches, err := anno.Mentions(src)
	is.NoErr(err)
	is.OK(matches)
	is.Equal(len(matches), 2)
	is.Equal(matches[0].Val, []byte("@matryer"))
	is.Equal(matches[0].Start, 8)
	is.Equal(matches[0].End(), 8+len(matches[0].Val))
	is.Equal(matches[0].Kind, "mention")
	is.Equal(matches[1].Val, []byte("@downlistapp"))
	is.Equal(matches[1].Start, 39)
	is.Equal(matches[1].End(), 39+len(matches[1].Val))
	is.Equal(matches[1].Kind, "mention")
}

func TestHashtag(t *testing.T) {
	is := is.New(t)
	src := []byte("I love programming in #golang - it's #lovely.")
	matches, err := anno.Hashtags(src)
	is.NoErr(err)
	is.OK(matches)
	is.Equal(len(matches), 2)
	is.Equal(matches[0].Val, []byte("#golang"))
	is.Equal(matches[0].Start, 22)
	is.Equal(matches[0].End(), 22+len(matches[0].Val))
	is.Equal(matches[0].Kind, "hashtag")
	is.Equal(matches[1].Val, []byte("#lovely"))
	is.Equal(matches[1].Start, 37)
	is.Equal(matches[1].End(), 37+len(matches[1].Val))
	is.Equal(matches[1].Kind, "hashtag")
}
