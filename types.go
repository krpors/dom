package dom

import (
	"errors"
)

// This file contains the definitions of errors, interfaces and other constants
// of the DOM Level 3 spec.

// ErrorHierarchyRequest is the error which can be returned when the node
// is of a type that does not allow children, if the node to append to is
// one of this node's ancestors or this node itself, or if this node is of
// type Document and the DOM application attempts to append a second
// DocumentType or Element node.
var ErrorHierarchyRequest = errors.New("HIERARCHY_REQUEST_ERR: an attempt was made to insert a node where it is not permitted")

// ErrorInvalidCharacter is returned when an invalid character is used for
// for example an element or attribute name.
var ErrorInvalidCharacter = errors.New("INVALID_CHARACTER_ERR: an invalid or illegal XML character is specified")

// ErrorNotSupported is returned when this implementation does not support
// the requested operation or object.
var ErrorNotSupported = errors.New("NOT_SUPPORTED_ERR: this implementation does not support the requested type of object or operation")

// ErrorNotFound is returned when a specified Node is not found, for instance
// during an attempt to delete a child Node from another Node.
var ErrorNotFound = errors.New("NOT_FOUND_ERR: the given child is not found in the current context")

// ErrorWrongDocument is returned when an insertion is attempted of a Node which was
// created from a different document instance.
var ErrorWrongDocument = errors.New("WRONG_DOCUMENT_ERR: the child was created from a different Document instance")

// ErrorAttrInUse is returned when an attribute is already an attribute of another Element object.
// The DOM user must explicitly create/clone Attr nodes to re-use them in other elements.
var ErrorAttrInUse = errors.New("INUSE_ATTRIBUTE_ERR: the attribute is already an attribute of another Element")

// XMLDeclaration is the usually default XML processing instruction at the
// start of XML documents. This is merely added as a convenience. It's the
// same declaration which the encoding/xml package has, except it does not
// have a trailing newline.
var XMLDeclaration = `<?xml version="1.0" encoding="UTF-8"?>`

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
	GetNodeName() string   // Returns the node name.
	GetNodeType() NodeType // Returns the node type (e.g. ElementNode, CommentNode, ...)
	GetNodeValue() string  // Returns the node value.
	GetLocalName() string  // Returns the local part of the node name.

	GetParentNode() Node // Gets the parent Node, or nil of the Node has no parent.

	GetAttributes() NamedNodeMap // Returns the NamedNodeMap containing attributes, if any.
	HasAttributes() bool         // Returns true if the Node has any attributes, false if othwerwise.

	GetOwnerDocument() Document // Gets the owner document (the Document instance which was used to create the Node).

	GetChildNodes() []Node                              // Gets the list of child nodes this Node has, if any
	GetFirstChild() Node                                // Gets the first child Node of this Node. May return nil if no child nodes exist.
	GetLastChild() Node                                 // Gets the last child Node of this Node. May return nil if there is no such Node.
	AppendChild(Node) error                             // Appends a child to this Node. Returns an error when the Node does not allow the child Node.
	RemoveChild(oldChild Node) (Node, error)            // Removes oldChild and returns it.
	ReplaceChild(newChild, oldChild Node) (Node, error) // Replaces oldChild with newChild, and returns newChild.
	InsertBefore(newChild, refChild Node) (Node, error) // Inserts newChild before the refChild.
	HasChildNodes() bool                                // Returns true when the Node has one or more children.

	CloneNode(deep bool) Node          // Creates a duplicate of the current node.
	ImportNode(n Node, deep bool) Node // Imports a node from another document to this document, without altering or removing the source node from the original document.

	GetPreviousSibling() Node // Gets the Node immediately preceding this Node. Returns nil if no previous sibling exists.
	GetNextSibling() Node     // Gets the Node immediately following this Node. Returns nil if no following sibling exists.

	GetNamespaceURI() string    // Returns the namespace URI of this node.
	GetNamespacePrefix() string // Returns the namespace prefix of this node.

	LookupPrefix(namespace string) (string, bool) // Look up the prefix associated to the given namespace URI, starting from this node.
	LookupNamespaceURI(pfx string) (string, bool) // LookupNamespaceURI looks up the namespace URI associated to the given prefix.

	GetTextContent() string // Gets the text content of the current Node.
	SetTextContent(string)  // Sets the text content of the current Node. Any possible children are removed.

	setParentNode(Node)        // Sets the parent node of this Node.
	setOwnerDocument(Document) // Sets the owner document of the Node. Used by ImportNode() for example.
}

// ProcessingInstruction interface represents a "processing instruction", used
// in XML as a way to keep processor-specific information in the text of the document.
// A processing instruction has the following form in an XML document:
//	<?target data?>
// The target of a processing instruction can be anything except the string [XxMmLl].
// The data can be anything, except the string ?> since that denotes the end of the
// processing instruction. If that happens, a fatal error should occur.
type ProcessingInstruction interface {
	Node

	GetTarget() string   // Gets the target of the processing instruction.
	GetData() string     // Gets the data of the processing instruction.
	SetData(data string) // Sets the data.
}

// Attr represents an attribute in an Element. It implements the Node interface.
type Attr interface {
	Node

	GetName() string
	IsSpecified() bool
	GetValue() string
	SetValue(string)
	GetOwnerElement() Element

	setOwnerElement(Element) // setOwnerElement is necessary to add an owner after creation.
	setName(string)          // setName sets the attribute name. Used for normalizing attributes.
}

// Element represents an element in an HTML or XML document. It implements the Node interface.
type Element interface {
	Node

	GetTagName() string                                            // Gets the tag name of this element.
	SetAttribute(name, value string) error                         // Convenience function to add an attribute.
	SetAttributeNode(a Attr) error                                 // Sets an attribute based on the Attr type.
	GetAttribute(name string) string                               // Convenience function to get an attribute value.
	GetElementsByTagName(string) []Element                         // Find all descendant elements of the current element.
	GetElementsByTagNameNS(namespaceURI, tagname string) []Element // Like GetElementsByTagName, except with a namespace URI.

	setTagName(string)                // Sets the tagname when necessary.
	normalizeNamespaces(counter *int) // Normalizes namespaces. See https://www.w3.org/TR/2004/REC-DOM-Level-3-Core-20040407/namespaces-algorithms.html#normalizeDocumentAlgo
}

// Text represents character data within an element. It implements the Node interface.
// Note that the name of the methods defined on this interface are not aligned with the specifications,
// due to the fact the Go's interfaces will not see a correct difference between this Text
// interface, or the ProcessingInstruction interface when the methods have the same signatures.
type Text interface {
	Node

	GetText() string  // Gets the character data of this Text node.
	SetText(s string) // Sets the character data of this Text node.

	IsElementContentWhitespace() bool // Return true if the Text node contains "ignorable whitespace".
}

// DocumentType belongs to a Document, but can also be nil. The DocumentType
// interface in the DOM Core provides an interface to the list of entities
// that are defined for the document, and little else because the effect of
// namespaces and the various XML schema efforts on DTD representation are
// not clearly understood as of this writing. (Direct copy of the spec).
type DocumentType interface {
	Node

	GetName() string     // Gets the name of the DTD; i.e. the name immediately following the DOCTYPE keyword.
	GetPublicID() string // Returns the public identifier of the external subset.
	GetSystemID() string // Returns the system identifier of the external subset. This may be an absolute URI or not.
}

// Document is the root of the Document Object Model. It implements the Node interface. As per the spec,
// all child nodes must be created through an instance of a Document object.
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

	NormalizeDocument() // Puts the Document in 'normal form'.
}

// Comment represents a comment node in an XML tree (e.g. <!-- ... -->). It implements
// the Node interface.
type Comment interface {
	Node

	GetComment() string        // Returns the comment text of this node.
	SetComment(comment string) // Sets the comment text of this node.
}

// NamedNodeMap represents collections of nodes that can be accessed by name.
type NamedNodeMap interface {
	GetNamedItem(string) Node  // Gets a named item identified by the given string. Returns nil if nothing is found.
	SetNamedItem(Node) error   // Adds a new item. The node's NodeName is used as a key.
	RemoveNamedItem(string)    // Removes the item identified by the given string.
	GetItems() map[string]Node // Gets the items as a Go map.
	Length() int               // Gets the amount of items in the named node map.
}

// Configuration contains fields which can control the output of the Parser
// and Serializer. Note that not (all configuration are specified or used (yet).
type Configuration struct {
	CDataSections            bool   // Keep CDataSection Nodes in the Document.
	Comments                 bool   // Keep Comment nodes in the Document.
	ElementContentWhitespace bool   // Keep all whitespaces in the Document.
	Namespaces               bool   // Perform namespace processing as defined in https://www.w3.org/TR/2004/REC-DOM-Level-3-Core-20040407/namespaces-algorithms.html#normalizeDocumentAlgo
	NamespaceDeclarations    bool   // Include (true) or discard (false) namespace declaration attributes.
	NormalizeCharacters      bool   // Perform or do not perform character normalization.
	OmitXMLDeclaration       bool   // Omits XML declaration during serialization. Default: false.
	PrettyPrint              bool   // Pretty print during serialization. Default: false.
	IndentCharacter          string // Indent character, if pretty printing. Default is four spaces.
}

// NewConfiguration creates a Configuration object with the defaults as per the DOM spec.
func NewConfiguration() Configuration {
	return Configuration{
		CDataSections:            true,
		Comments:                 true,
		ElementContentWhitespace: true,
		Namespaces:               true,
		NamespaceDeclarations:    true,
		NormalizeCharacters:      false,
		OmitXMLDeclaration:       false,
		PrettyPrint:              false,
		IndentCharacter:          "    ",
	}
}
