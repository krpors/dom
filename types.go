package dom

import (
	"errors"
)

// This file contains the definitions of errors, interfaces and other constants
// of the DOM Level 3 spec.

// Error definitions:
var (
	// ErrorHierarchyRequest is the error which can be returned when the node
	// is of a type that does not allow children, if the node to append to is
	// one of this node's ancestors or this node itself, or if this node is of
	// type Document and the DOM application attempts to append a second
	// DocumentType or Element node.
	ErrorHierarchyRequest = errors.New("HIERARCHY_REQUEST_ERR: an attempt was made to insert a node where it is not permitted")

	// ErrorInvalidCharacter is returned when an invalid character is used for
	// for example an element or attribute name.
	ErrorInvalidCharacter = errors.New("INVALID_CHARACTER_ERR: an invalid or illegal XML character is specified")

	// ErrorNotSupported is returned when this implementation does not support
	// the requested operation or object.
	ErrorNotSupported = errors.New("NOT_SUPPORTED_ERR: this implementation does not support the requested type of object or operation")
)

var (
	// XMLDeclaration is the usually default XML processing instruction at the
	// start of XML documents. This is merely added as a convenience. It's the
	// same declaration which the encoding/xml package has, except it does not
	// have a trailing newline.
	XMLDeclaration = `<?xml version="1.0" encoding="UTF-8"?>`
)

// NodeType defines the types of nodes which exist in the DOM.
type NodeType uint8

// Enumeration of all types of Nodes in the DOM.
const (
	ElementNode NodeType = iota
	AttributeNode
	TextNode
	CDATASectionNode
	EntityReferenceNode
	EntityNode
	ProcessingInstructionNode
	CommentNode
	DocumentNode
	DocumentTypeNode
	DocumentFragmentNode
)

// String returns the string representation of the NodeType, using the default
// representation by the W3 specification.
func (n NodeType) String() string {
	switch n {
	case ElementNode:
		return "ELEMENT_NODE"
	case AttributeNode:
		return "ATTRIBUTE_NODE"
	case TextNode:
		return "TEXT_NODE"
	case CDATASectionNode:
		return "CDATA_SECTION_NODE"
	case EntityReferenceNode:
		return "ENTITY_REFERENCE_NODE"
	case EntityNode:
		return "ENTITY_NODE"
	case ProcessingInstructionNode:
		return "PROCESSING_INSTRUCTION_NODE"
	case CommentNode:
		return "COMMENT_NODE"
	case DocumentNode:
		return "DOCUMENT_NODE"
	case DocumentTypeNode:
		return "DOCUMENT_TYPE_NODE"
	case DocumentFragmentNode:
		return "DOCUMENT_FRAGMENT_NODE"
	default:
		return "???"
	}
}

// Node is the primary interface for the entire Document Object Model. It represents
// a single node in the document tree. While all objects implementing the Node
// interface expose methods for dealing with children, not all objects implementing
// the Node interface may have children.
type Node interface {
	// Gets the node name. Depending on the type (Attr, CDATASection, Element etc.)
	// the result of this call differs.
	GetNodeName() string
	// Gets the type of node.
	GetNodeType() NodeType
	// Gets the node value. Like GetNodeName(), the output differs depending on the type.
	GetNodeValue() string
	// Returns the local part of the qualified name of this node.
	GetLocalName() string
	// Gets the list of child nodes.
	GetChildNodes() []Node
	// Gets the parent node. May be nil if none was assigned.
	GetParentNode() Node
	// Gets the first child Node of this Node. May return nil if no child nodes exist.
	GetFirstChild() Node
	// GetAttributes will return the attributes belonging to this node. In the current
	// spec, only Element nodes will return something sensible (i.e. non nil) when this
	// function is called.
	GetAttributes() NamedNodeMap
	// Gets the owner document (the Document instance which was used to create the Node).
	GetOwnerDocument() Document
	// Appends a child to this Node. Will return an error when this Node is not
	// able to have any (more) children, like Text nodes.
	AppendChild(Node) error
	// Returns true when the Node has one or more children.
	HasChildNodes() bool
	// GetPreviousSibling gets the Node immediately preceding this Node. If there is no such
	// node, this method returns nil.
	GetPreviousSibling() Node
	// GetNextSibling gets the Node immediately following this Node. If there is no such node,
	// this methods returns nil.
	GetNextSibling() Node
	// Returns the namespace URI of this node.
	GetNamespaceURI() string
	// GetNamespacePrefix returns the prefix of this node, or an empty string if it
	// does not have a prefix.
	GetNamespacePrefix() string
	// LookupPrefix up the prefix associated to the given namespace URI, starting from this node.
	// The default namespace declarations are ignored by this method.
	LookupPrefix(namespace string) string
	// LookupNamespaceURI looks up the namespace URI associated to the given prefix, starting
	// from this node. See Namespace Prefix Lookup for details on the algorithm used by this method:
	// https://www.w3.org/TR/2004/REC-DOM-Level-3-Core-20040407/namespaces-algorithms.html#lookupNamespacePrefixAlgo
	LookupNamespaceURI(pfx string) string
	// TODO: SetTextContent(string) implementation.
	// SetTextContent sets the text content of the current node. On setting, any
	// possible children this node may have are removed and, if the new string
	// is not empty, replaced by a single Text node containing the string this
	// attribute is set to.
	// SetTextContent(string)

	// Private functions
	setParentNode(Node)
	setNamespaceURI(string)
}

// ProcessingInstruction interface represents a "processing instruction", used
// in XML as a way to keep processor-specific information in the text of the document.
type ProcessingInstruction interface {
	Node

	// The content of this processing instruction. This is from the first non white
	// space character after the target to the character immediately preceding the ?>.
	// Target can be anything except the [XxMmLl] string.
	GetTarget() string
	// The target of this processing instruction. XML defines this as being the first
	// token following the markup that begins the processing instruction.
	GetData() string
}

// Attr represents an attribute in an Element. It implements the Node interface.
type Attr interface {
	Node

	GetName() string
	IsSpecified() bool
	GetValue() string
	SetValue(string)
	GetOwnerElement() Element

	setName(string)
	setOwnerElement(Element)
}

// Element represents an element in an HTML or XML document. It implements the Node interface.
type Element interface {
	Node

	// Sets the tag name of this element.
	SetTagName(tagname string)
	// Gets the tag name of this element.
	GetTagName() string
	// Convenience function to add an attribute.
	SetAttribute(name, value string)
	// Convenience function to get an attribute value.
	GetAttribute(name string) string
	// Sets an attribute based on the Attr type.
	SetAttributeNode(a Attr)
	// GetElementsByTagName finds all descendant elements of the current element,
	// with the given tag name, in document order.
	GetElementsByTagName(string) []Element
	// GetElementsByTagNameNS finds all descendant elements of the current element,
	// with the given tag name and namespace URI, in document order.
	GetElementsByTagNameNS(namespaceURI, tagname string) []Element
}

// Text represents character data within an element. It implements the Node interface.
type Text interface {
	Node

	// GetData gets the character data of this Text node.
	GetData() string
	// SetData sets the character data of this Text node.
	SetData(s string)
}

// DocumentType belongs to a Document, but can also be nil. The DocumentType
// interface in the DOM Core provides an interface to the list of entities
// that are defined for the document, and little else because the effect of
// namespaces and the various XML schema efforts on DTD representation are
// not clearly understood as of this writing. (Direct copy of the spec).
type DocumentType interface {
	Node

	// GetName gets the name of the DTD; i.e.  the name immediately following the DOCTYPE keyword.
	GetName() string
	// GetPublicID returns public identifier of the external subset.
	GetPublicID() string
	// GetSystemID returns the system identifier of the external subset. This may be an absolute URI or not.
	GetSystemID() string
}

// Document is the root of the Document Object Model. It implements the Node interface.
type Document interface {
	Node

	// Creates an element with the given tagname and returns it. Will return
	// an ErrorInvalidCharacter if the specified name is not an XML name according
	// to the XML version in use, specified in the XMLVersion attribute.
	CreateElement(tagName string) (Element, error)
	// Creates an element of the givens qualified name and namespace URI, and
	// returns it. Use an empty string if no namespace is necessary. See
	// CreateElement(string).
	CreateElementNS(namespaceURI, tagName string) (Element, error)
	// Creates a Text node given the specified string and returns it.
	CreateText(string) Text
	// Creates an Attr of the given name and returns it.
	CreateAttribute(name string) (Attr, error)
	// Creates an Attr using the given namespace URI and name.
	CreateAttributeNS(namespaceURI, name string) (Attr, error)
	// CreateComment creates a Comment node with the given comment content. If
	// the comment contains a double hyphen (--), this should generate an error.
	CreateComment(comment string) (Comment, error)
	// CreateProcessingInstruction creates a processing instruction and returns it.
	CreateProcessingInstruction(target, data string) (ProcessingInstruction, error)
	// Gets the document element, which should be the first (and only) child Node
	// of the Document. Can be nil if none is set yet.
	GetDocumentElement() Element
	// GetElementsByTagName finds all descendant elements of the current element,
	// with the given tag name, in document order.
	GetElementsByTagName(string) []Element
	// GetElementsByTagNameNS finds all descendant elements of the current element,
	// with the given tag name and namespace URI, in document order.
	GetElementsByTagNameNS(namespaceURI, tagname string) []Element
}

// Comment represents a comment node in an XML tree (e.g. <!-- ... -->). It implements
// the Node interface.
type Comment interface {
	Node

	// GetComment gets the comment text of this node.
	GetComment() string
	// SetComment gets the comment text of this node.
	SetComment(comment string)
}

// NamedNodeMap represent collections of nodes that can be accessed by name.
type NamedNodeMap interface {
	GetNamedItem(string) Node
	SetNamedItem(Node) error
	GetItems() map[string]Node
	Length() int
}
