package dom

import (
	"fmt"
	"io"
	"strings"
)

// Serializer defines the type that can be used to serialize a Node + its children. The struct configuration
// can be used to control the output of the serialization to a certain degree.
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

func isNamespaceAttrDefined(attr Attr) bool {
	if strings.HasPrefix(attr.GetName(), "xmlns:") {
		if element, ok := attr.GetOwnerElement().GetParentNode().(Element); ok {
			// check if the namespace is already declared, but NOT IN THE OWNER ELEMENT!!!
			// of else it will always be true :|
			ns, foundpfx := element.LookupPrefix(attr.GetValue())
			pfx, foundns := element.LookupNamespaceURI(attr.GetLocalName())
			if foundpfx && foundns && ns == attr.GetLocalName() && pfx == attr.GetValue() {
				return true
			}
		}
	} else if attr.GetName() == "xmlns" {
		// FIXME: owner element can be nil in case of root element
		if element, ok := attr.GetOwnerElement().GetParentNode().(Element); ok {
			return element.IsDefaultNamespace(attr.GetValue())
		}

	}

	return false
}

// Serialize writes the node plus its children to the writer w. The Serializer does not do any
// specific mutations on the given Node to serialize, i.e. it will write it as-is. No normalizations,
// alterations etc are done.
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
			// When pretty printing, indent the <element> string with the specified amount of indent chars.
			if s.Configuration.PrettyPrint {
				fmt.Fprintf(w, "%s", indent)
			}
			// In any case, write the tagname <x>.
			fmt.Fprintf(w, "<%s", t.GetTagName())

			// if this is the document element, make sure to write the namespace declaration,
			// and do not try to look it up using Lookup*
			if t.GetParentNode().GetNodeType() == DocumentNode {
				if t.GetNamespacePrefix() != "" {
					fmt.Fprintf(w, ` xmlns:%s="%s"`, t.GetNamespacePrefix(), t.GetNamespaceURI())
				} else if t.GetNamespacePrefix() == "" && t.GetNamespaceURI() != "" {
					fmt.Fprintf(w, ` xmlns="%s"`, t.GetNamespaceURI())
				}
			} else {
				// in other cases, look stuff up in ancestors
				if t.GetNamespacePrefix() != "" {
					_, found := t.GetParentNode().LookupNamespaceURI(t.GetNamespacePrefix())
					if !found && t.GetNamespaceURI() != "" {
						fmt.Fprintf(w, ` xmlns:%s="%s"`, t.GetNamespacePrefix(), t.GetNamespaceURI())
					}
				}
			}

			// Add any attributes
			if t.GetAttributes() != nil {
				for _, val := range t.GetAttributes().GetItems() {
					// TODO: configurable sort on attributes to make things more deterministic
					attr := val.(Attr)
					if !isNamespaceAttrDefined(attr) {
						fmt.Fprintf(w, " %s=\"%s\"", attr.GetNodeName(), attr.GetNodeValue())
					}
				}
			}
			// If the current element has any children, do not end the element, e.g. <element>
			if t.HasChildNodes() {
				fmt.Fprintf(w, ">")
			} else {
				// Write the element as <element/>, because no elements follow.
				fmt.Fprintf(w, "/>")
			}

			// Add a newline after element start, if pretty printing, and the node doesn't contain text only nodes.
			if s.Configuration.PrettyPrint && !s.nodeContainsTextOnly(n) {
				fmt.Fprintf(w, "\n")
			}

		case Text:
			// Contains only whitespaces? If so, write the text as-is.
			if strings.TrimSpace(t.GetText()) == "" {
				fmt.Fprintf(w, "%s", t.GetText())
			} else {
				// Else escape any text where necessary.
				fmt.Fprintf(w, "%s", escape(t.GetText()))
			}
		case Comment:
			// When pretty printing, indent the comment with the indent level.
			if s.Configuration.PrettyPrint {
				fmt.Fprintf(w, "%s", indent)
			}
			fmt.Fprintf(w, "<!-- %s -->\n", t.GetComment())
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
				// Serialize this child. Call traverse again with an increased indent character.
				traverse(node, indent+s.Configuration.IndentCharacter)
			}
		}

		// Check if and how we should write an element ending: </element>
		switch t := n.(type) {
		case Element:
			if t.HasChildNodes() {
				// Are we pretty printing, and the Element does not contain text only nodes? Then just write the
				// indent characters. Example:
				//
				// <element>
				//   <child>
				//     <other/>
				//   </child> <== indent character at this point.
				// </element>
				if s.Configuration.PrettyPrint && !s.nodeContainsTextOnly(n) {
					fmt.Fprintf(w, "%s", indent)
				}
				// In any case, write the 'end element'.
				fmt.Fprintf(w, "</%s>", t.GetTagName())
				// When pretty printing, be sure to write a trailing newline.
				if s.Configuration.PrettyPrint {
					fmt.Fprint(w, "\n")
				}
			}
		}
	}

	traverse(node, "")
}
