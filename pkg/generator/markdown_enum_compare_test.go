package generator

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"ones-cfg2md/pkg/model"
	"ones-cfg2md/pkg/parser"
	"ones-cfg2md/pkg/testutil"
)

func loadReference(t *testing.T) string {
	t.Helper()
	// compute path relative to this test file's testdata
	_, thisFile, _, _ := runtime.Caller(0)
	refPath := filepath.Join(filepath.Dir(thisFile), "..", "..", "fixtures", "output", "Перечисление_СостоянияЗаказов.md")
	data, err := os.ReadFile(refPath)
	if err != nil {
		t.Fatalf("failed to read reference markdown: %v", err)
	}
	// normalize newlines and remove BOM
	return testutil.Normalize(string(data))
}

// legacy test removed; use unified TestMarkdownFromFixtureMatchesReference

// TestMarkdownFromFixtureMatchesReference runs the same comparison for both EDT and CFG
// formats using subtests so the test names are identical for both formats.
func TestMarkdownFromFixtureMatchesReference(t *testing.T) {
	// compute testdata directory relative to this test file to avoid depending on working directory
	// parser fixtures are kept under pkg/parser/testdata/input
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..", "fixtures")
	repoRoot, _ = filepath.Abs(repoRoot)

	ref := loadReference(t)

	cases := []struct {
		name string
		kind string // "edt" or "cfg"
	}{
		{name: "EDT", kind: "edt"},
		{name: "CFG", kind: "cfg"},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			var root string
			switch c.kind {
			case "edt":
				root = filepath.Join(repoRoot, "input", "edt")
			case "cfg":
				root = filepath.Join(repoRoot, "input", "cfg")
			default:
				t.Fatalf("unknown kind: %s", c.kind)
			}

			var (
				enums []model.MetadataObject
				perr  error
			)

			if c.kind == "edt" {
				p, err := parser.NewEDTParser(root)
				if err != nil {
					t.Fatalf("failed to create EDT parser: %v", err)
				}
				enums, perr = p.ParseEnums()
			} else {
				p, err := parser.NewCFGParser(root)
				if err != nil {
					t.Fatalf("failed to create CFG parser: %v", err)
				}
				enums, perr = p.ParseEnums()
			}

			if perr != nil {
				t.Fatalf("ParseEnums %s error: %v", c.kind, perr)
			}
			if len(enums) == 0 {
				t.Fatalf("no enums parsed from %s fixtures", c.kind)
			}

			var targetEnum model.MetadataObject
			for _, e := range enums {
				if e.Name == "СостоянияЗаказов" {
					targetEnum = e
					break
				}
			}
			if targetEnum.Name == "" {
				t.Fatalf("target enum `СостоянияЗаказов` not found in %s fixtures", c.kind)
			}

			g := NewMarkdownGenerator("")
			got := testutil.Normalize(g.generateContent(targetEnum))

			if strings.TrimSpace(got) != strings.TrimSpace(ref) {
				t.Fatalf("generated markdown does not match reference\n--- got ---\n%s\n--- ref ---\n%s", got, ref)
			}
		})
	}
}
