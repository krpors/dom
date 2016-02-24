package dom

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

// ToXML is a utility function to serialize a Node and its children to the
// given writer 'w'. The whole tree is traversed using recursion, which is
// pretty much eligible for a refactor for optimalization reasons.
func ToXML(node Node, omitXMLDecl bool, w io.Writer) {
	// Must define the function here so we can refer to ourselves in
	// the traverse function.
	var traverse func(n Node)

	if !omitXMLDecl {
		fmt.Fprintf(w, "%s", XMLDeclaration)
	}

	traverse = func(n Node) {
		switch t := n.(type) {
		case Element:
			if t.HasChildNodes() {
				fmt.Fprintf(w, "<%s>", t.GetTagName())
			} else {
				fmt.Fprintf(w, "<%s/>", t.GetTagName())
			}
		case Text:
			fmt.Fprintf(w, "%s", escape(t.GetData()))
		case Comment:
			fmt.Fprintf(w, "<!-- %s -->", t.GetComment())
		}

		// For each child node, call traverse() again.
		for _, node := range n.NodeList() {
			traverse(node)
		}

		switch t := n.(type) {
		case Element:
			if t.HasChildNodes() {
				fmt.Fprintf(w, "</%s>", t.GetTagName())
			}
		}
	}

	traverse(node)
}

// escape is a convenience method to escape XML for serialization.
func escape(s string) string {
	var b bytes.Buffer
	// TODO: error is ignored for now.
	xml.EscapeText(&b, []byte(s))
	return b.String()
}
