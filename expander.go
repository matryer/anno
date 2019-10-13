package anno

// Expander holds a map of note kinds to expander functions that will be
// used to replace the text in the original string.
type Expander map[string]func(b string) string

// Expand generates a new string by calling the expander function
// for each note depending on its kind.
func (e Expander) Expand(s string, notes Notes) string {
	for i, note := range notes {
		fn, present := e[note.Kind]
		if !present {
			continue
		}
		insert := fn(string(note.Val))
		s = s[0:note.Start] + insert + s[note.End():]
		//	log.Printf("----- s after: %v", s)
		if i < len(notes)-1 {
			// update the offset of remaining notes
			offset := len(insert) - len(note.Val)
			for j := range notes[i+1:] {
				jj := j + i + 1
				notes[jj].Start += offset
			}
		}

	}
	return s
}
