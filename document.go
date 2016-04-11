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
}

// NewDocument creates a new Document which can be used to create
// custom documents using the methods supplied.
func NewDocument() Document {
	return &domDocument{}
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
	if dd.HasChildNodes() {
		return dd.nodes[0]
	}
	return nil
}

func (dd *domDocument) GetLastChild() Node {
	if dd.HasChildNodes() {
		return dd.nodes[len(dd.nodes)-1]
	}
	return nil
}

func (dd *domDocument) GetAttributes() NamedNodeMap {
	return nil
}

func (dd *domDocument) HasAttributes() bool {
	return false
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

	if child == dd {
		return fmt.Errorf("%v: adding a node to itself as a child", ErrorHierarchyRequest)
	}

	if child.GetOwnerDocument() != dd {
		return ErrorWrongDocument
	}

	if child.GetNodeType() == ElementNode {
		// Check if a Document element is already appended.
		if dd.GetDocumentElement() != nil {
			return fmt.Errorf("%v: a Document element already exists (<%v>)", ErrorHierarchyRequest, dd.GetDocumentElement())
		}
	}

	if child.GetNodeType() == AttributeNode || child.GetNodeType() == TextNode {
		return ErrorHierarchyRequest
	}

	// Child already has a parent. Remove it!
	cparent := child.GetParentNode()
	if cparent != nil {
		cparent.RemoveChild(child)
	}

	child.setParentNode(dd)
	dd.nodes = append(dd.nodes, child)
	return nil
}

func (dd *domDocument) RemoveChild(oldChild Node) (Node, error) {
	if oldChild == nil {
		return nil, nil
	}

	for i, child := range dd.GetChildNodes() {
		if child == oldChild {
			// Slice trickery to remove the node at the found index:
			dd.nodes = append(dd.nodes[:i], dd.nodes[i+1:]...)
			return child, nil
		}
	}

	return nil, ErrorNotFound
}

func (dd *domDocument) ReplaceChild(newChild, oldChild Node) (Node, error) {
	if newChild == nil {
		return nil, fmt.Errorf("%v: given new child is nil", ErrorHierarchyRequest)
	}
	if oldChild == nil {
		return nil, fmt.Errorf("%v: given old child is nil", ErrorHierarchyRequest)
	}
	if newChild.GetNodeType() == AttributeNode || newChild.GetNodeType() == TextNode {
		return nil, ErrorHierarchyRequest
	}

	// newChild must be created by the same owner document of this element.
	if newChild.GetOwnerDocument() != dd {
		return nil, ErrorWrongDocument
	}

	// Replacing a Node (which is not an element) with an element when there's already an element, should fail.
	if dd.GetDocumentElement() != nil && newChild.GetNodeType() == ElementNode && oldChild.GetNodeType() != ElementNode {
		return nil, ErrorHierarchyRequest
	}

	// Find the old child, and replace it with the new child.
	for i, child := range dd.GetChildNodes() {
		if child == oldChild {
			// Check if newChild has a parent (i.e., it's in the tree).
			ncParent := newChild.GetParentNode()
			if ncParent != nil {
				// Remove the newChild from its parent.
				ncParent.RemoveChild(newChild)
			}

			// Slice trickery, again. It will make a new underlying slice with one element,
			// the 'newChild', and then append the rest of the de.nodes to that.
			dd.nodes = append(dd.nodes[:i], append([]Node{newChild}, dd.nodes[i+1:]...)...)
			// Change the parent node:
			newChild.setParentNode(dd)

			return oldChild, nil
		}
	}

	return nil, ErrorNotFound
}

func (dd *domDocument) InsertBefore(newChild, refChild Node) (Node, error) {
	// If a document element is already specified, and another element is attempted
	// to insert, return an error.
	if newChild == nil {
		// FIXME: what in this case? Is an error ok?
		return nil, ErrorHierarchyRequest
	}

	// If refChild is nil, append to the end, and return.
	if refChild == nil {
		err := dd.AppendChild(newChild)
		if err != nil {
			return nil, err
		}
		return newChild, nil
	}

	// Cannot insert an element if there's already one element.
	if newChild.GetNodeType() == ElementNode && dd.GetDocumentElement() != nil {
		return nil, fmt.Errorf("%v: a Document element already exists (<%v>)", ErrorHierarchyRequest, dd.GetDocumentElement())
	}

	if newChild.GetNodeType() == AttributeNode || newChild.GetNodeType() == TextNode {
		return nil, ErrorHierarchyRequest
	}

	if newChild.GetOwnerDocument() != dd {
		return nil, ErrorWrongDocument
	}

	// Find the reference child, insert newChild before that one.
	for i, child := range dd.GetChildNodes() {
		if child == refChild {
			// Check if newChild is in the tree already. If so, remove it.
			ncParent := newChild.GetParentNode()
			if ncParent != nil {
				ncParent.RemoveChild(newChild)
			}
			newChild.setParentNode(dd)
			dd.nodes = append(dd.nodes[:i], append([]Node{newChild}, dd.nodes[i:]...)...)
			return newChild, nil
		}
	}

	return nil, ErrorNotFound
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

func (dd *domDocument) setOwnerDocument(doc Document) {
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

func (dd *domDocument) NormalizeDocument() {
	counter := 0
	for _, c := range dd.GetChildNodes() {
		if e, ok := c.(Element); ok {
			e.normalizeNamespaces(&counter)
		}
	}
}

func (dd *domDocument) LookupPrefix(namespace string) string {
	return ""
}

func (dd *domDocument) LookupNamespaceURI(pfx string) (string, bool) {
	if dd.GetDocumentElement() != nil {
		return dd.GetDocumentElement().LookupNamespaceURI(pfx)
	}
	return "", false
}

// GetTextContent should return null, but Go doesn't allow null strings so this method
// will return an empty string.
func (dd *domDocument) GetTextContent() string {
	return ""
}

// SetTextContent does nothing on a Document Node.
func (dd *domDocument) SetTextContent(content string) {
	// no-op.
}

// CloneNode creates a copy of the Document instance. When deep is true, it will create a complete copy
// of the whole Document, recursively. When false, it's pretty useless since it will return just a plain new
// empty Document.
func (dd *domDocument) CloneNode(deep bool) Node {
	cloneDoc := NewDocument()

	if deep {
		for _, c := range dd.GetChildNodes() {
			cloneChild := cloneDoc.ImportNode(c, true)
			err := cloneDoc.AppendChild(cloneChild)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return cloneDoc
}

func (dd *domDocument) ImportNode(n Node, deep bool) Node {
	return importNode(dd, n, deep)
}

func (dd *domDocument) String() string {
	return fmt.Sprintf("%s", dd.GetNodeType())
}
