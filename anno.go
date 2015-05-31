package anno

import (
	"bytes"
	"sync"
)

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

// FindManyString runs all finders against the source and returns a
// slice of notes or an error.
func FindManyString(src string, finders ...Finder) ([]*Note, error) {
	return FindMany([]byte(src), finders...)
}

// FindMany runs all finders against the source and returns a
// slice of notes or an error.
func FindMany(src []byte, finders ...Finder) ([]*Note, error) {

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

	var notes []*Note
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
	return func(src []byte) ([]*Note, error) {
		var notes []*Note
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
							End:   s + len(match),
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
