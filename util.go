package dom

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

// TODO: namespaces for prefixes must be predeclared somehow before serializing
// or else it must generate an error.

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
			fmt.Fprintf(w, "<%s", t.GetTagName())
			// Add any attributes
			if t.GetAttributes() != nil {
				for _, val := range t.GetAttributes().GetItems() {
					attr := val.(Attr)
					fmt.Fprintf(w, " %s=\"%s\"", attr.GetNodeName(), attr.GetNodeValue())
				}
			}
			if t.HasChildNodes() {
				fmt.Fprintf(w, ">")
			} else {
				fmt.Fprintf(w, "/>")
			}
		case Text:
			fmt.Fprintf(w, "%s", escape(t.GetData()))
		case Comment:
			fmt.Fprintf(w, "<!-- %s -->", t.GetComment())
		case ProcessingInstruction:
			// TODO: proper serialization of target/data. Must include valid chars etc.
			// Also, if target/data contains '?>', generate a fatal error.
			fmt.Fprintf(w, "<?%v %v?>", t.GetTarget(), t.GetData())
		}

		// For each child node, call traverse() again.
		for _, node := range n.GetChildNodes() {
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

// getElementsBy finds descandant elements with the given (optional) namespaceURI and tagname.
// When the 'includeNamespace' is set to true, the namespace URI is explicitly checked for
// equality. If false, no namespace check will be done. The elements are returned as a 'live'
// slice. The nodes are searched using the given parent node, which must not be nil.
func getElementsBy(parent Node, namespaceURI, tagname string, includeNamespace bool) []Element {
	if parent == nil {
		panic("parent node cannot be nil")
	}

	var elements []Element

	var traverse func(n Node)
	traverse = func(n Node) {
		for _, child := range n.GetChildNodes() {
			// only check elements:
			if elem, ok := child.(Element); ok {
				if includeNamespace && elem.GetNodeName() == tagname && elem.GetNamespaceURI() == namespaceURI {
					// include namespace equality, if chosen.
					elements = append(elements, elem)
				} else if !includeNamespace && elem.GetNodeName() == tagname {
					// do not include namespace equality, just the tagname
					elements = append(elements, elem)
				}

			}

			traverse(child)
		}
	}
	traverse(parent)

	return elements
}
