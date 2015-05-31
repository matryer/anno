package anno

import "bytes"

var punctuation = "."

func trimPunctuation(s []byte) []byte {
	return bytes.TrimRight(s, punctuation)
}

// URL finds web addresses.
var URL = FieldFunc("url", func(s []byte) (bool, []byte) {
	return bytes.HasPrefix(s, []byte("http")) || bytes.HasPrefix(s, []byte("www")), s
})

// Mentions finds @twitter style mentions.
var Mentions = FieldFunc("mention", func(s []byte) (bool, []byte) {
	return bytes.HasPrefix(s, []byte("@")), trimPunctuation(s)
})

// Hashtags finds #hashtags.
var Hashtags = FieldFunc("hashtag", func(s []byte) (bool, []byte) {
	return bytes.HasPrefix(s, []byte("#")), trimPunctuation(s)
})
