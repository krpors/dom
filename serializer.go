package dom

import (
	"fmt"
	"io"
	"strings"
)

// Serializes defines the type that can be used to serialize a Node + its children.
type Serializer struct {
	Configuration Configuration // Serializer's configuration.
}

// NewSerializer creates a new Serializer using the default configuration.
func NewSerializer() *Serializer {
	s := &Serializer{}
	s.Configuration = NewConfiguration()
	return s
}

func (s *Serializer) nodeContainsTextOnly(n Node) bool {
	if !n.HasChildNodes() {
		return false
	}

	for _, c := range n.GetChildNodes() {
		if c.GetNodeType() != TextNode {
			return false
		}
	}
	return true
}

// Serialize write the node plus its children to the writer w.
func (s *Serializer) Serialize(node Node, w io.Writer) {
	// Must define the function here so we can refer to ourselves in
	// the traverse function.
	var traverse func(n Node, indent string)

	if !s.Configuration.OmitXMLDeclaration {
		fmt.Fprintf(w, "%s", XMLDeclaration)
		if s.Configuration.PrettyPrint {
			fmt.Fprintln(w)
		}
	}

	traverse = func(n Node, indent string) {
		switch t := n.(type) {
		case Element:
			if s.Configuration.PrettyPrint {
				fmt.Fprintf(w, "%s", indent)
			}

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

			// Add a newline after element start, if pretty printing, and the node doesn't contain text only nodes.
			if s.Configuration.PrettyPrint && !s.nodeContainsTextOnly(n) {
				fmt.Fprintf(w, "\n")
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
			// Don't indent the first element when the first node is a DocumentNode.
			if n.GetNodeType() == DocumentNode {
				traverse(node, "")
			} else {
				traverse(node, indent+s.Configuration.IndentCharacter)
			}
		}

		switch t := n.(type) {
		case Element:
			if t.HasChildNodes() {
				if s.Configuration.PrettyPrint && !s.nodeContainsTextOnly(n) {
					fmt.Fprintf(w, "%s", indent)
				}
				fmt.Fprintf(w, "</%s>\n", t.GetTagName())
			}
		}
	}

	traverse(node, "")
}
