package dom

import (
	"testing"
)

// TestXMLNameIsValid tests whether the IsValid function on XMLName behaves as
// expected.
func TestXMLNameIsValid(t *testing.T) {
	var tests = []struct {
		element string
		actual  bool
	}{
		{"", false},
		{"\xff\xfdand_the_rest", false},
		{"valid", true},
		{" with_leading_space", false},
		{"with space", false},
		{"åmal", true},
		{"Ball", true},
		{"éééé", true},
		{"!cruft", false},
		{"_underscore", true},
		{":element", true},
		{"abc012", true},
		{"other-element", true},
		{".dotstart", false},
		{"elem.with.dot", true},
		{"¾¾", false},
		{"Éomër-From-lord-ÖF.THERINGS", true},
		{"ALLCAPSSHOULDWORKASWELL", true},
		{"namepaceprefix:someelement", true},
		{"_______", true},
		{"\x00\x0A", false},
		{"ok\xff\xfd", false},
		{"pfx:some", true},
		// {"ns1:double:colon", false}, //  TODO: fix this one
		{"hi:erf\\asd", false},
	}

	for _, test := range tests {
		elemName := XMLName(test.element)
		valid := elemName.IsValid()
		if valid != test.actual {
			t.Errorf("element '%v' is expected to be '%t', but was '%t'", elemName, test.actual, valid)
		}
	}
}
