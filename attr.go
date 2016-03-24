package dom

import (
	"fmt"
)

type domAttr struct {
	localName     string
	attributes    NamedNodeMap
	ownerDocument Document
	namespaceURI  string

	// Attr specific things:
	ownerElement Element
	attrName     XMLName
	attrValue    string
}

func newAttr(owner Document, name string, namespaceURI string) Attr {
	a := &domAttr{}
	a.ownerDocument = owner
	a.attrName = XMLName(name)
	a.namespaceURI = namespaceURI
	return a
}

func (da *domAttr) GetNodeName() string {
	return string(da.attrName)
}

func (da *domAttr) GetNodeType() NodeType {
	return AttributeNode
}

func (da *domAttr) GetNodeValue() string {
	return da.attrValue
}

func (da *domAttr) GetLocalName() string {
	return da.attrName.GetLocalPart()
}

// GetChildNodes() returns an empty list of nodes for the Attr type.
func (da *domAttr) GetChildNodes() []Node {
	return []Node{}
}

// ParentNode returns nil, since the spec says Attr objects cannot have parents.
func (da *domAttr) GetParentNode() Node {
	return nil
}

// GetLastChild will return nil, since Attr objects cannot contain children.
func (da *domAttr) GetFirstChild() Node {
	return nil
}

// GetFirstChild will return nil, since Attr objects cannot contain children.
func (da *domAttr) GetLastChild() Node {
	return nil
}

// GetAttributes will return nil, since this will be called on an instance of Attr.
// Only Element objects can have attributes.
func (da *domAttr) GetAttributes() NamedNodeMap {
	return nil
}

// OwnerDocument returns the owner document of this Attr.
func (da *domAttr) GetOwnerDocument() Document {
	return da.ownerDocument
}

// AppendChild returns a hierarchy error for Attr objects.
func (da *domAttr) AppendChild(child Node) error {
	return fmt.Errorf("%v: attributes do not allow children", ErrorHierarchyRequest)
}

func (da *domAttr) RemoveChild(oldChild Node) (Node, error) {
	return nil, ErrorHierarchyRequest
}
func (da *domAttr) ReplaceChild(newChild, oldChild Node) (Node, error) {
	return nil, ErrorHierarchyRequest
}
func (da *domAttr) InsertBefore(newChild, refChild Node) (Node, error) {
	return nil, ErrorHierarchyRequest
}

// HasChildNodes returns false since Attr objects do not contain children.
func (da *domAttr) HasChildNodes() bool {
	return false
}

// GetPreviousSibling always returns nil for Attr nodes.
func (da *domAttr) GetPreviousSibling() Node {
	return nil
}

// GetNextSibling always returns nil for Attr nodes.
func (da *domAttr) GetNextSibling() Node {
	return nil
}

func (da *domAttr) GetNamespaceURI() string {
	return da.namespaceURI
}

func (da *domAttr) GetNamespacePrefix() string {
	return da.attrName.GetPrefix()
}

func (da *domAttr) LookupPrefix(namespace string) string {
	if da.GetOwnerElement() != nil {
		return da.GetOwnerElement().LookupPrefix(namespace)
	}
	return ""
}

func (da *domAttr) LookupNamespaceURI(pfx string) string {
	if da.GetOwnerElement() != nil {
		return da.GetOwnerElement().LookupNamespaceURI(pfx)
	}
	return ""
}

func (da *domAttr) GetTextContent() string {
	return da.GetValue()
}

func (da *domAttr) SetTextContent(content string) {
	da.SetValue(content)
}

// Private functions:
func (da *domAttr) setParentNode(parent Node) {
	// no-op
}

func (da *domAttr) GetName() string {
	return da.GetNodeName()
}

func (da *domAttr) IsSpecified() bool {
	// TODO: what?
	return true
}

// GetOwnerElement returns the Element that owns this Attr, or nil if the attribute is
// not in use.
func (da *domAttr) GetOwnerElement() Element {
	return da.ownerElement
}

func (da *domAttr) GetValue() string {
	return da.attrValue
}

func (da *domAttr) SetValue(val string) {
	da.attrValue = val
}

func (da *domAttr) setOwnerElement(owner Element) {
	da.ownerElement = owner
}

func (da *domAttr) String() string {
	return fmt.Sprintf("%v, %v=%v", da.GetNodeType(), da.attrName, da.attrValue)
}
