package dom

import (
	"fmt"
	"strings"
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

// NewDocument creates a new Document which can be used to create
// custom documents using the methods supplied.
func NewDocument() Document {
	d := &domDocument{}
	return d
}

// NODE SPECIFIC FUNCTIONS

func (dd *domDocument) GetNodeName() string {
	return "#document"
}

func (dd *domDocument) GetNodeType() NodeType {
	return DocumentNode
}

// NodeValue should return null/nil for Document types like the spec says,
// but Go does not permit nil strings which are not pointers. So for now we
// just return an empty string at all times.
func (dd *domDocument) GetNodeValue() string {
	return ""
}

func (dd *domDocument) GetLocalName() string {
	// TODO: what?
	return dd.tagName
}

func (dd *domDocument) GetChildNodes() []Node {
	return dd.nodes
}

func (dd *domDocument) GetParentNode() Node {
	return nil
}

func (dd *domDocument) GetFirstChild() Node {
	return dd.nodes[0]
}

func (dd *domDocument) GetAttributes() NamedNodeMap {
	return nil
}

func (dd *domDocument) GetOwnerDocument() Document {
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
		if len(dd.GetChildNodes()) <= 0 {
			child.setParentNode(dd)
			dd.nodes = append(dd.nodes, child)
			return nil
		}
	case ProcessingInstruction:
		// Processing instructions are legal children of a DOM Document and can appear
		// anywhere, even before the Document element.
		// TODO processing instructions in their own node list?
		// child.setParentNode(dd)
		// dd.nodes = append(dd.nodes, child)
		// fmt.Println(dd.nodes)
		// fmt.Println(child.GetNodeName(), child.GetNodeValue())
		return nil
	default:
		return fmt.Errorf("only nodes of type (%v | %v) can be added (tried '%v')",
			ElementNode, ProcessingInstructionNode, typ.GetNodeType())
	}

	return fmt.Errorf("%v: document can only have one child, which must be of type Element", ErrorHierarchyRequest)
}

func (dd *domDocument) HasChildNodes() bool {
	return len(dd.nodes) > 0
}

// NamespaceURI should return nil as per the spec, but Go doesn't allow that for
// non-pointer types, so return an empty string instead.
func (dd *domDocument) GetNamespaceURI() string {
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
	name := XMLName(tagName)
	if !name.IsValid() {
		return nil, fmt.Errorf("%v; tagname '%v'", ErrorInvalidCharacter, tagName)
	}

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

func (dd *domDocument) CreateComment(comment string) (Comment, error) {
	if strings.ContainsAny(comment, "--") {
		return nil, fmt.Errorf("%v: comments may not contain a double hyphen (--)", ErrorInvalidCharacter)
	}

	c := newComment()
	c.setOwnerDocument(dd)
	c.SetComment(comment)
	return c, nil
}

func (dd *domDocument) CreateAttribute(name string) (Attr, error) {
	return nil, nil
}

func (dd *domDocument) CreateProcessingInstruction(target, data string) (ProcessingInstruction, error) {
	pi := &domProcInst{}
	pi.setOwnerDocument(dd)
	pi.setParentNode(dd)
	pi.data = data
	pi.target = target
	return pi, nil
}

func (dd *domDocument) GetDocumentElement() Element {
	firstNode := dd.GetChildNodes()[0]
	bleh := firstNode.(Element)
	return bleh
}

func (dd *domDocument) String() string {
	return fmt.Sprintf("%s", dd.GetNodeType())
}
