package dom

import (
	"fmt"
	"strings"
)

// domText. We don't 'inherit' from CharacterData, that's a bit too convoluted...
// Maybe we'll implement that some other time.
type domText struct {
	localName     string
	parentNode    Node
	ownerDocument Document

	// Text specific things
	data string
}

func newText() Text {
	t := &domText{}
	return t
}

func (dt *domText) NodeName() string {
	return "#text"
}

func (dt *domText) NodeType() NodeType {
	return TextNode
}

// NodeValue returns the same as GetData, the content of the text node.
func (dt *domText) NodeValue() string {
	return dt.GetData()
}

func (dt *domText) LocalName() string {
	// TODO: huh? for text?
	return ""
}

func (dt *domText) NodeList() []Node {
	return nil
}

func (dt *domText) ParentNode() Node {
	return dt.parentNode
}

func (dt *domText) FirstChild() Node {
	return nil
}

func (dt *domText) GetAttributes() NamedNodeMap {
	return nil
}

func (dt *domText) OwnerDocument() Document {
	return dt.ownerDocument
}

func (dt *domText) AppendChild(child Node) error {
	return fmt.Errorf("%v: %v does not allow children", ErrorHierarchyRequest, dt.NodeType())
}

func (dt *domText) HasChildNodes() bool {
	return false
}

func (dt *domText) NamespaceURI() string {
	return ""
}

// Private functions:
func (dt *domText) setParentNode(parent Node) {
	dt.parentNode = parent
}

func (dt *domText) setOwnerDocument(d Document) {
	dt.ownerDocument = d
}

func (dt *domText) setNamespaceURI(uri string) {
	// no-op
}

// Text specifics:

// GetData returns the character data of this text node, unescaped.
func (dt *domText) GetData() string {
	return dt.data
}

// SetData sets the character data of the XML node. The data can be unescaped
// XML, since GetData() will take care of conversion.
func (dt *domText) SetData(data string) {
	dt.data = data
}

func (dt *domText) String() string {
	return fmt.Sprintf("%s: '%s'", dt.NodeType(), strings.TrimSpace(dt.data))
}
