package util

import "testing"

type testCase struct {
	utf8, other, otherEncoding string
}

var testCases = []testCase{
	{"Résumé", "Résumé", "utf8"},
	{"Résumé", "R\xe9sum\xe9", "latin1"},
	{"これは漢字です。", "S0\x8c0o0\"oW[g0Y0\x020", "UTF-16LE"},
	{"これは漢字です。", "0S0\x8c0oo\"[W0g0Y0\x02", "UTF-16BE"},
	{"Hello, world", "Hello, world", "ASCII"},
	{"Gdańsk", "Gda\xf1sk", "ISO-8859-2"},
}

func TestDecodeToUtf8(t *testing.T) {
	for _, tc := range testCases {

		actual, err := DecodeToUtf8(tc.other, tc.otherEncoding)
		if err != nil {
			t.Errorf("%s: encode %q: %s", tc.otherEncoding, tc.utf8, err)
			continue
		}

		if actual != tc.utf8 {
			t.Errorf("%s: got %q, want %q", tc.otherEncoding, actual, tc.other)
		}
	}
}

func TestEncodeUtf8(t *testing.T) {
	for _, tc := range testCases {

		actual, err := EncodeUtf8(tc.utf8, tc.otherEncoding)
		if err != nil {
			t.Errorf("%s: encode %q: %s", tc.otherEncoding, tc.utf8, err)
			continue
		}

		if actual != tc.other {
			t.Errorf("%s: got %q, want %q", tc.otherEncoding, actual, tc.other)
		}
	}
}
