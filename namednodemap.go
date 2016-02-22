package dom

type domNamedNodeMap struct {
}

func newNamedNodeMap() NamedNodeMap {
	nnm := &domNamedNodeMap{}
	return nnm
}

func (nnm *domNamedNodeMap) GetItems() map[string]Node {
	return nil
}

func (nnm *domNamedNodeMap) GetNamedItem(name string) Node {
	return nil
}

func (nnm *domNamedNodeMap) SetNamedItem(n Node) error {
	return nil
}

func (nnm *domNamedNodeMap) Length() int {
	return 0
}
