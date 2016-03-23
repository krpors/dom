package dom

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
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
			// Contains only whitespaces?
			if strings.TrimSpace(t.GetText()) == "" {
				fmt.Fprintf(w, "%s", t.GetText())
			} else {
				fmt.Fprintf(w, "%s", escape(t.GetText()))
			}
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
