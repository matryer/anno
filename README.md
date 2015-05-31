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
