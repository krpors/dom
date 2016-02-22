package dom

import (
	"fmt"
	"strings"
)

// domComment. We don't 'inherit' from CharacterData, that's a bit too convoluted...
// Maybe we'll implement that some other time.
type domComment struct {
	localName     string
	parentNode    Node
	ownerDocument Document

	// Comment specific things
	comment string
}

func newComment() Comment {
	t := &domComment{}
	return t
}

func (dc *domComment) NodeName() string {
	return "#comment"
}

func (dc *domComment) NodeType() NodeType {
	return CommentNode
}

// NodeValue returns the same as GetData, the content of the text node.
func (dc *domComment) NodeValue() string {
	return dc.GetComment()
}

func (dc *domComment) LocalName() string {
	// TODO: huh? for text?
	return ""
}

func (dc *domComment) NodeList() []Node {
	return nil
}

func (dc *domComment) ParentNode() Node {
	return dc.parentNode
}

func (dc *domComment) FirstChild() Node {
	return nil
}

func (dc *domComment) GetAttributes() NamedNodeMap {
	return nil
}

func (dc *domComment) OwnerDocument() Document {
	return dc.ownerDocument
}

func (dc *domComment) AppendChild(child Node) error {
	return fmt.Errorf("%v: %v does not allow children", ErrorHierarchyRequest, dc.NodeType())
}

func (dc *domComment) HasChildNodes() bool {
	return false
}

func (dc *domComment) NamespaceURI() string {
	return ""
}

// Private functions:
func (dc *domComment) setParentNode(parent Node) {
	dc.parentNode = parent
}

func (dc *domComment) setOwnerDocument(d Document) {
	dc.ownerDocument = d
}

func (dc *domComment) setNamespaceURI(uri string) {
	// no-op
}

// Text specifics:

// GetComment returns the comment content.
func (dc *domComment) GetComment() string {
	return dc.comment
}

// SetComment sets the character comment data of the XML node.
func (dc *domComment) SetComment(comment string) {
	dc.comment = comment
}

func (dc *domComment) String() string {
	return fmt.Sprintf("%s: '%s'", dc.NodeType(), strings.TrimSpace(dc.comment))
}
