package testutil

import "testing"

func TestNormalize(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "BOM and CRLF", in: "\uFEFFline1\r\n\r\nline2\r\n", want: "line1\n\nline2"},
		{name: "CR only and trailing spaces", in: "a b\r\n  \r\n", want: "a b"},
		{name: "multiple empty lines", in: "x\n\n\n\n y\n", want: "x\n\n y"},
		{name: "empty input", in: "", want: ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Normalize(c.in)
			if got != c.want {
				t.Fatalf("%s: got %q want %q", c.name, got, c.want)
			}
		})
	}
}
