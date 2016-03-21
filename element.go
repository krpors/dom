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

func (de *domElement) GetAttributes() NamedNodeMap {
	return de.attributes
}

func (de *domElement) GetOwnerDocument() Document {
	return de.ownerDocument
}

func (de *domElement) AppendChild(child Node) error {
	// TODO: if an Attr is attempted to be appended, return error.
	if de == child {
		return fmt.Errorf("%v: adding a node to itself as a child", ErrorHierarchyRequest)
	}

	// Uh, we can do type assertion, or this.
	if child.GetNodeType() == AttributeNode {
		return fmt.Errorf("%v: an attempt was made to insert a node where it is not permitted", ErrorHierarchyRequest)
	}

	child.setParentNode(de)
	de.nodes = append(de.nodes, child)
	return nil
}

func (de *domElement) RemoveChild(oldChild Node) (Node, error) {
	panic("not implemented yet")
}
func (de *domElement) ReplaceChild(oldChild Node) (Node, error) {
	panic("not implemented yet")
}
func (de *domElement) InsertBefore(newChild, refChild Node) (Node, error) {
	panic("not implemented yet")
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

func (de *domElement) SetAttribute(name, value string) {
	if de.attributes == nil {
		de.attributes = newNamedNodeMap()
	}

	attr := newAttr(de.GetOwnerDocument(), name, "")
	attr.SetValue(value)
	attr.setOwnerElement(de)
	de.attributes.SetNamedItem(attr)
}

// SetAttributeNode adds a new attribute node. If an attribute with that name (nodeName) is
// already present in the element, it is replaced by the new one. Replacing an attribute node
// by itself has no effect. To add a new attribute node with a qualified name and namespace
// URI, use the setAttributeNodeNS method.
// TODO: implement above
func (de *domElement) SetAttributeNode(a Attr) {
	if de.attributes == nil {
		de.attributes = newNamedNodeMap()
	}

	a.setOwnerElement(de)
	de.attributes.SetNamedItem(a)
}

func (de *domElement) GetAttribute(name string) string {
	if de.attributes == nil {
		// TODO: no attributes, return empty string??
		return ""
	}
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

// LookupPrefix looks up the prefix associated to the given namespace URI, starting from this node.
// The default namespace declarations are ignored by this method. See Namespace Prefix Lookup for
// details on the algorithm used by this method:
// https://www.w3.org/TR/2004/REC-DOM-Level-3-Core-20040407/namespaces-algorithms.html#lookupNamespacePrefixAlgo
func (de *domElement) LookupPrefix(namespace string) string {
	if namespace == "" {
		return ""
	}

	// Check if the element has a namespace URI declared, and if there's a
	// namespace.
	pfx := de.GetNamespacePrefix()
	if de.GetNamespaceURI() == namespace && pfx != "" {
		return pfx
	}

	// Iterate over attributes with xmlns declarations.
	if de.GetAttributes() != nil {
		attrs := de.GetAttributes().GetItems()
		for _, node := range attrs {
			a := node.(Attr)
			attrpfx := a.GetNamespacePrefix() // xmlns : ... = .........
			attrloc := a.GetLocalName()       // ..... : pfx = .........
			attrval := a.GetNodeValue()       // ..... : ... = namespace

			if attrpfx == "xmlns" && attrval == namespace {
				return attrloc
			}
		}
	}

	// Nothing found in this element, maybe something is declared up in the tree?
	if parentElement, ok := de.GetParentNode().(Element); ok {
		return parentElement.LookupPrefix(namespace)
	}

	return ""
}

// LookupNamespaceURI looks up the namespace URI belonging to the prefix pfx. See
// https://www.w3.org/TR/2004/REC-DOM-Level-3-Core-20040407/namespaces-algorithms.html#lookupNamespaceURIAlgo
// for more information on the implementation of this method.
func (de *domElement) LookupNamespaceURI(pfx string) string {
	if de.GetNamespaceURI() != "" && de.GetNamespacePrefix() == pfx {
		return de.GetNamespaceURI()
	}

	// Check the element's xmlns declarations.
	if de.GetAttributes() != nil {
		attrs := de.GetAttributes().GetItems()
		for _, node := range attrs {
			a := node.(Attr)
			// <elem xmlns="..." />, and prefix is empty:
			if a.GetNodeName() == "xmlns" && pfx == "" {
				return a.GetNodeValue()
			}

			// <elem xmlnsanycharacter="..." />, and prefix is empty:
			//
			// This seems to be according to spec. Anything starting with xmlns is just a namespace declaration.
			// Xerces DOM also works like this.
			if strings.HasPrefix(a.GetNodeName(), "xmlns") && !strings.Contains(a.GetNodeName(), ":") && pfx == "" {
				return a.GetNodeValue()
			}

			// <pfx:elem xmlns:pfx="..." />, with a given prefix:
			//
			// First, get the last index of the 'xmlns:pfx' part. The node name can possibly contain multiple
			// colon characters, like 'xmlns:bla:cruft:pfx'. In the Xerces implementation of the DOM, this will
			// result in the local name 'pfx'.
			s := strings.LastIndex(a.GetNodeName(), ":")
			if strings.HasPrefix(a.GetNodeName(), "xmlns") && s >= 0 && a.GetNodeName()[s+1:] == pfx {
				return a.GetNodeValue()
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
	return ""
}

// Private functions:
func (de *domElement) setParentNode(parent Node) {
	de.parentNode = parent
}

func (de *domElement) String() string {
	return fmt.Sprintf("%s, <%s>, ns=%s, attrs=%v",
		de.GetNodeType(), de.tagName, de.namespaceURI, de.attributes)
}
