package anno

import "bytes"

var punctuation = ".?"

func trimPunctuation(s []byte) []byte {
	return bytes.TrimRight(s, punctuation)
}

var tlds = [][]byte{
	[]byte(".com"),
	[]byte(".net"),
	[]byte(".org"),
	[]byte(".edu"),
	[]byte(".co"),
	[]byte(".io"),
	[]byte(".uk"),
	[]byte(".us"),
	[]byte(".ca"),
	[]byte(".de"),
	[]byte(".jp"),
	[]byte(".fr"),
	[]byte(".au"),
	[]byte(".ru"),
	[]byte(".ch"),
	[]byte(".it"),
	[]byte(".nl"),
	[]byte(".se"),
	[]byte(".no"),
	[]byte(".es"),
	[]byte(".mil"),
}

// URL finds web addresses.
var URLs = FieldFunc("url", func(s []byte) (bool, []byte) {
	trimmedS := trimPunctuation(s)
	if bytes.HasPrefix(s, []byte("http")) || bytes.HasPrefix(s, []byte("www")) {
		return true, trimmedS
	}
	for _, tld := range tlds {
		if bytes.HasSuffix(s, tld) {
			return true, trimmedS
		}
	}
	return false, s
})

// Mentions finds @twitter style mentions.
var Mentions = FieldFunc("mention", func(s []byte) (bool, []byte) {
	return bytes.HasPrefix(s, []byte("@")), trimPunctuation(s)
})

// Hashtags finds #hashtags.
var Hashtags = FieldFunc("hashtag", func(s []byte) (bool, []byte) {
	return bytes.HasPrefix(s, []byte("#")), trimPunctuation(s)
})
