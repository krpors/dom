# DOM Level 3 for Go

Experimental attempt for a complete implementation of the DOM, level 3,
using the specs [from W3.org](https://www.w3.org/TR/2004/REC-DOM-Level-3-Core-20040407/core.html#i-Document),
in the Go language.

There are several other implementations floating about. This is currently just
an experiment to see how far I can get, and perhaps it can grow into a usable
implementation.

## Interfaces

Since the spec for DOM3 regards everything as interfaces, I tried to do the
same in this package. Whether that will work out well or not is the question.
Right now, the Document type is (per spec) the main entry point for creating
any kind of Node. The following interfaces have (partial) implementations:

* `Document`: the entry point for creating Nodes.
* `Element`: for example: `<pfx:element/>`
* `Attr`: attributes of elements, for example: `<pfx:element pfx:attribute="hi"/>`
* `ProcessingInstruction`: for example: `<?spacing true?>`
* `Comment`: for example: `<!-- comment node -->`
* `Text`: basic text as a child of an Element

The following are omitted:

* `NodeList`: is just too convoluted to implement this as well IMO. A slice is sufficient.
* `CDATASection`: it's just `Text`.
* `CharacterData`: which is just a String type with methods defined on it.

## Example code

Reading XML documents is pretty much working, and the DOM can be traversed using the
supplied methods. An example program using this DOM implementation:

```go
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/krpors/dom"
)

func main() {
	resp, err := http.Get("https://www.reddit.com/r/golang/.xml")
	if err != nil {
		fmt.Println("failed:", err)
		os.Exit(1)
	}

	builder := dom.NewBuilder(resp.Body)
	doc, err := builder.Parse()
	if err != nil {
		os.Exit(2)
	}

	entries := doc.GetElementsByTagName("entry")
	for _, entry := range entries {
		username := entry.GetElementsByTagName("name")[0]
		title := entry.GetElementsByTagName("title")[0]

		fmt.Printf("'%s' posted '%s'\n", username.GetFirstChild().GetNodeValue(), title.GetFirstChild().GetNodeValue())
	}
}
```

This will get an XML feed from the [r/golang](https://reddit.com/r/golang) subreddit,
and prints out the `name` and `title` nodes from each `entry` node. It's obviously a bit
verbose and probably not really 'idiomatic' Go, but the goal is mostly to be consistent
with the DOM specification.

## Challenges

1. No inheritance, meaning double implementations of the same functionality (Element, Comment,
	ProcessingInstruction, ...). More code, but better separation.
1. Spec mentions things like 'if x is null, then ...', but Go doesn't have null pointers for
	Strings, for example.
