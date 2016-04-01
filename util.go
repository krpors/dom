package dom

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

// TODO: namespaces for prefixes must be predeclared somehow before serializing
// or else it must generate an error.

// PrintTree prints the whole tree starting from the given Node 'node' to the
// writer w.
func PrintTree(start Node, w io.Writer) {
	var traverse func(Node, string)

	traverse = func(n Node, indent string) {
		fmt.Fprintf(w, "%s%v\n", indent, n)
		for _, child := range n.GetChildNodes() {
			traverse(child, indent+"  ")
		}
	}
	traverse(start, "")
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
//
// This function does not take the xmlns attributes into account. In other words, if an
// Element is created without the Document's CreateElementNS() method, it will NOT have a
// namespace. Even after an xmlns attribute is added. Xerces does it like this too.
func getElementsBy(parent Node, namespaceURI, tagname string, includeNamespace bool) []Element {
	var elements []Element

	var traverse func(n Node)
	traverse = func(n Node) {
		for _, child := range n.GetChildNodes() {
			// only check elements:
			if elem, ok := child.(Element); ok {
				if includeNamespace && elem.GetLocalName() == tagname && elem.GetNamespaceURI() == namespaceURI {
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

// GetPreviousSibling gets the previous sibling, using the specified Node as
// a reference. If there is no previous sibling (i.e. this NOde is the first), it
// will return nil. If the Node has no parent this will return nil.
//
// This implementation is rather naive: it iterates through all child elements
// of the parent, checks the position of this element, then uses that index - 1.
func getPreviousSibling(node Node) Node {
	if node.GetParentNode() == nil {
		return nil
	}

	siblings := node.GetParentNode().GetChildNodes()
	for i, n := range siblings {
		if n == node && i > 0 {
			return siblings[i-1]
		}
	}

	return nil
}

// getNextSibling gets the next sibling, using the specified Node as a reference.
// If there is no next sibling (i.e. this Node is the last), it will return nil.
// If the Node has no parent, this will return nil.
//
// This implementation is rather naive: it iterates through all child elements
// of the parent, checks the position of this Node, then uses that index + 1.
func getNextSibling(node Node) Node {
	if node.GetParentNode() == nil {
		return nil
	}

	siblings := node.GetParentNode().GetChildNodes()
	for i, n := range siblings {
		if n == node && i < len(siblings)-1 {
			return siblings[i+1]
		}
	}

	return nil
}
