package parser

import (
	"path/filepath"
	"testing"

	"ones-cfg2md/pkg/model"
)

func TestNewParser_ExplicitCFG(t *testing.T) {
	p, err := NewParser(".", model.FormatCFG)
	if err != nil {
		t.Fatalf("NewParser CFG returned error: %v", err)
	}
	if _, ok := p.(*CFGParser); !ok {
		t.Fatalf("expected *CFGParser, got %T", p)
	}
}

func TestNewParser_ExplicitEDT(t *testing.T) {
	p, err := NewParser(".", model.FormatEDT)
	if err != nil {
		t.Fatalf("NewParser EDT returned error: %v", err)
	}
	if _, ok := p.(*EDTParser); !ok {
		t.Fatalf("expected *EDTParser, got %T", p)
	}
}

func TestNewParser_AutoDetect_FromFixtures(t *testing.T) {
	cfgFixtures := filepath.Join("..", "..", "fixtures", "input", "cfg")
	p, err := NewParser(cfgFixtures, "")
	if err != nil {
		t.Fatalf("NewParser autodetect cfg returned error: %v", err)
	}
	if _, ok := p.(*CFGParser); !ok {
		t.Fatalf("expected *CFGParser for cfg fixtures, got %T", p)
	}

	edtFixtures := filepath.Join("..", "..", "fixtures", "input", "edt")
	p2, err := NewParser(edtFixtures, "")
	if err != nil {
		t.Fatalf("NewParser autodetect edt returned error: %v", err)
	}
	if _, ok := p2.(*EDTParser); !ok {
		t.Fatalf("expected *EDTParser for edt fixtures, got %T", p2)
	}
}
