package dom

import (
	"fmt"
	"strings"
)

type domElement struct {
	localName     string       // The local part, without the (optional) namespace prefix.
	nodes         []Node       // Child nodes.
	parentNode    Node         // Parent node
	attributes    NamedNodeMap // Attributes on this element.
	ownerDocument Document     // Owner document.
	namespaceURI  string       // Namespace uri.

	// Element specific things:
	tagName XMLName // The complete tagname given, with prefix.
}

func newElement(owner Document, tagname string, namespaceURI string) Element {
	e := &domElement{}
	e.ownerDocument = owner
	e.tagName = XMLName(tagname)
	e.namespaceURI = namespaceURI
	e.attributes = newNamedNodeMap()
	return e
}

func (de *domElement) GetNodeName() string {
	return string(de.tagName)
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
	return de.tagName.GetLocalPart()
}

func (de *domElement) GetChildNodes() []Node {
	return de.nodes
}

func (de *domElement) GetParentNode() Node {
	return de.parentNode
}

func (de *domElement) GetFirstChild() Node {
	if de.HasChildNodes() {
		return de.nodes[0]
	}
	return nil
}

func (de *domElement) GetLastChild() Node {
	if de.HasChildNodes() {
		return de.nodes[len(de.nodes)-1]
	}
	return nil
}

func (de *domElement) GetAttributes() NamedNodeMap {
	return de.attributes
}

func (de *domElement) HasAttributes() bool {
	return len(de.attributes.GetItems()) > 0
}

func (de *domElement) GetOwnerDocument() Document {
	return de.ownerDocument
}

func (de *domElement) AppendChild(child Node) error {
	if de == child {
		return fmt.Errorf("%v: adding a node to itself as a child", ErrorHierarchyRequest)
	}

	// Uh, we can do type assertion, or this.
	if child.GetNodeType() == AttributeNode || child.GetNodeType() == DocumentNode {
		return fmt.Errorf("%v: an attempt was made to insert a node where it is not permitted", ErrorHierarchyRequest)
	}

	// Remove child from it's exisiting parent, if any.
	parent := child.GetParentNode()
	if parent != nil {
		parent.RemoveChild(child)
	}

	child.setParentNode(de)
	de.nodes = append(de.nodes, child)
	return nil
}

// RemoveChild removes the child node indicated by oldChild from the list of children of ref, and returns it.
// The returned error will be non nil in case the oldChild is not a child of the current Node.
func (de *domElement) RemoveChild(oldChild Node) (Node, error) {
	if oldChild == nil {
		return nil, nil
	}

	for i, child := range de.GetChildNodes() {
		if child == oldChild {
			// Slice trickery to remove the node at the found index:
			de.nodes = append(de.nodes[:i], de.nodes[i+1:]...)
			return child, nil
		}
	}

	return nil, ErrorNotFound
}

// ReplaceChild replaces the child node oldChild with newChild in the list of children, and
// returns the oldChild node. If newChild is a DocumentFragment object, oldChild is replaced
// by all of the DocumentFragment children, which are inserted in the same order. If the
// newChild is already in the tree, it is first removed.
func (de *domElement) ReplaceChild(newChild, oldChild Node) (Node, error) {
	if newChild == nil {
		return nil, fmt.Errorf("%v: given new child is nil", ErrorHierarchyRequest)
	}
	if oldChild == nil {
		return nil, fmt.Errorf("%v: given old child is nil", ErrorHierarchyRequest)
	}

	// newChild must be created by the same owner document of this element.
	if newChild.GetOwnerDocument() != de.GetOwnerDocument() {
		return nil, ErrorWrongDocument
	}

	// Find the old child, and replace it with the new child.
	for i, child := range de.GetChildNodes() {
		if child == oldChild {
			// Check if newChild has a parent (i.e., it's in the tree).
			ncParent := newChild.GetParentNode()
			if ncParent != nil {
				// Remove the newChild from its parent.
				ncParent.RemoveChild(newChild)
			}

			// Slice trickery, again. It will make a new underlying slice with one element,
			// the 'newChild', and then append the rest of the de.nodes to that.
			de.nodes = append(de.nodes[:i], append([]Node{newChild}, de.nodes[i+1:]...)...)
			// Change the parent node:
			newChild.setParentNode(de)

			return oldChild, nil
		}
	}

	return nil, ErrorNotFound
}

// InsertBefore inserts the Node newChild before the reference child, refChild.
// If the refChild is nil, the newChild will simply be appended at the end of
// the list of children.
func (de *domElement) InsertBefore(newChild, refChild Node) (Node, error) {
	if newChild == nil {
		// FIXME: what in this case? Is an error ok?
		return nil, ErrorHierarchyRequest
	}

	if newChild.GetNodeType() == AttributeNode {
		return nil, ErrorHierarchyRequest
	}

	// New child must have the same owner Document as this element's document.
	if newChild.GetOwnerDocument() != de.GetOwnerDocument() {
		return nil, ErrorWrongDocument
	}

	// If refChild is nil, append to the end, and return.
	if refChild == nil {
		err := de.AppendChild(newChild)
		if err != nil {
			return nil, err
		}
		return newChild, nil
	}

	// Find the reference child, insert newChild before that one.
	for i, child := range de.GetChildNodes() {
		if child == refChild {
			// Check if newChild is in the tree already. If so, remove it.
			ncParent := newChild.GetParentNode()
			if ncParent != nil {
				ncParent.RemoveChild(newChild)
			}
			newChild.setParentNode(de)
			de.nodes = append(de.nodes[:i], append([]Node{newChild}, de.nodes[i:]...)...)
			return newChild, nil
		}
	}

	// The reference child is given, but not found. We got no information where
	// to insert the newChild at.
	return nil, ErrorNotFound
}

func (de *domElement) HasChildNodes() bool {
	return len(de.nodes) > 0
}

func (de *domElement) GetPreviousSibling() Node {
	return getPreviousSibling(de)
}

func (de *domElement) GetNextSibling() Node {
	return getNextSibling(de)
}

func (de *domElement) GetNamespaceURI() string {
	return de.namespaceURI
}

func (de *domElement) GetNamespacePrefix() string {
	return de.tagName.GetPrefix()
}

func (de *domElement) GetTagName() string {
	return string(de.tagName)
}

func (de *domElement) SetAttribute(name, value string) error {
	attr, err := de.GetOwnerDocument().CreateAttribute(name)
	if err != nil {
		return err
	}

	var namespaceFound = false

	// Get the prefix (if any), double check if its declared. If it is not, then setting that
	// attribute generates an error. Attributes with prefix "xmlns" must always be able
	// to be set, and no lookup is necessary. Attributes without a prefix should be ok
	// as well.
	if attr.GetNamespacePrefix() == "xmlns" || attr.GetNamespacePrefix() == "" {
		namespaceFound = true
	} else if attr.GetNamespacePrefix() != "" {
		// Not xmlns, but a different prefix. Look it up, see if it's declared somewhere
		// up in the tree.
		_, namespaceFound = de.LookupNamespaceURI(attr.GetNamespacePrefix())
	}

	if !namespaceFound {
		return fmt.Errorf("the namespace for prefix '%v' has not been declared", attr.GetNamespacePrefix())
	}

	attr.SetValue(value)
	attr.setOwnerElement(de)
	de.attributes.SetNamedItem(attr)

	return nil
}

// SetAttributeNode adds a new attribute node. If an attribute with that name (nodeName) is
// already present in the element, it is replaced by the new one. Replacing an attribute node
// by itself has no effect. To add a new attribute node with a qualified name and namespace
// URI, use the SetAttributeNodeNS method.
func (de *domElement) SetAttributeNode(a Attr) error {
	// Attribute and Element must share the same owner document.
	if a.GetOwnerDocument() != de.GetOwnerDocument() {
		return ErrorWrongDocument
	}

	// Is the Attribute is already owned by another Element?
	if a.GetOwnerElement() != nil {
		return ErrorAttrInUse
	}

	a.setOwnerElement(de)
	de.attributes.SetNamedItem(a)
	return nil
}

// TODO: SetAttributeNodeNS

func (de *domElement) GetAttribute(name string) string {
	if theAttr := de.attributes.GetNamedItem(name); theAttr != nil {
		return theAttr.GetNodeValue()
	}

	// Not found, can return an empty string as per spec.
	return ""
}

// GetElementsByTagName finds all descendant element with the given tagname.
// This implementation does a recursive search.
func (de *domElement) GetElementsByTagName(tagname string) []Element {
	return getElementsBy(de, "", tagname, false)
}

func (de *domElement) GetElementsByTagNameNS(namespaceURI, tagname string) []Element {
	return getElementsBy(de, namespaceURI, tagname, true)
}

// setTagName is only used internally, when the tagname needs to change. One example is during parsing:
// The Go encoding/xml package does not directly take prefixes into account, so we do some hackery to
// make that work. After we found a prefx<->namespace match, we need to change the tagname.
//
// Is it assumed that the tagname is XML valid at that point. For now.
func (de *domElement) setTagName(tagname string) {
	de.tagName = XMLName(tagname)
}

// normalizeNamespaces normalizes namespace declaration attributes and prefixes, as part of the NormalizeDocument
// method of the Document interface.
func (de *domElement) normalizeNamespaces(counter *int) {
	parent := de.GetParentNode()
	// fmt.Printf("==> Parent node: %v\n", parent)
	// fmt.Printf("    Normalizing: %v\n", de)

	// Does the current element have a namespace URI defined?
	if de.GetNamespaceURI() != "" {
		// Never inherited if the parent is the Document itself. In that case,
		// this element is the root element.
		if parent.GetNodeType() == DocumentNode {
			de.removeNSDecl()
			de.createXmlnsDecl()
		} else {
			// Check if the declaration is inherited by looking up the prefix.
			ns, found := parent.LookupNamespaceURI(de.GetNamespacePrefix())
			if !found || ns != de.GetNamespaceURI() {
				// The namespace is inherited from a top level element.
				// fmt.Println("    Inherited?")
				de.removeNSDecl()
				de.createXmlnsDecl()
			} else {
				// fmt.Printf("    Found: %s, %s\n", ns, de.GetNamespacePrefix())
			}
			// Else, it's inherited from the parent, and nothing needs to be done.
		}
	}

	// fmt.Printf("==> Result:      %v\n", de)
	// fmt.Println()

	for _, c := range de.GetChildNodes() {
		if e, ok := c.(Element); ok {
			e.normalizeNamespaces(counter)
		}
	}
}

// LookupPrefix looks up the prefix associated to the given namespace URI, starting from this node.
// The default namespace declarations are ignored by this method. See Namespace Prefix Lookup for
// details on the algorithm used by this method:
// https://www.w3.org/TR/2004/REC-DOM-Level-3-Core-20040407/namespaces-algorithms.html#lookupNamespacePrefixAlgo
func (de *domElement) LookupPrefix(requestedNamespace string) (string, bool) {
	if requestedNamespace == "" {
		return "", false
	}

	// Check if the element has a namespace URI declared, and if there's a
	// namespace.
	elementPrefix := de.GetNamespacePrefix()
	lookedupNamespace, found := de.LookupNamespaceURI(elementPrefix)
	if de.GetNamespaceURI() == requestedNamespace &&
		elementPrefix != "" &&
		found &&
		lookedupNamespace == requestedNamespace {

		return elementPrefix, true
	}

	// Iterate over attributes with xmlns declarations.
	if de.GetAttributes() != nil {
		attrs := de.GetAttributes().GetItems()
		for _, node := range attrs {
			a := node.(Attr)
			attrpfx := a.GetNamespacePrefix() // xmlns : ... = .........
			attrloc := a.GetLocalName()       // ..... : pfx = .........
			attrval := a.GetNodeValue()       // ..... : ... = namespace

			if attrpfx == "xmlns" && attrval == requestedNamespace {
				return attrloc, true
			}
		}
	}

	// Nothing found in this element, maybe something is declared up in the tree?
	if parentElement, ok := de.GetParentNode().(Element); ok {
		return parentElement.LookupPrefix(requestedNamespace)
	}

	return "", false
}

// LookupNamespaceURI looks up the namespace URI belonging to the prefix pfx. See
// https://www.w3.org/TR/2004/REC-DOM-Level-3-Core-20040407/namespaces-algorithms.html#lookupNamespaceURIAlgo
// for more information on the implementation of this method.
func (de *domElement) LookupNamespaceURI(pfx string) (string, bool) {
	if de.GetNamespaceURI() != "" && de.GetNamespacePrefix() == pfx {
		return de.GetNamespaceURI(), true
	}

	// Check the element's xmlns declarations.
	if de.GetAttributes() != nil {
		attrs := de.GetAttributes().GetItems()
		for _, node := range attrs {
			a := node.(Attr)
			// <elem xmlns="..." />, and prefix is empty:
			if a.GetNodeName() == "xmlns" && pfx == "" {
				return a.GetNodeValue(), true
			}

			// <elem xmlnsanycharacter="..." />, and prefix is empty:
			//
			// This seems to be according to spec. Anything starting with xmlns is just a namespace declaration.
			// Xerces DOM also works like this.
			if strings.HasPrefix(a.GetNodeName(), "xmlns") && !strings.Contains(a.GetNodeName(), ":") && pfx == "" {
				return a.GetNodeValue(), true
			}

			// <pfx:elem xmlns:pfx="..." />, with a given prefix:
			//
			// First, get the last index of the 'xmlns:pfx' part. The node name can possibly contain multiple
			// colon characters, like 'xmlns:bla:cruft:pfx'. In the Xerces implementation of the DOM, this will
			// result in the local name 'pfx'.
			s := strings.LastIndex(a.GetNodeName(), ":")
			if strings.HasPrefix(a.GetNodeName(), "xmlns") && s >= 0 && a.GetNodeName()[s+1:] == pfx {
				return a.GetNodeValue(), true
			}
		}
	}

	// Found no declarations in the attributes of this element, therefore we check the ancestor. We must only check
	// if the parent element is an Element itself. If we don't, we can get in an infinite loop when the parent node
	// is a Document, since the Document will use the GetDocumentElement() to lookup the prefix.
	if parentElement, ok := de.GetParentNode().(Element); ok {
		return parentElement.LookupNamespaceURI(pfx)
	}

	// In the end, nothing is found.
	return "", false
}

func (de *domElement) IsDefaultNamespace(namespace string) bool {
	if de.GetNamespacePrefix() == "" {
		return de.GetNamespaceURI() == namespace
	}

	// TODO verify this loop
	for k, v := range de.GetAttributes().GetItems() {
		if k == "xmlns" {
			return v.GetNodeValue() == namespace
		}
	}

	if parentElement, ok := de.GetParentNode().(Element); ok {
		return parentElement.IsDefaultNamespace(namespace)
	}

	return false
}

func (de *domElement) GetTextContent() string {
	if !de.HasChildNodes() {
		return ""
	}

	textContent := ""
	for _, child := range de.GetChildNodes() {
		// Skip comments and PIs.
		if child.GetNodeType() == CommentNode || child.GetNodeType() == ProcessingInstructionNode {
			continue
		}

		textContent += child.GetTextContent()
	}

	return textContent
}

// SetTextContent will remove any possible children this node may have if the content string is not empty. The children
// will be replaced by a single Text node containing the content.
func (de *domElement) SetTextContent(content string) {
	if content == "" {
		return
	}
	// Remove existing nodes from this element by initializing an empty Node slice.
	de.nodes = make([]Node, 0)

	text := de.GetOwnerDocument().CreateText(content)
	de.AppendChild(text)
}

// CloneNode for Elements clones this element. If deep is set to false, it will create a clone of the
// Element, plus its attributes. If deep is set to true, it will create a clone of all its children (and so on).
func (de *domElement) CloneNode(deep bool) Node {
	// Clone element. The clone does not have a parent.
	cloneElement, err := de.ownerDocument.CreateElementNS(de.namespaceURI, string(de.tagName))
	if err != nil {
		panic("CreateElement returned an error, but should be impossible at this point")
	}
	// Then its attributes.
	for _, attrNode := range de.GetAttributes().GetItems() {
		cloneAttr := attrNode.CloneNode(deep).(Attr)
		cloneElement.SetAttributeNode(cloneAttr)
	}

	if !deep {
		// No deep clone, so return it.
		return cloneElement
	}
	// Do a deep clone.
	for _, child := range de.GetChildNodes() {
		childClone := child.CloneNode(true)
		cloneElement.AppendChild(childClone)
	}

	return cloneElement
}

func (de *domElement) ImportNode(n Node, deep bool) Node {
	return importNode(de.ownerDocument, n, deep)
}

// Private functions:
func (de *domElement) setParentNode(parent Node) {
	de.parentNode = parent
}

func (de *domElement) setOwnerDocument(doc Document) {
	de.ownerDocument = doc
}

// removeNSDeclAndSet finds xmlns:prefix declarations, and removes them. New namespace declarations will be
// created once we see we need them.
func (de *domElement) removeNSDecl() {
	for nsdecl := range de.GetAttributes().GetItems() {
		if strings.HasPrefix(nsdecl, "xmlns") {
			// remove it.
			//fmt.Printf("Removing xmlns attribute '%s'\n", nsdecl)
			de.GetAttributes().RemoveNamedItem(nsdecl)
		}
	}
}

func (de *domElement) createXmlnsDecl() {
	if de.GetNamespacePrefix() == "" {
		de.SetAttribute("xmlns", de.GetNamespaceURI())
	} else {
		de.SetAttribute("xmlns:"+de.GetNamespacePrefix(), de.GetNamespaceURI())
	}
}

func (de *domElement) String() string {
	return fmt.Sprintf("%s, <%s>, ns=%s, attrs={%v}",
		de.GetNodeType(), de.tagName, de.namespaceURI, de.attributes)
}
