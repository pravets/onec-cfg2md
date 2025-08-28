package detector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectFormat_CFG(t *testing.T) {
	dir := t.TempDir()
	// create Configuration.xml to simulate CFG
	f, err := os.Create(filepath.Join(dir, "Configuration.xml"))
	if err != nil {
		t.Fatalf("failed create fixture: %v", err)
	}
	f.Close()

	format, err := DetectFormat(dir)
	if err != nil {
		t.Fatalf("DetectFormat returned error: %v", err)
	}
	if format != "cfg" {
		t.Fatalf("expected cfg format, got %s", format)
	}
}

func TestDetectFormat_EDT(t *testing.T) {
	dir := t.TempDir()
	// create .project and src to simulate EDT
	f, err := os.Create(filepath.Join(dir, ".project"))
	if err != nil {
		t.Fatalf("failed create fixture: %v", err)
	}
	f.Close()
	if err := os.Mkdir(filepath.Join(dir, "src"), 0755); err != nil {
		t.Fatalf("failed create src: %v", err)
	}

	format, err := DetectFormat(dir)
	if err != nil {
		t.Fatalf("DetectFormat returned error: %v", err)
	}
	if format != "edt" {
		t.Fatalf("expected edt format, got %s", format)
	}
}

func TestValidateFormat_Errors(t *testing.T) {
	dir := t.TempDir()
	if err := ValidateFormat(dir, "cfg"); err == nil {
		t.Fatalf("expected error when validating cfg without Configuration.xml")
	}
	if err := ValidateFormat(dir, "edt"); err == nil {
		t.Fatalf("expected error when validating edt without .project and src")
	}
}
