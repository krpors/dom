package dom

import "fmt"

type domProcInst struct {
	ownerDocument Document
	data          string
	target        string
}

func newProcInst() ProcessingInstruction {
	pi := &domProcInst{}
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

func (pi *domProcInst) GetAttributes() NamedNodeMap {
	return nil
}

func (pi *domProcInst) GetOwnerDocument() Document {
	return pi.ownerDocument
}

func (pi *domProcInst) AppendChild(child Node) error {
	return ErrorHierarchyRequest
}

func (pi *domProcInst) HasChildNodes() bool {
	return false
}

func (pi *domProcInst) GetNamespaceURI() string {
	return ""
}

func (pi *domProcInst) GetData() string {
	return pi.data
}

func (pi *domProcInst) GetTarget() string {
	return pi.target
}

// Private functions of Node:
func (pi *domProcInst) setNamespaceURI(s string) {
}

func (pi *domProcInst) setOwnerDocument(d Document) {
	pi.ownerDocument = d
}

func (pi *domProcInst) setParentNode(parent Node) {
	// no-op
}

func (pi *domProcInst) String() string {
	return fmt.Sprintf("%s: '%s'='%s'", pi.GetNodeType(), pi.target, pi.data)
}
