package dom

import (
	"fmt"
)

type domElement struct {
	localName     string
	nodes         []Node
	parentNode    Node
	firstChild    Node
	attributes    NamedNodeMap
	ownerDocument Document
	namespaceURI  string

	// Element specific things:
	tagName string
}

func newElement() Element {
	e := &domElement{}
	return e
}

func (de *domElement) GetNodeName() string {
	return de.tagName
}

func (de *domElement) GetNodeType() NodeType {
	return ElementNode
}

// NodeValue should return null/nil for Element types like the spec says,
// but Go does not permit nil strings which are not pointers. So for now we
// just return an empty string at all times.
func (de *domElement) GetNodeValue() string {
	return ""
}

func (de *domElement) GetLocalName() string {
	// TODO: what?
	return de.tagName
}

func (de *domElement) GetChildNodes() []Node {
	return de.nodes
}

func (de *domElement) GetParentNode() Node {
	return de.parentNode
}

func (de *domElement) GetFirstChild() Node {
	return de.nodes[0]
}

func (de *domElement) GetAttributes() NamedNodeMap {
	return de.attributes
}

func (de *domElement) GetOwnerDocument() Document {
	return de.ownerDocument
}

func (de *domElement) AppendChild(child Node) error {
	if de == child {
		return fmt.Errorf("%v: adding a node to itself as a child", ErrorHierarchyRequest)
	}
	child.setParentNode(de)
	de.nodes = append(de.nodes, child)
	return nil
}

func (de *domElement) HasChildNodes() bool {
	return len(de.nodes) > 0
}

func (de *domElement) GetNamespaceURI() string {
	return de.namespaceURI
}

func (de *domElement) SetTagName(name string) {
	de.tagName = name
}

func (de *domElement) GetTagName() string {
	return de.tagName
}

func (de *domElement) SetAttribute(name, value string) {
	if de.attributes == nil {
		de.attributes = newNamedNodeMap()
	}

	attr := newAttr()
	attr.SetName(name)
	attr.SetValue(value)
	attr.setOwnerElement(de)
	de.attributes.SetNamedItem(attr)
}

func (de *domElement) GetAttribute(name string) string {
	return ""
}

// Private functions:
func (de *domElement) setParentNode(parent Node) {
	de.parentNode = parent
}

func (de *domElement) setOwnerDocument(d Document) {
	de.ownerDocument = d
}

func (de *domElement) setNamespaceURI(uri string) {
	de.namespaceURI = uri
}

func (de *domElement) String() string {
	return fmt.Sprintf("%s, <%s>, ns=%s, attrs=%v",
		de.GetNodeType(), de.tagName, de.namespaceURI, de.attributes)
}
