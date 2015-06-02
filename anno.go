package anno

import (
	"bytes"
	"sync"
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

// FieldFunc returns a FinderFunc that finds notes on a per field basis.
// The fn returns true if it's a match, and optionally a subset of the the
// match.
// The returned []byte must be contained in the source string, otherwise
// ErrNoMatch will be returned.
func FieldFunc(kind string, fn func(b []byte) (bool, []byte)) FinderFunc {
	return func(src []byte) (Notes, error) {
		var notes Notes
		fields := bytes.Fields(src)
		noteChan := make(chan *Note, 0)
		errChan := make(chan error, 0)

		go func() {
			var wg sync.WaitGroup
			for _, f := range fields {
				wg.Add(1)
				go func(f []byte) {
					if ok, match := fn(f); ok {
						s := bytes.Index(src, match)
						if s == -1 {
							// true was returned without the returned bytes
							// appearing in the match.
							errChan <- ErrNoMatch(match)
						}
						noteChan <- &Note{
							Val:   match,
							Start: s,
							Kind:  kind,
						}
					}
					wg.Done()
				}(f)
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
