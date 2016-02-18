package dom

type NodeType uint8

const (
	UndefinedNode NodeType = iota
	ElementNode
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
	AppendChild(Node)
	HasChildNodes() bool
	NamespaceUri() string

	// 'Setters':
	setNodeType(NodeType)
	setParentNode(Node)
}

type Element interface {
	Node
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

func NewNode(name string) Node {
	dn := &domNode{}
	dn.nodeType = UndefinedNode
	dn.nodeName = name
	return dn
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
	return dn.firstChild
}

func (dn *domNode) Attributes() NamedNodeMap {
	return dn.attributes
}

func (dn *domNode) OwnerDocument() Document {
	return dn.ownerDocument
}

func (dn *domNode) AppendChild(node Node) {
	dn.nodes = append(dn.nodes, node)
	node.setParentNode(dn)
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

type domElement struct {
	Node
}

func NewElement() Element {
	n := &domNode{}
	de := &domElement{n}
	de.setNodeType(ElementNode)
	de.HasChildNodes()
	return de

}
