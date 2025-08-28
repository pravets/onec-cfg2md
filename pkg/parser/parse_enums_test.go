package parser

import (
	"ones-cfg2md/pkg/model"
	"path/filepath"
	"runtime"
	"testing"
)

func TestParseEnums_FromFixture(t *testing.T) {
	// run subtests for EDT and CFG
	t.Run("EDT", func(t *testing.T) {
		// compute path relative to this test file to avoid depending on working directory
		_, thisFile, _, _ := runtime.Caller(0)
		testRoot := filepath.Join(filepath.Dir(thisFile), "testdata", "input", "edt")
		edtFixture := testRoot

		p, err := NewEDTParser(edtFixture)
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
		e := enums[0]
		if e.Type != model.ObjectTypeEnum {
			t.Fatalf("expected ObjectTypeEnum, got %s", e.Type)
		}
	})

	t.Run("CFG", func(t *testing.T) {
		// compute path relative to this test file to avoid depending on working directory
		_, thisFile, _, _ := runtime.Caller(0)
		testRoot := filepath.Join(filepath.Dir(thisFile), "testdata", "input", "cfg")
		cfgFixture := testRoot

		p, err := NewCFGParser(cfgFixture)
		if err != nil {
			t.Fatalf("failed to create CFG parser: %v", err)
		}

		enums, err := p.ParseEnums()
		if err != nil {
			t.Fatalf("ParseEnums CFG error: %v", err)
		}
		if len(enums) == 0 {
			t.Fatalf("no enums parsed from CFG fixtures")
		}
		e := enums[0]
		if e.Type != model.ObjectTypeEnum {
			t.Fatalf("expected ObjectTypeEnum, got %s", e.Type)
		}
	})
}
