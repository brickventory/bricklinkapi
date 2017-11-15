package bricklinkapi

import (
	"testing"
)

func TestStringInSlice(t *testing.T) {
	testCases := []struct {
		desc string
		s    string
		sl   []string
		exp  bool
	}{
		{desc: "testing simple slice for true",
			s:   "foo",
			sl:  []string{"foo", "bar", "baz"},
			exp: true,
		},
		{desc: "testing simple slice for false",
			s:   "test",
			sl:  []string{"foo", "bar", "baz"},
			exp: false,
		},
		{desc: "testing empty slice for false",
			s:   "test",
			sl:  []string{},
			exp: false,
		},
	}
	for _, tc := range testCases {
		ok := stringInSlice(tc.s, tc.sl)
		if ok != tc.exp {
			t.Errorf("\n%v, want: %v, got: %v\n", tc.desc, tc.exp, ok)
		}
	}
}
