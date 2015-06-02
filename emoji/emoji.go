package emoji

import (
	"bytes"

	"github.com/matryer/anno"
	"gopkg.in/kyokomi/emoji.v1"
)

// Find finds raw :emoji: fields.
var Find = anno.FieldFunc("emoji", func(s []byte) (bool, []byte) {
	return bytes.HasPrefix(s, []byte(":")) && bytes.HasSuffix(s, []byte(":")), s
})

// Expand takes a single emoji source and returns the actual
// characters.
func Expand(b string) string {
	return emoji.CodeMap()[b]
}
