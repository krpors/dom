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
	tagName string // The complete tagname given, with prefix.
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
	if index := strings.Index(de.tagName, ":"); index >= 0 {
		return de.tagName[index+1:]
	}
	return de.tagName
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

func (de *domElement) GetNamespacePrefix() string {
	// TODO: namespace prefix
	if index := strings.Index(de.tagName, ":"); index >= 0 {
		return de.tagName[0:index]
	}
	return ""
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
	attr.setName(name)
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
			if strings.HasPrefix(a.GetNodeName(), "xmlns") && s >= 0 && a.GetNodeName()[s:] == pfx {
				return a.GetNodeValue()
			}
		}
	}

	// Found no declarations in the attributes of this element, therefore we check the ancestor. This could be another
	// Element, or a Document. In any case, it's a Node.
	if de.GetParentNode() != nil {
		return de.GetParentNode().LookupNamespaceURI(pfx)
	}

	// In the end, nothing is found.
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
