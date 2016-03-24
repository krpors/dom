package dom

import "fmt"

type domProcInst struct {
	ownerDocument Document
	data          string
	target        string
}

func newProcInst(owner Document, target string, data string) ProcessingInstruction {
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
	return pi.ownerDocument
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

func (pi *domProcInst) LookupPrefix(namespace string) string {
	if pi.GetParentNode() != nil {
		return pi.LookupPrefix(namespace)
	}
	return ""
}

func (pi *domProcInst) LookupNamespaceURI(pfx string) string {
	if pi.GetParentNode() != nil {
		return pi.LookupNamespaceURI(pfx)
	}
	return ""
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
	// TODO: implement
	return ""
}

func (pi *domProcInst) SetTextContent(content string) {
	// TODO: implement
}

func (pi *domProcInst) setParentNode(parent Node) {
	// no-op
}

func (pi *domProcInst) String() string {
	return fmt.Sprintf("%s: '%s'='%s'", pi.GetNodeType(), pi.target, pi.data)
}
