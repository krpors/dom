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

func newText(owner Document) Text {
	t := &domText{}
	t.ownerDocument = owner
	return t
}

func (dt *domText) GetNodeName() string {
	return "#text"
}

func (dt *domText) GetNodeType() NodeType {
	return TextNode
}

// NodeValue returns the same as GetData, the content of the text node.
func (dt *domText) GetNodeValue() string {
	return dt.GetText()
}

func (dt *domText) GetLocalName() string {
	return ""
}

func (dt *domText) GetChildNodes() []Node {
	return nil
}

func (dt *domText) GetParentNode() Node {
	return dt.parentNode
}

func (dt *domText) GetFirstChild() Node {
	return nil
}

func (dt *domText) GetLastChild() Node {
	return nil
}

func (dt *domText) GetAttributes() NamedNodeMap {
	return nil
}

func (dt *domText) GetOwnerDocument() Document {
	return dt.ownerDocument
}

func (dt *domText) AppendChild(child Node) error {
	return fmt.Errorf("%v: %v does not allow child nodes", ErrorHierarchyRequest, TextNode)
}

func (dt *domText) RemoveChild(oldChild Node) (Node, error) {
	return nil, fmt.Errorf("%v: %v does not allow child nodes - nothing to remove", ErrorHierarchyRequest, TextNode)
}
func (dt *domText) ReplaceChild(newChild, oldChild Node) (Node, error) {
	return nil, fmt.Errorf("%v: %v does not allow child nodes - nothing to replace", ErrorHierarchyRequest, TextNode)
}
func (dt *domText) InsertBefore(newChild, refChild Node) (Node, error) {
	return nil, fmt.Errorf("%v: %v does not allow child nodes - nothing to insert", ErrorHierarchyRequest, TextNode)
}

func (dt *domText) HasChildNodes() bool {
	return false
}

func (dt *domText) GetPreviousSibling() Node {
	return getPreviousSibling(dt)
}

func (dt *domText) GetNextSibling() Node {
	return getNextSibling(dt)
}

// GetNamespaceURI returns an empty string for Text nodes.
func (dt *domText) GetNamespaceURI() string {
	return ""
}

// GetNamespacePrefix returns an empty string for Text nodes.
func (dt *domText) GetNamespacePrefix() string {
	return ""
}

func (dt *domText) LookupPrefix(namespace string) string {
	return ""
}

func (dt *domText) LookupNamespaceURI(pfx string) string {
	return ""
}

func (dt *domText) GetTextContent() string {
	return dt.GetNodeValue()
}

func (dt *domText) SetTextContent(content string) {
	dt.SetText(content)
}

// Private functions:
func (dt *domText) setParentNode(parent Node) {
	dt.parentNode = parent
}

// Text specifics:

// GetText returns the character data of this text node, unescaped.
func (dt *domText) GetText() string {
	return dt.data
}

// SetText sets the character data of the XML node. The data can be unescaped
// XML, since GetText() will take care of conversion.
func (dt *domText) SetText(data string) {
	dt.data = data
}

// IsElementContentWhitespace returns true when the Text node contains ignorable
// whitespace, like any combinations of \t, \n, \r and space characters.
func (dt *domText) IsElementContentWhitespace() bool {
	return strings.TrimSpace(dt.GetText()) == ""
}

func (dt *domText) String() string {
	maxlen := 30
	var d string
	if len(dt.data) > maxlen {
		d = strings.TrimSpace(dt.data[0:maxlen] + " [...]")
	}
	return fmt.Sprintf("%s: '%s'", dt.GetNodeType(), d)
}
