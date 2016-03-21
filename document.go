package dom

import (
	"fmt"
	"strings"
)

// TODO: implement name list, which is just a map of namespaces + prefixes
// known to the Document. The key should be a prefix, and the namespace URI
// is a value for the given prefix. The namespace is actually unique, but
// that means one namespace can only have one prefix, but that's not something
// XML forces. E.g. xmlns:pfx="urn" and xmlns:bleh="urn" is the same.
//
// After creating a node, find its prefix + namespace, and add it to this map.
// OR: just add namespaces in a map, use some clever thing to create a prefix
// from it? Meh.

type domDocument struct {
	nodes []Node

	nameList map[string]string // map of namespace + prefix
}

// NewDocument creates a new Document which can be used to create
// custom documents using the methods supplied.
func NewDocument() Document {
	d := &domDocument{}
	d.nameList = make(map[string]string)
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
	return ""
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

func (dd *domDocument) RemoveChild(oldChild Node) (Node, error) {
	panic("not implemented yet")
}
func (dd *domDocument) ReplaceChild(oldChild Node) (Node, error) {
	panic("not implemented yet")
}
func (dd *domDocument) InsertBefore(newChild, refChild Node) (Node, error) {
	panic("not implemented yet")
}

func (dd *domDocument) HasChildNodes() bool {
	return len(dd.nodes) > 0
}

func (dd *domDocument) GetPreviousSibling() Node {
	return nil
}

func (dd *domDocument) GetNextSibling() Node {
	return nil
}

// NamespaceURI should return nil as per the spec, but Go doesn't allow that for
// non-pointer types, so return an empty string instead.
func (dd *domDocument) GetNamespaceURI() string {
	return ""
}

func (dd *domDocument) GetNamespacePrefix() string {
	return ""
}

// Private functions:
func (dd *domDocument) setParentNode(parent Node) {
	// no-op
}

// DOCUMENT SPECIFIC FUNCTIONS
func (dd *domDocument) CreateElement(tagName string) (Element, error) {
	name := XMLName(tagName)
	if !name.IsValid() {
		return nil, fmt.Errorf("%v; tagname '%v'", ErrorInvalidCharacter, tagName)
	}

	e := newElement(dd, tagName, "")
	return e, nil
}

// CreateelementNS creates an element with the given namespace URI and tagname. If I recall correctly,
// the DOM spec mentions something about not caring about namespace URIs. As long as they are escaped,
// it's okay. Even the Xerces implementation in Java doesn't care about the namespace URI, and will be
// serialized just fine.
func (dd *domDocument) CreateElementNS(namespaceURI, tagName string) (Element, error) {
	name := XMLName(tagName)
	if !name.IsValid() {
		return nil, fmt.Errorf("%v; tagname '%v'", ErrorInvalidCharacter, tagName)
	}

	e := newElement(dd, tagName, namespaceURI)
	return e, nil
}

func (dd *domDocument) CreateText(text string) Text {
	t := newText(dd)
	t.SetText(text)
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

	c := newComment(dd)
	c.SetComment(comment)
	return c, nil
}

func (dd *domDocument) CreateAttribute(name string) (Attr, error) {
	xmlname := XMLName(name)
	if !xmlname.IsValid() {
		return nil, fmt.Errorf("%v: '%v'", ErrorInvalidCharacter, xmlname)
	}

	attr := newAttr(dd, name, "")
	return attr, nil
}

func (dd *domDocument) CreateAttributeNS(namespaceURI, name string) (Attr, error) {
	xmlname := XMLName(name)
	if !xmlname.IsValid() {
		return nil, fmt.Errorf("%v: '%v'", ErrorInvalidCharacter, xmlname)
	}

	attr := newAttr(dd, name, namespaceURI)
	return attr, nil
}

func (dd *domDocument) CreateProcessingInstruction(target, data string) (ProcessingInstruction, error) {
	pi := newProcInst(dd, target, data)
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

// GetElementsByTagName finds all descendant elements of the current element,
// with the given tag name, in document order.
func (dd *domDocument) GetElementsByTagName(tagname string) []Element {
	return getElementsBy(dd, "", tagname, false)
}

// GetElementsByTagNameNS finds all descendant elements of the current element,
// with the given tag name and namespace URI, in document order.
func (dd *domDocument) GetElementsByTagNameNS(namespaceURI, tagname string) []Element {
	return getElementsBy(dd, namespaceURI, tagname, true)
}

func (dd *domDocument) LookupPrefix(namespace string) string {
	return ""
}

func (dd *domDocument) LookupNamespaceURI(pfx string) string {
	if dd.GetDocumentElement() != nil {
		return dd.GetDocumentElement().LookupNamespaceURI(pfx)
	}
	return ""
}

func (dd *domDocument) String() string {
	return fmt.Sprintf("%s", dd.GetNodeType())
}
