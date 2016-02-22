package dom

import (
	"errors"
)

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
	// XmlDeclaration is the usually default XML processing instruction at the
	// start of XML documents. This is merely added as a convenience. It's the
	// same declaration which the encoding/xml package has, except it does not
	// have a trailing newline.
	XmlDeclaration = `<?xml version="1.0" encoding="UTF-8"?>`
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

// NamedNodeMap represent collections of nodes that can be accessed by name.
type NamedNodeMap interface {
	GetNamedItem(string) Node
	SetNamedItem(Node) error
	GetItems() map[string]Node
	Length() int
}

// Node is the primary interface for the entire Document Object Model. It represents
// a single node in the document tree. While all objects implementing the Node
// interface expose methods for dealing with children, not all objects implementing
// the Node interface may have children.
type Node interface {
	// Gets the node name. Depending on the type (Attr, CDATASection, Element etc.)
	// the result of this call differs.
	NodeName() string
	// Gets the type of node.
	NodeType() NodeType
	// Gets the node value. Like NodeName(), the output differs depending on the type.
	NodeValue() string
	// Returns the local part of the qualified name of this node.
	LocalName() string
	// Gets the list of child nodes.
	NodeList() []Node
	// Gets the parent node. May be nil if none was assigned.
	ParentNode() Node
	// Gets the first child Node of this Node. May return nil if no child nodes
	// exist.
	FirstChild() Node
	// GetAttributes will return the attributes belonging to this node. In the current
	// spec, only Element nodes will return something sensible (i.e. non nil) when this
	// function is called.
	GetAttributes() NamedNodeMap
	// Gets the owner document (the Document instance which was used to create
	// the Node).
	OwnerDocument() Document
	// Appends a child to this Node. Will return an error when this Node is not
	// able to have any (more) children, like Text nodes.
	AppendChild(Node) error
	// Returns true when the Node has one or more children.
	HasChildNodes() bool
	// Returns the namespace URI of this node.
	NamespaceURI() string

	// Private functions
	setParentNode(Node)
	setOwnerDocument(Document)
	setNamespaceURI(string)
}

// Attr represents an attribute in an Element.
type Attr interface {
	Node

	GetName() string
	SetName(string)
	IsSpecified() bool
	GetValue() string
	SetValue(string)
	GetOwnerElement() Element

	setOwnerElement(Element)
}

// Element represents an element in an HTML or XML document.
type Element interface {
	Node

	// Sets the tag name of this element.
	SetTagName(tagname string)
	// Gets the tag name of this element.
	GetTagName() string

	SetAttribute(name, value string)
	GetAttribute(name string) string
}

// Text is a Node that represents character data.
type Text interface {
	Node

	GetData() string
	SetData(s string)
}

// DocumentType belongs to a Document, but can also be nil. The DocumentType
// interface in the DOM Core provides an interface to the list of entities
// that are defined for the document, and little else because the effect of
// namespaces and the various XML schema efforts on DTD representation are
// not clearly understood as of this writing. (Direct copy of the spec).
type DocumentType interface {
	Node

	// Gets the name of the DTD; i.e.  the name immediately following the DOCTYPE keyword.
	GetName() string
	// The public identifier of the external subset.
	GetPublicID() string
	// The system identifier of the external subset. This may be an absolute URI or not.
	GetSystemID() string
}

// Document is the root of the Document Object Model.
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
	CreateTextNode(string) Text
	// Creates an Attr of the given name and returns it.
	CreateAttribute(name string) (Attr, error)
	// CreateComment creates a Comment node with the given comment content. If
	// the comment contains a double hyphen (--), this should generate an error.
	CreateComment(comment string) (Comment, error)
	// Gets the document element, which should be the first (and only) child Node
	// of the Document. Can be nil if none is set yet.
	GetDocumentElement() Element
}

type Comment interface {
	Node

	GetComment() string
	SetComment(comment string)
}
