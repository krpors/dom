package dom

import (
	"fmt"
)

type domDocument struct {
	localName    string
	nodes        []Node
	firstChild   Node
	attributes   NamedNodeMap
	namespaceURI string

	// Element specific things:
	tagName string
}

func NewDocument() Document {
	d := &domDocument{}
	return d
}

// NODE SPECIFIC FUNCTIONS

func (dd *domDocument) NodeName() string {
	return "#document"
}

func (dd *domDocument) NodeType() NodeType {
	return DocumentNode
}

// NodeValue should return null/nil for Document types like the spec says,
// but Go does not permit nil strings which are not pointers. So for now we
// just return an empty string at all times.
func (dd *domDocument) NodeValue() string {
	return ""
}

func (dd *domDocument) LocalName() string {
	// TODO: what?
	return dd.tagName
}

func (dd *domDocument) NodeList() []Node {
	return dd.nodes
}

func (dd *domDocument) ParentNode() Node {
	return nil
}

func (dd *domDocument) FirstChild() Node {
	return dd.nodes[0]
}

func (dd *domDocument) GetAttributes() NamedNodeMap {
	return nil
}

func (dd *domDocument) OwnerDocument() Document {
	return nil
}

func (dd *domDocument) AppendChild(child Node) error {
	if child == nil {
		return nil
	}

	if dd == child {
		return fmt.Errorf("%v: adding a node to itself as a child", ErrorHierarchyRequest)
	}

	// Only allow elements to be append as a child... for now!
	switch typ := child.(type) {
	case Element:
		if len(dd.NodeList()) <= 0 {
			child.setParentNode(dd)
			dd.nodes = append(dd.nodes, child)
			return nil
		}
	default:
		return fmt.Errorf("only nodes of type %v can be added (tried '%v')", ElementNode, typ.NodeType())
	}

	return fmt.Errorf("%v: document can only have one child, which must be of type Element", ErrorHierarchyRequest)
}

func (dd *domDocument) HasChildNodes() bool {
	return len(dd.nodes) > 0
}

// NamespaceURI should return nil as per the spec, but Go doesn't allow that for
// non-pointer types, so return an empty string instead.
func (dd *domDocument) NamespaceURI() string {
	return ""
}

// Private functions:
func (dd *domDocument) setParentNode(parent Node) {
	// no-op
}

func (dd *domDocument) setOwnerDocument(d Document) {
	// no-op
}

func (dd *domDocument) setNamespaceURI(uri string) {
	// no-op
}

// DOCUMENT SPECIFIC FUNCTIONS
func (dd *domDocument) CreateElement(tagName string) (Element, error) {
	e := newElement()
	e.setOwnerDocument(dd)
	e.SetTagName(tagName)
	return e, nil
}

func (dd *domDocument) CreateElementNS(namespaceURI, tagName string) (Element, error) {
	e, err := dd.CreateElement(tagName)
	if err != nil {
		return nil, err
	}
	e.setNamespaceURI(namespaceURI)
	return e, nil
}

func (dd *domDocument) CreateTextNode(text string) Text {
	t := newText()
	t.setOwnerDocument(dd)
	t.SetData(text)
	return t
}

func (dd *domDocument) CreateAttribute(name string) (Attr, error) {
	return nil, nil
}

func (dd *domDocument) GetDocumentElement() Element {
	firstNode := dd.NodeList()[0]
	bleh := firstNode.(Element)
	return bleh
}

func (dd *domDocument) String() string {
	return fmt.Sprintf("%s", dd.NodeType())
}

// ===

func ToXML(node Node) string {
	var xml string

	var xtree func(n Node, padding string)

	xtree = func(n Node, padding string) {
		switch t := n.(type) {
		case Element:
			xml += fmt.Sprintf("<%s>", t.GetTagName())
		//case Comment: // TODO!!
		case Text:
			xml += t.GetData()
		}

		for _, node := range n.NodeList() {
			xtree(node, padding+"  ")
		}

		switch t := n.(type) {
		case Element:
			xml += fmt.Sprintf("</%s>", t.GetTagName())
		}
	}

	xtree(node, "")

	return xml

}
