package generator

import (
	"ones-cfg2md/pkg/model"
	"ones-cfg2md/pkg/parser"
	"ones-cfg2md/pkg/testutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateContent_Enum(t *testing.T) {
	// Path to fixtures, relative to the current test file
	fixtureDir := filepath.Join("..", "..", "fixtures")

	// 1. Parse the enum from the CFG fixture
	cfgRoot := filepath.Join(fixtureDir, "input", "cfg")
	p, err := parser.NewCFGParser(cfgRoot)
	if err != nil {
		t.Fatalf("failed to create CFG parser: %v", err)
	}

	enums, err := p.ParseEnums()
	if err != nil {
		t.Fatalf("ParseEnums CFG error: %v", err)
	}

	var targetEnum model.MetadataObject
	for _, e := range enums {
		if e.Name == "СостоянияЗаказов" {
			targetEnum = e
			break
		}
	}
	if targetEnum.Name == "" {
		t.Fatal("target enum `СостоянияЗаказов` not found in cfg fixtures")
	}

	// 2. Generate the markdown content
	g := NewMarkdownGenerator("") // output dir doesn't matter for content generation
	got := testutil.Normalize(g.generateContent(targetEnum))

	// 3. Load the reference ("golden") file
	refPath := filepath.Join(fixtureDir, "output", "Перечисление_СостоянияЗаказов.md")
	refBytes, err := os.ReadFile(refPath)
	if err != nil {
		t.Fatalf("failed to read reference file: %v", err)
	}
	want := testutil.Normalize(string(refBytes))

	// 4. Compare generated content with the reference
	if got != want {
		t.Fatalf("generated markdown does not match reference\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}
