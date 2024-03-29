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

func newComment(owner Document) Comment {
	t := &domComment{}
	t.ownerDocument = owner
	return t
}

func (dc *domComment) GetNodeName() string {
	return "#comment"
}

func (dc *domComment) GetNodeType() NodeType {
	return CommentNode
}

// NodeValue returns the same as GetData, the content of the text node.
func (dc *domComment) GetNodeValue() string {
	return dc.GetComment()
}

func (dc *domComment) GetLocalName() string {
	// TODO: huh? for text?
	return ""
}

func (dc *domComment) GetChildNodes() []Node {
	return nil
}

func (dc *domComment) GetParentNode() Node {
	return dc.parentNode
}

func (dc *domComment) GetFirstChild() Node {
	return nil
}

func (dc *domComment) GetLastChild() Node {
	return nil
}

func (dc *domComment) GetAttributes() NamedNodeMap {
	return nil
}

func (dc *domComment) HasAttributes() bool {
	return false
}

func (dc *domComment) GetOwnerDocument() Document {
	return dc.ownerDocument
}

func (dc *domComment) AppendChild(child Node) error {
	return fmt.Errorf("%v: %v does not allow children", ErrorHierarchyRequest, dc.GetNodeType())
}

func (dc *domComment) RemoveChild(oldChild Node) (Node, error) {
	return nil, ErrorHierarchyRequest
}
func (dc *domComment) ReplaceChild(newChild, oldChild Node) (Node, error) {
	return nil, ErrorHierarchyRequest
}
func (dc *domComment) InsertBefore(newChild, refChild Node) (Node, error) {
	return nil, ErrorHierarchyRequest
}

func (dc *domComment) HasChildNodes() bool {
	return false
}

func (dc *domComment) GetPreviousSibling() Node {
	return getPreviousSibling(dc)
}

func (dc *domComment) GetNextSibling() Node {
	return getNextSibling(dc)
}

func (dc *domComment) GetNamespaceURI() string {
	return ""
}

// GetNamespacePrefix returns an empty string for comments.
func (dc *domComment) GetNamespacePrefix() string {
	return ""
}

func (dc *domComment) LookupPrefix(namespace string) (string, bool) {
	if dc.GetParentNode() != nil {
		return dc.GetParentNode().LookupPrefix(namespace)
	}
	return "", false
}

func (dc *domComment) LookupNamespaceURI(pfx string) (string, bool) {
	if dc.GetParentNode() != nil {
		return dc.GetParentNode().LookupNamespaceURI(pfx)
	}
	return "", false
}

func (dc *domComment) IsDefaultNamespace(namespace string) bool {
	// TODO ?
	return false
}

func (dc *domComment) GetTextContent() string {
	return dc.GetNodeValue()
}

func (dc *domComment) SetTextContent(content string) {
	dc.SetComment(content)
}

func (dc *domComment) CloneNode(deep bool) Node {
	cloneComment, err := dc.ownerDocument.CreateComment(dc.comment)
	if err != nil {
		panic("CreateComment returned error, but was unexpected at this point")
	}
	return cloneComment
}

func (dc *domComment) ImportNode(n Node, deep bool) Node {
	return importNode(dc.ownerDocument, n, deep)
}

// Private functions:
func (dc *domComment) setParentNode(parent Node) {
	dc.parentNode = parent
}

func (dc *domComment) setOwnerDocument(doc Document) {
	dc.ownerDocument = doc
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
	return fmt.Sprintf("%s: '%s'", dc.GetNodeType(), strings.TrimSpace(dc.comment))
}
