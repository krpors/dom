package dom

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

// TODO: Return proper errors.
// NUSE_ATTRIBUTE_ERR: Raised if arg is an Attr that is already an attribute of
// another Element object. The DOM user must explicitly clone Attr nodes to
// re-use them in other elements.
// HIERARCHY_REQUEST_ERR: Raised if an attempt is made to add a node doesn't
// belong in this NamedNodeMap. Examples would include trying to insert something
// other than an Attr node into an Element's map of attributes, or a non-Entity
// node into the DocumentType's map of Entities.
func (nnm *domNamedNodeMap) SetNamedItem(n Node) error {
	nnm.nodes[n.GetNodeName()] = n
	return nil
}

func (nnm *domNamedNodeMap) Length() int {
	return len(nnm.nodes)
}
