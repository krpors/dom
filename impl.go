package dom

import (
	"errors"
	"fmt"
)

var (
	// ErrorHierarchyRequest is the error which can be returned when the node
	// is of a type that does not allow children, if the node to append to is
	// one of this node's ancestors or this node itself, or if this node is of
	// type Document and the DOM application attempts to append a second
	// DocumentType or Element node.
	ErrorHierarchyRequest = errors.New("HIERARCHY_REQUEST_ERR: an attempt was made to insert a node where it is not permitted")

	// ErrorInvalidCharacter is returned when an invalid character is used for
	// for example an element or attribute name.
	ErrorInvalidCharacter = errors.New("INVALID_CHARACTER_ERR: an invalid or illegal XML character is specified")

	// ErrorNotSupported is returned when this implementation does not support
	// the requested operation or object.
	ErrorNotSupported = errors.New("NOT_SUPPORTED_ERR: this implementation does not support the requested type of object or operation")
)

/*
TODO: as I suspected, making an attempt at inheritance by embedding interfaces
in other interfaces works out poorly. Especially when setting parent nodes when
appending children to nodes. The pointers are not done correctly. One BIG fix
should be the overhaul of the types.

Define all the necessary interfaces (interface embedding can be done though),
but all the concrete types must manually implement ALL interface methods. This
will lead to copy/pasting, but it gets the job done.

Example:

type Node interface {
	// Node is abstract is it can get, so there will be no implementation of it.
	// ... All functions of Node
}

type Element interface {
	Node // <-- embedded interface
	// ... All other functions of Element
}

Concrete types:

type domElement struct {
	// All members which are necessary, such as attributes, child nodes,
	// parent node, owner document, etc. etc.
}
// Add all functions of Node and Element here.

type domDocument struct {
	// All member vars, such as child, processing instructions...
}
// Add all functions of Node and Document here

*/

type domNamedNodeMap struct {
	Attrs map[string]Node
}

func newNamedNodeMap() NamedNodeMap {
	nnm := &domNamedNodeMap{}
	nnm.Attrs = make(map[string]Node)
	return nnm
}

func (dn *domNamedNodeMap) GetNamedItem(s string) Node {
	return dn.Attrs[s]
}

func (dn *domNamedNodeMap) SetNamedItem(node Node) error {
	// Node must be an AttributeNode, or else it will return a hierarchy error.
	if node.NodeType() != AttributeNode {
		return ErrorHierarchyRequest
	}
	dn.Attrs[node.NodeName()] = node
	return nil
}

func (dn *domNamedNodeMap) GetItems() map[string]Node {
	return dn.Attrs
}

func (dn *domNamedNodeMap) Length() int {
	return len(dn.Attrs)
}

func (dn *domNamedNodeMap) String() string {
	s := "{"
	for name, val := range dn.Attrs {
		s += fmt.Sprintf("%v: %v, ", name, val.NodeValue())
	}
	s += "}"
	return s
}

//================================================================================

type domNode struct {
	localName     string
	nodeValue     string
	nodeType      NodeType
	nodes         []Node
	parentNode    Node
	firstChild    Node
	attributes    NamedNodeMap
	ownerDocument Document
	namespaceURI  string
}

func newNode() Node {
	node := &domNode{}
	node.attributes = newNamedNodeMap()
	return node
}

func (dn *domNode) NodeName() string {
	// TODO: this is incorrect, see "Definition group NodeType" in the spec
	return dn.nodeType.String()
}

// NodeValue returns an empty string. A Node by itself does not have a value,
// unless implemented/embedded by other DOM objects, such as Element. The
// specification mention for instance that Document, Element, Entity - and more -
// should return null on a call to  NodeValue(). Go doesn't allow nils like this,
// so we return an empty string for now.
func (dn *domNode) NodeValue() string {
	return ""
}

func (dn *domNode) NodeType() NodeType {
	return dn.nodeType
}

func (dn *domNode) LocalName() string {
	return dn.localName
}

func (dn *domNode) NodeList() []Node {
	return dn.nodes
}

func (dn *domNode) ParentNode() Node {
	return dn.parentNode
}

func (dn *domNode) FirstChild() Node {
	return dn.NodeList()[0]
}

func (dn *domNode) GetAttributes() NamedNodeMap {
	return dn.attributes
}

func (dn *domNode) OwnerDocument() Document {
	return dn.ownerDocument
}

func (dn *domNode) AppendChild(node Node) error {
	node.setParentNode(dn)
	dn.nodes = append(dn.nodes, node)
	return nil
}

func (dn *domNode) HasChildNodes() bool {
	return len(dn.nodes) > 0
}

func (dn *domNode) NamespaceURI() string {
	return dn.namespaceURI
}

func (dn *domNode) setParentNode(node Node) {
	dn.parentNode = node
}

func (dn *domNode) setNodeValue(s string) {
	dn.nodeValue = s
}

func (dn *domNode) setNodeType(t NodeType) {
	dn.nodeType = t
}

func (dn *domNode) setOwnerDocument(d Document) {
	dn.ownerDocument = d
}

func (dn *domNode) setNamespaceURI(namespaceURI string) {
	dn.namespaceURI = namespaceURI
}

//================================================================================

type domDocumentType struct {
	Node

	Name     string
	PublicID string
	SystemID string
}

func newDocumentType() DocumentType {
	dt := &domDocumentType{}
	dt.Node = &domNode{}
	dt.setNodeType(DocumentTypeNode)
	return dt
}

func (dt *domDocumentType) GetName() string {
	return dt.Name
}
func (dt *domDocumentType) GetPublicID() string {
	return dt.PublicID
}
func (dt *domDocumentType) GetSystemID() string {
	return dt.SystemID
}

//================================================================================

type domDocument struct {
	Node
}

// NewDocument creates a new document.
func NewDocument() Document {
	d := &domDocument{}
	d.Node = &domNode{}
	d.setNodeType(DocumentNode)
	return d
}

func (dd *domDocument) CreateElement(tagName string) (Element, error) {
	elem := newElement()
	elem.SetTagName(tagName)
	elem.setParentNode(dd)
	elem.setOwnerDocument(dd)
	return elem, nil
}

func (dd *domDocument) CreateElementNS(namespaceURI, tagName string) (Element, error) {
	elem, err := dd.CreateElement(tagName)
	if err != nil {
		return nil, err
	}
	elem.setNamespaceURI(namespaceURI)
	return elem, nil
}

func (dd *domDocument) CreateTextNode(t string) Text {
	text := newText()
	text.setOwnerDocument(dd)
	text.SetData(t)
	return text
}

func (dd *domDocument) CreateAttribute(name string) (Attr, error) {
	attr := newAttr(name)
	attr.setOwnerDocument(dd)
	return attr, nil
}

// 'Override' the AppendChild() function from the Node interface. One child can
// be appended when the node list is empty. The first child of the document will
// be the document element. Subsequent calls will result in an error.
func (dd *domDocument) AppendChild(n Node) error {
	if len(dd.NodeList()) <= 0 {
		n.setParentNode(dd)
		dd.Node.AppendChild(n)
		return nil
	}
	return ErrorHierarchyRequest
}

func (dd *domDocument) GetDocumentElement() Element {
	if len(dd.NodeList()) == 1 {
		node := dd.FirstChild()
		if node.NodeType() == ElementNode {
			elem := node.(Element)
			return elem
		}
	}
	return nil
}

func (dd *domDocument) String() string {
	return fmt.Sprintf("%s", dd.NodeType())
}

//================================================================================

type domElement struct {
	Node

	TagName string
}

func newElement() Element {
	e := &domElement{}
	e.Node = newNode()
	e.setNodeType(ElementNode)
	return e
}

func (de *domElement) SetTagName(name string) {
	de.TagName = name
}

func (de *domElement) GetTagName() string {
	return de.TagName
}

func (de *domElement) SetAttribute(name, value string) {
	attr := newAttr(name)
	attr.setParentNode(de)
	attr.SetValue(value)
	de.Node.GetAttributes().SetNamedItem(attr)
}

func (de *domElement) GetAttribute(name string) string {
	return ""
}

func (de *domElement) String() string {
	return fmt.Sprintf("%s, <%s>, ns=%s, attrs=%v",
		de.NodeType(), de.GetTagName(), de.NamespaceURI(), de.GetAttributes())
}

//================================================================================

type domAttr struct {
	Node

	Name         string
	Specified    bool
	Value        string
	OwnerElement Element
}

func newAttr(name string) Attr {
	a := &domAttr{}
	a.Node = &domNode{}
	a.setNodeType(AttributeNode)
	a.Name = name
	return a
}

// NodeName is an override from Node. As per the spec, the NodeName() function
// should return the same thing as GetName().
func (da *domAttr) NodeName() string {
	return da.Name
}

// Identical to NodeName().
func (da *domAttr) GetName() string {
	return da.Name
}

// NodeValue is an override from Node.
func (da *domAttr) NodeValue() string {
	return da.Value
}

// Identical to NodeValue()
func (da *domAttr) GetValue() string {
	return da.Value
}
func (da *domAttr) SetValue(val string) {
	da.Specified = true
	da.Value = val
}

func (da *domAttr) IsSpecified() bool {
	return da.Specified
}

func (da *domAttr) GetOwnerElement() Element {
	return da.OwnerElement
}

func (da *domAttr) setOwnerElement(e Element) {
	da.OwnerElement = e
}

//================================================================================

type domText struct {
	Node

	Data string
}

func newText() Text {
	t := &domText{}
	t.Node = &domNode{}
	t.setNodeType(TextNode)
	return t
}

func (dt *domText) NodeName() string {
	return "#text"
}

// NodeValue is an override from the embedded Node. It returns the character
// data, identical to the GetData() function.
func (dt *domText) NodeValue() string {
	return dt.Data
}

func (dt *domText) SetData(s string) {
	dt.Data = s
}

// GetData return the character data. Identical to NodeValue().
func (dt *domText) GetData() string {
	return dt.Data
}

func (dt *domText) String() string {
	return fmt.Sprintf("%s, %s, %s", dt.NodeType(), dt.NodeName(), dt.GetData())
}
