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
	attrName     string
	attrValue    string
}

func newAttr() Attr {
	a := &domAttr{}
	return a
}

func (da *domAttr) GetNodeName() string {
	return da.attrName
}

func (da *domAttr) GetNodeType() NodeType {
	return AttributeNode
}

// NodeValue should return null/nil for Element types like the spec says,
// but Go does not permit nil strings which are not pointers. So for now we
// just return an empty string at all times.
func (da *domAttr) GetNodeValue() string {
	return da.attrValue
}

func (da *domAttr) GetLocalName() string {
	// TODO: what?
	return ""
}

// GetChildNodes() returns an empty list of nodes for the Attr type.
func (da *domAttr) GetChildNodes() []Node {
	return []Node{}
}

// ParentNode returns nil, since the spec says Attr objects cannot have parents.
func (da *domAttr) GetParentNode() Node {
	return nil
}

// FirstChild will return nil, since Attr objects cannot contain children.
func (da *domAttr) GetFirstChild() Node {
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

// HasChildNodes returns false since Attr objects do not contain children.
func (da *domAttr) HasChildNodes() bool {
	return false
}

func (da *domAttr) GetNamespaceURI() string {
	return da.namespaceURI
}

func (da *domAttr) GetNamespacePrefix() string {
	// TODO: namespace prefix of an attr.
	return ""
}

// Private functions:
func (da *domAttr) setParentNode(parent Node) {
	// no-op
}

func (da *domAttr) setOwnerDocument(d Document) {
	da.ownerDocument = d
}

func (da *domAttr) setNamespaceURI(uri string) {
	da.namespaceURI = uri
}

// Attr specific functions:

func (da *domAttr) setName(name string) {
	da.attrName = name
}

func (da *domAttr) GetName() string {
	return da.attrName
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
	// TODO: this
	return fmt.Sprintf("%v, %v=%v", da.GetNodeType(), da.attrName, da.attrValue)
}
