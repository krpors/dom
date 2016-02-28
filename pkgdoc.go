// Package dom is an implementation of the DOM specification, level 3.
// Since elements, text nodes, comments, processing instructions, etc. cannot exist
// outside the context of a Document, the Document interface also contains the
// factory methods needed to create these objects. The Node objects created have an
// ownerDocument attribute which associates them with the Document within whose
// context they were created.
//
// Example:
/*
	package dom

	import (
		"bytes"
		"fmt"
		"github.com/krpors/dom"
	)

	func main() {
		doc := dom.NewDocument()
		root, _ := doc.CreateElement("rootNode")
		text := doc.CreateText("some arbitrary text")
		root.AppendChild(text)
		doc.AppendChild(root)

		var b bytes.Buffer
		dom.ToXML(doc, false, &b)
		fmt.Println(b.String())
	}
*/
// Running this code will yield:
//	<?xml version="1.0" encoding="UTF-8"?><rootNode>some arbitrary text</rootNode>
// References:
//	https://www.w3.org/TR/2004/REC-DOM-Level-3-Core-20040407/Overview.html#contents
//	https://www.w3.org/TR/xml/
package dom
