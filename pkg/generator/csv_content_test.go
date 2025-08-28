package generator

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ones-cfg2md/pkg/model"
)

func TestGenerateCatalog_Escaping(t *testing.T) {
	out := t.TempDir()
	g := NewCSVGenerator(out)

	objs := []model.MetadataObject{
		{
			Type:    model.ObjectTypeCatalog,
			Name:    "Контрагенты,Имя",
			Synonym: "Синоним; с точкой с запятой",
		},
	}

	if err := g.GenerateCatalog(objs); err != nil {
		t.Fatalf("GenerateCatalog failed: %v", err)
	}

	csvPath := filepath.Join(out, "objects.csv")
	f, err := os.Open(csvPath)
	if err != nil {
		t.Fatalf("failed to open generated csv: %v", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("failed to read csv: %v", err)
	}

	if len(records) < 2 {
		t.Fatalf("expected header + 1 record, got %d rows", len(records))
	}

	// verify header exactly
	header := records[0]
	expectedHeader := []string{"Имя объекта", "Тип объекта", "Синоним", "Файл"}
	if len(header) != len(expectedHeader) {
		t.Fatalf("unexpected header length: got %d want %d", len(header), len(expectedHeader))
	}
	for i := range header {
		if header[i] != expectedHeader[i] {
			t.Fatalf("unexpected header column %d: got %q want %q", i, header[i], expectedHeader[i])
		}
	}

	rec := records[1]
	if len(rec) != len(expectedHeader) {
		t.Fatalf("unexpected record column count: got %d want %d", len(rec), len(expectedHeader))
	}
	// ObjectName should be like "Справочник.Контрагенты,Имя"
	if !strings.Contains(rec[0], "Справочник.Контрагенты,Имя") {
		t.Fatalf("unexpected object name in csv: %s", rec[0])
	}
	if !strings.Contains(rec[2], "Синоним; с точкой с запятой") {
		t.Fatalf("unexpected synonym in csv: %s", rec[2])
	}
}
