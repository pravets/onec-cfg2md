package generator

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"ones-cfg2md/pkg/parser"
	"ones-cfg2md/pkg/testutil"
)

func TestGenerateFilesFromEDTMatchesReference(t *testing.T) {
	// use parser testdata as source for input EDT fixtures
	_, thisFile, _, _ := runtime.Caller(0)
	edtRoot := filepath.Join(filepath.Dir(thisFile), "..", "parser", "testdata", "input", "edt")

	p, err := parser.NewEDTParser(edtRoot)
	if err != nil {
		t.Fatalf("failed to create EDT parser: %v", err)
	}

	enums, err := p.ParseEnums()
	if err != nil {
		t.Fatalf("ParseEnums EDT error: %v", err)
	}
	if len(enums) == 0 {
		t.Fatalf("no enums parsed from EDT fixtures")
	}

	out := t.TempDir()
	g := NewMarkdownGenerator(out)
	if err := g.GenerateFiles(enums); err != nil {
		t.Fatalf("GenerateFiles error: %v", err)
	}

	// expected filename for enum
	fileName := "Перечисление_ТестовоеПеречисление.md"
	genPath := filepath.Join(out, fileName)
	data, err := os.ReadFile(genPath)
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}

	got := testutil.Normalize(string(data))
	ref := loadReference(t)

	if got != ref {
		t.Fatalf("generated file content does not match reference\n--- got ---\n%s\n--- ref ---\n%s", got, ref)
	}
}
