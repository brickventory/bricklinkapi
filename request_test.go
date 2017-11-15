package bricklinkapi

import (
	"net/http"
	"testing"
)

func TestGenerateBaseURL(t *testing.T) {
	testCases := []struct {
		desc   string
		method string
		uri    string
		params []string
		expS   string
	}{
		{desc: "testing URI without params",
			method: "GET",
			uri:    "https://foo.com",
			params: []string{},
			expS:   "GET&https%3A%2F%2Ffoo.com"}, // without trailing &
		{desc: "testing URI with params",
			method: "GET",
			uri:    "https://foo.com",
			params: []string{
				"token=abcd",
				"secret=1234",
			},
			expS: "GET&https%3A%2F%2Ffoo.com&secret%3D1234%26token%3Dabcd"}, // params are sorted
	}
	for _, tc := range testCases {
		req, _ := http.NewRequest(tc.method, tc.uri, nil)
		result := generateBaseURL(req, tc.params)
		if result != tc.expS {
			t.Errorf("\n%v, want: %v, got: %v\n", tc.desc, tc.expS, result)
		}
	}
}

func TestGenerateSignature(t *testing.T) {
	testCases := []struct {
		desc           string
		base           string
		consumerSecret string
		tokenSecret    string
		expS           string
	}{
		{desc: "testing simple test string",
			base:           "foobar",
			consumerSecret: "foo",
			tokenSecret:    "bar",
			expS:           "uTIVP8RyuDOKi71kmTPV3t8%2BIfw%3D"},
		{desc: "testing simple uri",
			base:           "GET&GET&https%3A%2F%2Ffoo.com",
			consumerSecret: "abcd",
			tokenSecret:    "1234",
			expS:           "0wM5ydwpQqPq0%2FOgU%2FVHLzmAroM%3D"},
		{desc: "testing simple uri with parameter",
			base:           "GET&GET&https%3A%2F%2Ffoo.com&foo%3Dbar",
			consumerSecret: "abcd",
			tokenSecret:    "1234",
			expS:           "ipi8jiHZZl7T8GxkATgpxCJI5Nk%3D"},
	}
	for _, tc := range testCases {
		result := generateSignature(tc.base, tc.consumerSecret, tc.tokenSecret)
		if result != tc.expS {
			t.Errorf("\n%v, want: %v, got: %v\n", tc.desc, tc.expS, result)
		}
	}
}

func TestEncode(t *testing.T) {
	testCases := []struct {
		desc string
		s    string
		expS string
	}{
		{desc: "testing", s: "A", expS: "A"},
		{desc: "testing", s: "1", expS: "1"},
		{desc: "testing", s: "foo", expS: "foo"},
		{desc: "testing", s: "foo bar", expS: "foo%20bar"},
		{desc: "testing", s: "foo=bar", expS: "foo%3Dbar"},
		{desc: "testing", s: "foo+bar", expS: "foo%2Bbar"},
		{desc: "testing", s: "https://foo.com", expS: "https%3A%2F%2Ffoo.com"},
	}
	for _, tc := range testCases {
		result := encode(tc.s)
		if result != tc.expS {
			t.Errorf("%v \"%v\", want: %v, got: %v\n", tc.desc, tc.s, tc.expS, result)
		}
	}
}

func TestEncodable(t *testing.T) {
	testCases := []struct {
		desc string
		b    string
		exp  bool
	}{
		{desc: "testing letter", b: "A", exp: false},
		{desc: "testing letter", b: "M", exp: false},
		{desc: "testing letter", b: "Z", exp: false},
		{desc: "testing letter", b: "a", exp: false},
		{desc: "testing letter", b: "m", exp: false},
		{desc: "testing letter", b: "z", exp: false},
		{desc: "testing digit", b: "0", exp: false},
		{desc: "testing digit", b: "4", exp: false},
		{desc: "testing digit", b: "9", exp: false},
		{desc: "testing sign", b: "-", exp: false},
		{desc: "testing sign", b: ".", exp: false},
		{desc: "testing sign", b: "_", exp: false},
		{desc: "testing sign", b: "~", exp: false},
		{desc: "testing sign", b: "+", exp: true},
		{desc: "testing sign", b: "&", exp: true},
		{desc: "testing sign", b: "=", exp: true},
		{desc: "testing sign", b: "/", exp: true},
	}
	for _, tc := range testCases {
		result := encodable(tc.b[0])
		if result != tc.exp {
			t.Errorf("%v \"%v\", want: %v, got: %v\n", tc.desc, tc.b, tc.exp, result)
		}
	}
}
