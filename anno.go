package anno

import "bytes"

// Finder represents types capable of finding
// notes.
type Finder interface {
	Find(s []byte) ([]*Note, error)
}

// Note represents something interesting within
// text.
type Note struct {
	Val   []byte
	Start int
	End   int
	Kind  string
}

func (n *Note) String() string {
	return string(n.Val)
}

// FindString uses the Finder to find notes within
// the specified string.
func FindString(finder Finder, s string) ([]*Note, error) {
	return finder.Find([]byte(s))
}

// Many wraps multiple Finder types into one and its Find
// method will run them all in order.
func Many(finders ...Finder) Finder {
	return FinderFunc(func(s []byte) ([]*Note, error) {
		var notes []*Note
		for _, finder := range finders {
			ns, err := finder.Find(s)
			if err != nil {
				return notes, err
			}
			notes = append(notes, ns...)
		}
		return notes, nil
	})
}

// FinderFunc represents a function capable of finding
// notes.
type FinderFunc func(s []byte) ([]*Note, error)

var _ Finder = FinderFunc(nil)

// Find calls the FinderFunc.
func (fn FinderFunc) Find(s []byte) ([]*Note, error) {
	return fn(s)
}

// ErrNoMatch is returned when a match was expected but
// can not be found.
// Can be returned from FieldFunc.
type ErrNoMatch []byte

func (e ErrNoMatch) Error() string {
	return "no match for '" + string(e) + "'"
}

// FieldFunc returns a FinderFunc that finds notes on a per field basis.
// The fn returns true if it's a match, and optionally a subset of the the
// match.
// The returned []byte must be contained in the source string, otherwise
// ErrNoMatch will be returned.
func FieldFunc(kind string, fn func(b []byte) (bool, []byte)) FinderFunc {
	return func(src []byte) ([]*Note, error) {
		var notes []*Note
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
					End:   s + len(match),
					Kind:  kind,
				})
			}
		}
		return notes, nil
	}
}
