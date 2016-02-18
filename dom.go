package dom

const (
	ElementNode = iota
	AttributeNode
	TextNode
	CDataSectionNode
	EntityReferenceNode
	EntityNode
	ProcessingInstructionNode
	CommentNode
	DocumentNode
	DocumentTypeNode
	DocumentFragmentNode
)

type Node interface {
	NodeName() string
	NodeValue() string
	LocalName() string
	NodeList() []Node
	ParentNode() Node
	FirstChild() Node
	OwnerDocument() Document
	AppendChild(node Node)
	HasChildNodes() bool

	// 'Setters':
	setParentNode(node Node)
}

type Document interface {
	Node
}

type domNode struct {
	localName     string
	nodeName      string
	nodeValue     string
	nodes         []Node
	parentNode    Node
	firstChild    Node
	ownerDocument Document
}

func NewNode(name string) Node {
	dn := &domNode{}
	dn.nodeName = name
	return dn
}

func (dn *domNode) NodeName() string {
	return dn.nodeName
}

func (dn *domNode) NodeValue() string {
	return dn.nodeValue
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

func (dn *domNode) setParentNode(node Node) {
	dn.parentNode = node
}
