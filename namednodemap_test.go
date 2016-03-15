package dom

import "testing"

func TestNamedNodeMap(t *testing.T) {
	doc := NewDocument()

	nnm := newNamedNodeMap()
	if nnm.Length() != 0 {
		t.Error("expected length of 0")
	}

	comment, _ := doc.CreateComment("meh")
	err := nnm.SetNamedItem(comment)
	if err == nil {
		t.Error("expected error at this point, but got none")
	}

	attr, _ := doc.CreateAttribute("name")
	attr.SetValue("value")
	nnm.SetNamedItem(attr)

	if nnm.Length() != 1 {
		t.Errorf("expected length of 1, got %d", nnm.Length())
		t.FailNow()
	}

	ret := nnm.GetNamedItem("name")
	if v, ok := ret.(Attr); ok {
		if v.GetNodeValue() != "value" {
			t.Errorf("expected 'value', got '%v'", v.GetNodeValue())
		}
	} else {
		t.Error("type assertion for Attr failed")
		t.FailNow()
	}

	attrDuplicate, _ := doc.CreateAttribute("name")
	attrDuplicate.SetValue("dupe!")
	nnm.SetNamedItem(attrDuplicate)

	if nnm.Length() != 1 {
		t.Errorf("expected length of 1, got %d", nnm.Length())
		t.FailNow()
	}

	ret = nnm.GetNamedItem("name")
	if v, ok := ret.(Attr); ok {
		if v.GetNodeValue() != "dupe!" {
			t.Errorf("expected 'dupe!', got '%v'", v.GetNodeValue())
		}
	} else {
		t.Error("type assertion for Attr failed")
		t.FailNow()
	}

	m := nnm.GetItems()
	if len(m) != 1 {
		t.Errorf("expected length of 1, got %d", len(m))
		t.FailNow()
	}

	if v, found := m["name"]; found {
		if v.GetNodeValue() != "dupe!" {
			t.Errorf("expected 'dupe!' got '%v'", v.GetNodeValue())
		}
	} else {
		t.Error("expected to find key 'name', but got nothing")
	}
}
