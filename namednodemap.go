package dom

import "fmt"

type domNamedNodeMap struct {
	nodes map[string]Node
}

func newNamedNodeMap() NamedNodeMap {
	nnm := &domNamedNodeMap{}
	nnm.nodes = make(map[string]Node)
	return nnm
}

// TODO: GetItems()... not according to spec though.
func (nnm *domNamedNodeMap) GetItems() map[string]Node {
	return nnm.nodes
}

func (nnm *domNamedNodeMap) GetNamedItem(name string) Node {
	return nnm.nodes[name]
}

// TODO: Return proper errors. The current implementation is not entirely correct.
// The spec mentions that a NamedNodeMap can accept Nodes other than Attr, except
// when the 'parent' is an Element. We'll need specializations of the NamedNodeMap
// interface, for ex. attrNamedNodeMap, entityNamedNodeMap.
//
// See: https://www.w3.org/TR/2004/REC-DOM-Level-3-Core-20040407/core.html#ID-1780488922
//
// Errors:
//
// NUSE_ATTRIBUTE_ERR: Raised if arg is an Attr that is already an attribute of
// another Element object. The DOM user must explicitly clone Attr nodes to
// re-use them in other elements.
//
// HIERARCHY_REQUEST_ERR: Raised if an attempt is made to add a node doesn't
// belong in this NamedNodeMap. Examples would include trying to insert something
// other than an Attr node into an Element's map of attributes, or a non-Entity
// node into the DocumentType's map of Entities.
func (nnm *domNamedNodeMap) SetNamedItem(n Node) error {
	if _, ok := n.(Attr); ok {
		nnm.nodes[n.GetNodeName()] = n
		return nil
	}
	return fmt.Errorf("%v: can not set a non-Attr node as a named item", ErrorHierarchyRequest)
}

func (nnm *domNamedNodeMap) RemoveNamedItem(name string) {
	delete(nnm.nodes, name)
}

func (nnm *domNamedNodeMap) Length() int {
	return len(nnm.nodes)
}

func (nnm *domNamedNodeMap) String() string {
	s := ""
	for k, v := range nnm.GetItems() {
		s += fmt.Sprintf("%v=%v,", k, v.GetNodeValue())
	}
	return s
}
