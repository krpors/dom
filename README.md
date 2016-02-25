# DOM Level 3 for Go

Experimental attempt for a complete implementation of the DOM, level 3,
using the specs [from W3.org](https://www.w3.org/TR/2004/REC-DOM-Level-3-Core-20040407/core.html#i-Document),
in the Go language.

There are several other implementations floating about. This is currently just
an experiment to see how far I can get, and perhaps it can grow into a usable
implementation.

## Interfaces, interfaces everywhere

Since the spec for DOM3 regards everything as interfaces, I tried to do the
same in this package. Whether that will work out well or not is the question.
Right now, the Document type is (per spec) the main entrypoint for creating
any kind of Node.


## Example code

The API in the workings. Serialization and deserialization are obviously
a big todo so this only checks the API.

```go
package main

import (
	"fmt"
	"github.com/krpors/dom"
)

func tree(n dom.Node, padding string) {
	fmt.Printf("%s%v\n", padding, n)
	for _, child := range n.GetChildNodes() {
		tree(child, padding+"    ")
	}
}

func main() {
	doc := dom.NewDocument()

	root, _ := doc.CreateElement("root")

	sub1, _ := doc.CreateElement("one")
	txt1 := doc.CreateTextNode("sample text 1")
	sub1.AppendChild(txt1)
	sub1.SetAttribute("cruft", "twelve")
	sub1.SetAttribute("once", "twice")

	nnm := sub1.GetAttributes()
	fmt.Println(nnm.Length())
	fmt.Println(nnm.GetNamedItem("cruft").GetNodeValue())

	sub2, _ := doc.CreateElement("two")
	txt2 := doc.CreateTextNode("sample text 2")
	sub2.AppendChild(txt2)

	root.AppendChild(sub1)
	root.AppendChild(sub2)

	doc.AppendChild(root)

	tree(doc, "")
}
```
