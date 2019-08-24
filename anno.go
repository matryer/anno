package anno

import (
	"bytes"
	"sync"
	"unicode"
	"unicode/utf8"
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

	noteChan := make(chan *Note, 0)
	errChan := make(chan error, 0)

	go func() {
		var wg sync.WaitGroup
		for _, finder := range finders {
			wg.Add(1)

			go func(finder Finder) {

				notes, err := finder.Find(src)
				if err != nil {
					errChan <- err
				}
				for _, note := range notes {
					noteChan <- note
				}

				wg.Done()
			}(finder)

		}
		wg.Wait()
		close(noteChan)
	}()

	var notes Notes
	var err error
loop:
	for {
		select {
		case note := <-noteChan:
			if note == nil {
				break loop
			}
			notes = append(notes, note)
		case e := <-errChan:
			err = e
			break loop
		}
	}

	return notes, err

}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

/* An outright copy of bytes.Fields that returns a map of index=>field
	so that we don't have bugs created by duplicate fields and
	bytes.Index */
func annoFields(s []byte) map[int][]byte {
	// First count the fields.
	// This is an exact count if s is ASCII, otherwise it is an approximation.
	n := 0
	wasSpace := 1
	// setBits is used to track which bits are set in the bytes of s.
	setBits := uint8(0)
	for i := 0; i < len(s); i++ {
		r := s[i]
		setBits |= r
		isSpace := int(asciiSpace[r])
		n += wasSpace & ^isSpace
		wasSpace = isSpace
	}

	if setBits >= utf8.RuneSelf {
		// Some runes in the input slice are not ASCII.
		return annoFieldsFunc(s, unicode.IsSpace)
	}

	// ASCII fast path
	a := make(map[int][]byte, n)
	na := 0
	fieldStart := 0
	i := 0
	// Skip spaces in the front of the input.
	for i < len(s) && asciiSpace[s[i]] != 0 {
		i++
	}
	fieldStart = i
	for i < len(s) {
		if asciiSpace[s[i]] == 0 {
			i++
			continue
		}
		a[fieldStart] = s[fieldStart:i:i]
		na++
		i++
		// Skip spaces in between fields.
		for i < len(s) && asciiSpace[s[i]] != 0 {
			i++
		}
		fieldStart = i
	}
	if fieldStart < len(s) { // Last field might end at EOF.
		a[fieldStart] = s[fieldStart:len(s):len(s)]
	}
	return a
}

func annoFieldsFunc(s []byte, f func(rune) bool) map[int][]byte {
	// A span is used to record a slice of s of the form s[start:end].
	// The start index is inclusive and the end index is exclusive.
	type span struct {
		start int
		end   int
	}
	spans := make([]span, 0, 32)

	// Find the field start and end indices.
	wasField := false
	fromIndex := 0
	for i := 0; i < len(s); {
		size := 1
		r := rune(s[i])
		if r >= utf8.RuneSelf {
			r, size = utf8.DecodeRune(s[i:])
		}
		if f(r) {
			if wasField {
				spans = append(spans, span{start: fromIndex, end: i})
				wasField = false
			}
		} else {
			if !wasField {
				fromIndex = i
				wasField = true
			}
		}
		i += size
	}

	// Last field might end at EOF.
	if wasField {
		spans = append(spans, span{fromIndex, len(s)})
	}

	// Create subslices from recorded field indices.
	a := make(map[int][]byte, len(spans))
	for _, span := range spans {
		a[span.start] = s[span.start:span.end:span.end]
	}

	return a
}

// FieldFunc returns a FinderFunc that finds notes on a per field basis.
// The fn returns true if it's a match, and optionally a subset of the the
// match.
// The returned []byte must be contained in the source string, otherwise
// ErrNoMatch will be returned.
func FieldFunc(kind string, fn func(b []byte) (bool, []byte)) FinderFunc {
	return func(src []byte) (Notes, error) {
		var notes Notes
		fields := annoFields(src)
		noteChan := make(chan *Note, 0)
		errChan := make(chan error, 0)

		go func() {
			var wg sync.WaitGroup
			for i, f := range fields {
				wg.Add(1)
				go func(i int, f []byte) {
					if ok, match := fn(f); ok {
						s := bytes.Index(f, match)
						if s == -1 {
							// true was returned without the returned bytes
							// appearing in the match.
							errChan <- ErrNoMatch(match)
						}
						noteChan <- &Note{
							Val:   match,
							Start: i+s,
							Kind:  kind,
						}
					}
					wg.Done()
				}(i, f)
			}
			wg.Wait()
			close(noteChan)
		}()

		// read notes and build up the array
		var err error
	loop:
		for {
			select {
			case note := <-noteChan:
				if note == nil {
					break loop
				}
				notes = append(notes, note)
			case e := <-errChan:
				err = e
				break loop
			}
		}

		return notes, err
	}
}
