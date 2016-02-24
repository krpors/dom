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
		{"_______", true},
		{"\x00\x0A", false},
	}

	for _, test := range tests {
		elemName := XMLName(test.element)
		valid := elemName.IsValid()
		if valid != test.actual {
			t.Errorf("element '%v' is expected to be '%t', but was '%t'", elemName, test.actual, valid)
		}
	}
}
