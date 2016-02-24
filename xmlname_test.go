package dom

import (
	"testing"
)

func TestXMLNameIsValid(t *testing.T) {
	var tests = []struct {
		element        string
		expectedToBeOK bool
	}{
		{"valid", true},
		{" with_leading_space", false},
		{"with space", false},
	}
	_ = tests
}
