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

// GetFirstChild returns the first child in the document. Possible nodes are
// Comment, ProcessingInstruction, or Element.
func (dd *domDocument) GetFirstChild() Node {
	return dd.nodes[0]
}

func (dd *domDocument) GetAttributes() NamedNodeMap {
	return nil
}

func (dd *domDocument) GetOwnerDocument() Document {
	return nil
}

// AppendChild handles the appending of nodes to the document, and fails accordingly.
// Only 1 Element may be appended, but comments and processing instructions may appear
// in abundance.
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
		// Check if a Document element is already appended.
		docelem := dd.GetDocumentElement()
		if docelem == nil {
			child.setParentNode(dd)
			dd.nodes = append(dd.nodes, child)
			return nil
		}
		return fmt.Errorf("%v: a Document element already exists (<%v>)", ErrorHierarchyRequest, docelem)
	case ProcessingInstruction:
		// Processing instructions are legal children of a DOM Document and can appear
		// anywhere, even before the Document element.
		child.setParentNode(dd)
		dd.nodes = append(dd.nodes, child)
		return nil
	case Comment:
		child.setParentNode(dd)
		dd.nodes = append(dd.nodes, child)
		return nil
	default:
		return fmt.Errorf("only nodes of type (%v | %v | %v) can be added to a Document (tried '%v')",
			ElementNode, ProcessingInstructionNode, CommentNode, typ.GetNodeType())
	}
}

func (dd *domDocument) HasChildNodes() bool {
	return len(dd.nodes) > 0
}

// NamespaceURI should return nil as per the spec, but Go doesn't allow that for
// non-pointer types, so return an empty string instead.
func (dd *domDocument) GetNamespaceURI() string {
	return ""
}

// GetNamespacePrefix returns... ?
func (dd *domDocument) GetNamespacePrefix() string {
	// TODO: namespace prefix
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

// CreateelementNS creates an element with the given namespace URI and tagname. If I recall correctly,
// the DOM spec mentions something about not caring about namespace URIs. As long as they are escaped,
// it's okay. Even the Xerces implementation in Java doesn't care about the namespace URI, and will be
// serialized just fine.
func (dd *domDocument) CreateElementNS(namespaceURI, tagName string) (Element, error) {
	e, err := dd.CreateElement(tagName)
	if err != nil {
		return nil, err
	}
	e.setNamespaceURI(namespaceURI)
	return e, nil
}

func (dd *domDocument) CreateText(text string) Text {
	t := newText()
	t.setOwnerDocument(dd)
	t.SetData(text)
	return t
}

// CreateComment creates a comment node and returns it. When the comment string contains
// a double-hyphen (--) it will return an error and the Comment will be nil. The spec
// says something differently though:
//
// No lexical check is done on the content of a comment and it is therefore possible to
// have the character sequence "--" (double-hyphen) in the content, which is illegal in
// a comment per section 2.5 of [XML 1.0]. The presence of this character sequence must
// generate a fatal error **during serialization**.
//
// E.g. this implementation doesn't fail during serialization, but way before. This may
// be subject to change to get conform the spec. The Xerces implementation in Java 8
// doesn't fail serialization, for example , but simply replaces the '--' with '- -'.
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
	xmlname := XMLName(name)
	if !xmlname.IsValid() {
		return nil, fmt.Errorf("%v: '%v'", ErrorInvalidCharacter, xmlname)
	}

	attr := newAttr()
	attr.setName(name)
	attr.setOwnerDocument(dd)
	return attr, nil
}

func (dd *domDocument) CreateProcessingInstruction(target, data string) (ProcessingInstruction, error) {
	pi := newProcInst()
	pi.setOwnerDocument(dd)
	pi.setParentNode(dd)
	pi.setData(data)
	pi.setTarget(target)
	return pi, nil
}

// GetDocumentElement traverses through the child nodes and finds the first Element.
// That one will be returned as the Document element. The AppendChild function must
// take care that no two root nodes can be added to this Document.
func (dd *domDocument) GetDocumentElement() Element {
	// No nodes, so return nil.
	if len(dd.nodes) <= 0 {
		return nil
	}

	for _, node := range dd.nodes {
		if e, ok := node.(Element); ok {
			return e
		}
	}

	return nil
}

func (dd *domDocument) String() string {
	return fmt.Sprintf("%s", dd.GetNodeType())
}
