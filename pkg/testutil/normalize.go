package testutil

import "strings"

// Normalize applies common normalizations to avoid false diffs:
// - remove UTF-8 BOM
// - normalize CRLF -> LF
// - trim trailing spaces on lines
// - collapse multiple consecutive empty lines into one
// - trim leading/trailing empty lines
func Normalize(s string) string {
	// remove BOM
	s = strings.TrimPrefix(s, "\uFEFF")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	lines := strings.Split(s, "\n")
	var out []string
	emptySeq := 0
	for _, ln := range lines {
		// trim trailing spaces/tabs
		ln = strings.TrimRight(ln, " \t")
		if ln == "" {
			emptySeq++
			if emptySeq > 1 {
				// collapse multiple empty lines
				continue
			}
		} else {
			emptySeq = 0
		}
		out = append(out, ln)
	}
	// remove leading/trailing empty lines
	for len(out) > 0 && out[0] == "" {
		out = out[1:]
	}
	for len(out) > 0 && out[len(out)-1] == "" {
		out = out[:len(out)-1]
	}

	return strings.Join(out, "\n")
}
