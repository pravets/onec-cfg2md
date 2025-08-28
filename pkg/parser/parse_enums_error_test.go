package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseEnums_MalformedXML(t *testing.T) {
	dir := t.TempDir()
	// create minimal cfg layout with malformed enum XML
	fixturesDir := filepath.Join(dir, "Enums")
	if err := os.MkdirAll(fixturesDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	bad := `<?xml version="1.0" encoding="utf-8"?>\n<Enum><Name>BadEnum` // truncated
	if err := os.WriteFile(filepath.Join(fixturesDir, "BadEnum.xml"), []byte(bad), 0644); err != nil {
		t.Fatalf("write bad xml: %v", err)
	}

	p, err := NewCFGParser(dir)
	if err != nil {
		t.Fatalf("failed to create CFG parser: %v", err)
	}

	enums, err := p.ParseEnums()
	if err == nil && len(enums) > 0 {
		t.Fatalf("expected error or zero enums for malformed input, got %d enums", len(enums))
	}
}
