package anno

// Expander holds a map of note kinds to expander functions that will be
// used to replace the text in the original string.
type Expander map[string]func(b string) string

// Expand generates a new string by calling the expander function
// for each note depending on its kind.
func (e Expander) Expand(s string, notes []*Note) string {
	for _, note := range notes {
		fn, present := e[note.Kind]
		if !present {
			continue
		}
		s = s[0:note.Start] + fn(string(note.Val)) + s[note.End:]
	}
	return s
}
