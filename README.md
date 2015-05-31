# anno
Go package for text annotation.

## Usage

```
s := "Find http://www.websites.com/ and #hashtags and @mentions easily"
notes, err := anno.FindManyString(s, anno.URLs, anno.Mentions, anno.Hashtags)
if err != nil {
	log.Fatalln(err)
}
for _, note := range notes {
	log.Printf("Found a %s at [%d:%d]: \"%s\"", note.Kind, note.Start, note.End, note.Val)
}
```

Will output:

```
Found a url at [5:29]: "http://www.websites.com/"
Found a mention at [48:57]: "@mentions"
Found a hashtag at [34:43]: "#hashtags"
```

You can expand the notes using `anno.Expander`:

```
expander := anno.Expander{
	"url": func(b string) string {
		return fmt.Sprintf(`<a href="%[1]s">%[1]s</a>`, b)
	},
	"mention": func(b string) string {
		return fmt.Sprintf(`<a href="https://downlist.io/%[1]s">%[1]s</a>`, b)
	},
}
src := "This is a #long string written by @mat containing links to https://downlist.io/."
notes, err := anno.FindManyString(src, anno.URLs, anno.Mentions, anno.Hashtags)
is.NoErr(err)

log.Println(expander.Expand(src, notes))
```

Will output:

```
This is a #long string written by <a href="https://downlist.io/@mat">@mat</a> containing links to <a href="https://downlist.io/">https://downlist.io/</a>.
```