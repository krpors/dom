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

// NamedNodeMap represent collections of nodes that can be accessed by name.
type NamedNodeMap interface {
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
	Attributes() NamedNodeMap
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

// Element represents an element in an HTML or XML document.
type Element interface {
	Node

	// Sets the tag name of this element.
	SetTagName(tagname string)
	// Gets the tag name of this element.
	GetTagName() string
}

// Text is a Node that represents character data.
type Text interface {
	Node

	GetData() string
	SetData(s string)
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
	// Gets the document element, which should be the first (and only) child Node
	// of the Document. Can be nil if none is set yet.
	GetDocumentElement() Element
}

//================================================================================

type domNamedNodeMap struct {
	attributes map[string]Node
}

func newNamedNodeMap() NamedNodeMap {
	nnm := &domNamedNodeMap{}
	nnm.attributes = make(map[string]Node)
	return nnm
}

func (dn *domNamedNodeMap) GetNamedItem(s string) Node {
	return dn.attributes[s]
}

func (dn *domNamedNodeMap) SetNamedItem(node Node) {
	dn.attributes[node.NodeName()] = node
}

func (dn *domNamedNodeMap) Length() int {
	return len(dn.attributes)
}

//================================================================================

type domNode struct {
	localName     string
	nodeName      string
	nodeValue     string
	nodeType      NodeType
	nodes         []Node
	parentNode    Node
	firstChild    Node
	attributes    NamedNodeMap
	ownerDocument Document
	namespaceURI  string
}

func (dn *domNode) NodeName() string {
	return dn.nodeName
}

func (dn *domNode) NodeValue() string {
	return dn.nodeValue
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

func (dn *domNode) Attributes() NamedNodeMap {
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
	e.Node = &domNode{}
	e.setNodeType(ElementNode)
	return e
}

func (de *domElement) SetTagName(name string) {
	de.TagName = name
}

func (de *domElement) GetTagName() string {
	return de.TagName
}

func (de *domElement) String() string {
	return fmt.Sprintf("%s, <%s>", de.NodeType(), de.GetTagName())
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

func (dt *domText) SetData(s string) {
	dt.Data = s
}

func (dt *domText) GetData() string {
	return dt.Data
}

func (dt *domText) String() string {
	return fmt.Sprintf("%s, %s", dt.NodeType(), dt.GetData())
}
