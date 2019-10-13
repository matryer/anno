package anno

import (
	"bytes"
)

// Finder represents types capable of finding
// notes.
type Finder interface {
	Find(s []byte) (Notes, error)
}

// Notes is a sortable slice of Note objects.
type Notes []*Note

func (n Notes) Len() int           { return len(n) }
func (n Notes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Notes) Less(i, j int) bool { return n[i].Start < n[j].Start }

// Note represents something interesting within
// text.
type Note struct {
	Val   []byte `json:"val"`
	Start int    `json:"start"`
	Kind  string `json:"kind"`
}

// End calculates the end position of this note.
func (n *Note) End() int {
	return n.Start + len(n.Val)
}

func (n *Note) String() string {
	return string(n.Val)
}

// FindString uses the Finder to find notes within
// the specified string.
func FindString(finder Finder, s string) (Notes, error) {
	return finder.Find([]byte(s))
}

// FinderFunc represents a function capable of finding
// notes.
type FinderFunc func(s []byte) (Notes, error)

var _ Finder = FinderFunc(nil)

// Find calls the FinderFunc.
func (fn FinderFunc) Find(s []byte) (Notes, error) {
	return fn(s)
}

// ErrNoMatch is returned when a match was expected but
// can not be found.
// Can be returned from FieldFunc.
type ErrNoMatch []byte

func (e ErrNoMatch) Error() string {
	return "no match for '" + string(e) + "'"
}

// FindManyString runs all finders against the source and returns a
// slice of notes or an error.
func FindManyString(src string, finders ...Finder) (Notes, error) {
	return FindMany([]byte(src), finders...)
}

// FindMany runs all finders against the source and returns a
// slice of notes or an error.
func FindMany(src []byte, finders ...Finder) (Notes, error) {
	var allNotes Notes
	for _, finder := range finders {
		notes, err := finder.Find(src)
		if err != nil {
			return nil, err
		}
		allNotes = append(allNotes, notes...)
	}
	return allNotes, nil
}

// FieldFunc returns a FinderFunc that finds notes on a per field basis.
// The fn returns true if it's a match, and optionally a subset of the the
// match.
// The returned []byte must be contained in the source string, otherwise
// ErrNoMatch will be returned.
func FieldFunc(kind string, fn func(b []byte) (bool, []byte)) FinderFunc {
	return func(src []byte) (Notes, error) {
		var notes Notes
		fields := bytes.Fields(src)
		for _, f := range fields {
			if ok, match := fn(f); ok {
				s := bytes.Index(src, match)
				if s == -1 {
					// true was returned without the returned bytes
					// appearing in the match.
					return nil, ErrNoMatch(match)
				}
				notes = append(notes, &Note{
					Val:   match,
					Start: s,
					Kind:  kind,
				})
			}
		}
		return notes, nil
	}
}
