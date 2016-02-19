package dom

import (
	"errors"
	"fmt"
	"strings"
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
)

// NodeType defines the types of nodes which exist in the DOM.
type NodeType uint8

// Enumeration of all types of Nodes in the DOM.
const (
	ElementNode NodeType = iota
	AttributeNode
	TextNode
	CDATASectionNode
	EntityReferenceNode
	EntityNode
	ProcessingInstructionNode
	CommentNode
	DocumentNode
	DocumentTypeNode
	DocumentFragmentNode
)

// String returns the string representation of the NodeType, using the default
// representation by the W3 specification.
func (n NodeType) String() string {
	switch n {
	case ElementNode:
		return "ELEMENT_NODE"
	case AttributeNode:
		return "ATTRIBUTE_NODE"
	case TextNode:
		return "TEXT_NODE"
	case CDATASectionNode:
		return "CDATA_SECTION_NODE"
	case EntityReferenceNode:
		return "ENTITY_REFERENCE_NODE"
	case EntityNode:
		return "ENTITY_NODE"
	case ProcessingInstructionNode:
		return "PROCESSING_INSTRUCTION_NODE"
	case CommentNode:
		return "COMMENT_NODE"
	case DocumentNode:
		return "DOCUMENT_NODE"
	case DocumentTypeNode:
		return "DOCUMENT_TYPE_NODE"
	case DocumentFragmentNode:
		return "DOCUMENT_FRAGMENT_NODE"
	default:
		return "???"
	}
}

var (
	bleh = ":abcdefghijklmnopqrstuvwxyz_"
)

// IsValidName checks whether the given string s is a valid XML name for use
// in Elements and Attribute names.
func IsValidName(s string) bool {
	if s == "" {
		return false
	}
	if strings.ContainsAny(s[0:1], bleh) ||
		strings.ContainsAny(strings.ToUpper(s), bleh) ||
		(s[0] >= 0xC0 && s[0] <= 0xD6) ||
		(s[0] >= 0xD8 && s[0] <= 0xF6) {

		return true
	}

	return false
}

// NamedNodeMap represent collections of nodes that can be accessed by name.
type NamedNodeMap interface {
	Node

	GetNamedItem(string) Node
	SetNamedItem(Node)
	Length() int
}

// Node is the primary interface for the entire Document Object Model. It represents
// a single node in the document tree. While all objects implementing the Node
// interface expose methods for dealing with children, not all objects implementing
// the Node interface may have children.
type Node interface {
	NodeName() string
	NodeType() NodeType
	NodeValue() string
	LocalName() string
	// Gets the list of child nodes.
	NodeList() []Node
	// Gets the parent node. May be nil if none was assigned.
	ParentNode() Node
	// Gets the first child Node of this Node. May return nil if no child nodes
	// exist.
	FirstChild() Node
	GetAttributes() NamedNodeMap
	// Gets the owner document (the Document instance which was used to create
	// the Node).
	OwnerDocument() Document
	// Appends a child to this Node. Will return an error when this Node is not
	// able to have any (more) children, like Text nodes.
	AppendChild(Node) error
	// Returns true when the Node has one or more children.
	HasChildNodes() bool
	// Returns the namespace URI of this node.
	NamespaceURI() string

	setNodeType(NodeType)
	setParentNode(Node)
	setOwnerDocument(Document)
	setNamespaceURI(string)
}

// Attr represents an attribute in an Element.
type Attr interface {
	Node

	GetName() string
	IsSpecified() bool
	GetValue() string
	SetValue(string)
	GetOwnerElement() Element
}

// Element represents an element in an HTML or XML document.
type Element interface {
	Node

	// Sets the tag name of this element.
	SetTagName(tagname string)
	// Gets the tag name of this element.
	GetTagName() string

	SetAttribute(name, value string)
	GetAttribute(name string) string
}

// Text is a Node that represents character data.
type Text interface {
	Node

	GetData() string
	SetData(s string)
}

// DocumentType belongs to a Document, but can also be nil. The DocumentType
// interface in the DOM Core provides an interface to the list of entities
// that are defined for the document, and little else because the effect of
// namespaces and the various XML schema efforts on DTD representation are
// not clearly understood as of this writing. (Direct copy of the spec).
type DocumentType interface {
	Node

	// Gets the name of the DTD; i.e.  the name immediately following the DOCTYPE keyword.
	GetName() string
	// The public identifier of the external subset.
	GetPublicID() string
	// The system identifier of hte external subset. This may be an absolute URI or not.
	GetSystemID() string
}

// Document is the root of the Document Object Model.
type Document interface {
	Node

	// Creates an element with the given tagname and returns it. Will return
	// an ErrorInvalidCharacter if the specified name is not an XML name according
	// to the XML version in use, specified in the XMLVersion attribute.
	CreateElement(tagName string) (Element, error)
	// Creates an element of the givens qualified name and namespace URI, and
	// returns it. Use an empty string if no namespace is necessary. See
	// CreateElement(string).
	CreateElementNS(namespaceURI, tagName string) (Element, error)
	// Creates a Text node given the specified string and returns it.
	CreateTextNode(string) Text
	// Creates an Attr of the given name and returns it.
	CreateAttribute(name string) (Attr, error)

	// Gets the document element, which should be the first (and only) child Node
	// of the Document. Can be nil if none is set yet.
	GetDocumentElement() Element
}

//================================================================================

type domNamedNodeMap struct {
	Node

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

func (dn *domNamedNodeMap) SetNamedItem(node Node) {
	// assert that node must be an Attr
	dn.Attrs[node.NodeName()] = node
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

func (dn *domNode) NodeValue() string {
	return "TODO: what?"
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
	dn.nodes = append(dn.nodes, node)
	node.setParentNode(dn)
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

func (de *domDocument) CreateElement(tagName string) (Element, error) {
	elem := newElement()
	elem.SetTagName(tagName)
	elem.setOwnerDocument(de)
	return elem, nil
}

func (de *domDocument) CreateElementNS(namespaceURI, tagName string) (Element, error) {
	elem, err := de.CreateElement(tagName)
	if err != nil {
		return nil, err
	}
	elem.setNamespaceURI(namespaceURI)
	return elem, nil
}

func (de *domDocument) CreateTextNode(t string) Text {
	text := newText()
	text.setOwnerDocument(de)
	text.SetData(t)
	return text
}

func (dd *domDocument) CreateAttribute(name string) (Attr, error) {
	attr := newAttr(name)
	return attr, nil
}

// 'Override' the AppendChild() function from the Node interface. One child can
// be appended when the node list is empty. The first child of the document will
// be the document element. Subsequent calls will result in an error.
func (de *domDocument) AppendChild(n Node) error {
	if len(de.NodeList()) <= 0 {
		de.Node.AppendChild(n)
		return nil
	}
	return ErrorHierarchyRequest
}

func (de *domDocument) GetDocumentElement() Element {
	if len(de.NodeList()) == 1 {
		node := de.FirstChild()
		if node.NodeType() == ElementNode {
			elem := node.(Element)
			return elem
		}
	}
	return nil
}

func (de *domDocument) String() string {
	return fmt.Sprintf("%s", de.NodeType())
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

// NodeName is an override from Node.
func (da *domAttr) NodeName() string {
	return da.Name
}

func (da *domAttr) GetName() string {
	return da.Name
}

// NodeValue is an override from Node.
func (da *domAttr) NodeValue() string {
	return da.Value
}

func (da *domAttr) GetValue() string {
	return da.Value
}
func (da *domAttr) SetValue(val string) {
	da.Value = val
}

func (da *domAttr) IsSpecified() bool {
	return da.Specified
}

func (da *domAttr) GetOwnerElement() Element {
	return da.OwnerElement
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

func (dt *domText) SetData(s string) {
	dt.Data = s
}

func (dt *domText) GetData() string {
	return dt.Data
}

func (dt *domText) String() string {
	return fmt.Sprintf("%s, %s, %s", dt.NodeType(), dt.NodeName(), dt.GetData())
}
