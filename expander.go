package anno

import "sort"

// Expander holds a map of note kinds to expander functions that will be
// used to replace the text in the original string.
type Expander map[string]func(b string) string

// Expand generates a new string by calling the expander function
// for each note depending on its kind.
func (e Expander) Expand(s string, notes Notes) string {

	// put the notes in the order in which they appear
	// in the string
	sort.Sort(notes)

	for i, note := range notes {
		fn, present := e[note.Kind]
		if !present {
			continue
		}
		insert := fn(string(note.Val))
		s = s[0:note.Start] + insert + s[note.End():]
		if i < len(notes)-1 {
			// update the offste of remaining notes
			offset := len(insert) - len(note.Val)
			for _, upcomingNote := range notes[i+1:] {
				upcomingNote.Start += offset
			}
		}

	}
	return s
}
