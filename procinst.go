package dom

import "fmt"

type domProcInst struct {
	ownerDocument Document
	parentNode    Node
	data          string
	target        string
}

func newProcInst(owner Document, target string, data string) ProcessingInstruction {
	// TODO: validation of target and data, such as invalid characters (<? and ?> etc)
	pi := &domProcInst{}
	pi.ownerDocument = owner
	pi.target = target
	pi.data = data
	return pi
}

func (pi *domProcInst) GetNodeName() string {
	return pi.GetTarget()
}

func (pi *domProcInst) GetNodeType() NodeType {
	return ProcessingInstructionNode
}

// GetNodeValue returns the processing instruction's data.
func (pi *domProcInst) GetNodeValue() string {
	return pi.GetData()
}

func (pi *domProcInst) GetLocalName() string {
	// TODO: what?
	return ""
}

func (pi *domProcInst) GetChildNodes() []Node {
	return nil
}

func (pi *domProcInst) GetParentNode() Node {
	return pi.parentNode
}

func (pi *domProcInst) GetFirstChild() Node {
	return nil
}

func (pi *domProcInst) GetLastChild() Node {
	return nil
}

func (pi *domProcInst) GetAttributes() NamedNodeMap {
	return nil
}

func (pi *domProcInst) HasAttributes() bool {
	return false
}

func (pi *domProcInst) GetOwnerDocument() Document {
	return pi.ownerDocument
}

func (pi *domProcInst) AppendChild(child Node) error {
	return ErrorHierarchyRequest
}

func (pi *domProcInst) RemoveChild(oldChild Node) (Node, error) {
	return nil, ErrorHierarchyRequest
}
func (pi *domProcInst) ReplaceChild(newChild, oldChild Node) (Node, error) {
	return nil, ErrorHierarchyRequest
}
func (pi *domProcInst) InsertBefore(newChild, refChild Node) (Node, error) {
	return nil, ErrorHierarchyRequest
}

func (pi *domProcInst) HasChildNodes() bool {
	return false
}

func (pi *domProcInst) GetPreviousSibling() Node {
	return getPreviousSibling(pi)
}
func (pi *domProcInst) GetNextSibling() Node {
	return getNextSibling(pi)
}

func (pi *domProcInst) GetNamespaceURI() string {
	return ""
}

func (pi *domProcInst) GetNamespacePrefix() string {
	return ""
}

func (pi *domProcInst) LookupPrefix(namespace string) (string, bool) {
	if pi.GetParentNode() != nil {
		return pi.LookupPrefix(namespace)
	}
	return "", false
}

func (pi *domProcInst) LookupNamespaceURI(pfx string) (string, bool) {
	if pi.GetParentNode() != nil {
		return pi.LookupNamespaceURI(pfx)
	}
	return "", false
}

func (pi *domProcInst) IsDefaultNamespace(namespace string) bool {
	// TODO ?
	return false
}

// ProcessingInstruction methods
func (pi *domProcInst) GetData() string {
	return pi.data
}

func (pi *domProcInst) GetTarget() string {
	return pi.target
}

func (pi *domProcInst) SetData(data string) {
	pi.data = data
}

func (pi *domProcInst) setTarget(target string) {
	pi.target = target
}

func (pi *domProcInst) GetTextContent() string {
	return ""
}

func (pi *domProcInst) SetTextContent(content string) {
	// no-op
}

func (pi *domProcInst) CloneNode(deep bool) Node {
	clonePi, err := pi.ownerDocument.CreateProcessingInstruction(pi.target, pi.data)
	if err != nil {
		panic("CreateProcessingInstruction returned an unexpected error")
	}
	return clonePi
}

func (pi *domProcInst) ImportNode(n Node, deep bool) Node {
	return nil
}

func (pi *domProcInst) setParentNode(parent Node) {
	pi.parentNode = parent
}

func (pi *domProcInst) setOwnerDocument(doc Document) {
	pi.ownerDocument = doc
}

func (pi *domProcInst) String() string {
	return fmt.Sprintf("%s: '%s'='%s'", pi.GetNodeType(), pi.target, pi.data)
}
