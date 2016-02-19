package dom

import (
	"errors"
)

var (
	ErrorHierarchyRequest = errors.New("HIERARCHY_REQUEST_ERR: an attempt was made to insert a node where it is not permitted")
)

type NodeType uint8

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

type NamedNodeMap interface {
	GetNamedItem(string) Node
	SetNamedItem(Node)
	Length() int
}

type Node interface {
	NodeName() string
	NodeType() NodeType
	NodeValue() string
	LocalName() string
	NodeList() []Node
	ParentNode() Node
	FirstChild() Node
	Attributes() NamedNodeMap
	OwnerDocument() Document
	AppendChild(Node) error
	HasChildNodes() bool
	NamespaceUri() string

	// 'Setters':
	setNodeType(NodeType)
	setParentNode(Node)
}

type Element interface {
	Node

	SetTagName(tagname string)
	GetTagName() string
}

type Text interface {
	Node

	GetData() string
	SetData(s string)
}

type Document interface {
	Node
}

//================================================================================

type domNamedNodeMap struct {
	attributes map[string]Node
}

func NewNamedNodeMap() NamedNodeMap {
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
	namespaceUri  string
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

func (dn *domNode) NamespaceUri() string {
	return dn.namespaceUri
}

func (dn *domNode) setParentNode(node Node) {
	dn.parentNode = node
}

func (dn *domNode) setNodeType(t NodeType) {
	dn.nodeType = t
}

//================================================================================

type domDocument struct {
	Node
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

func NewDocument() Document {
	d := &domDocument{}
	d.Node = &domNode{}
	d.setNodeType(DocumentNode)
	return d
}

//================================================================================

type domElement struct {
	Node

	TagName string
}

func (de *domElement) SetTagName(name string) {
	de.TagName = name
}

func (de *domElement) GetTagName() string {
	return de.TagName
}

func NewElement() Element {
	e := &domElement{}
	e.Node = &domNode{}
	e.setNodeType(ElementNode)
	return e
}

//================================================================================

type domText struct {
	Node

	Data string
}

func NewText() Text {
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
