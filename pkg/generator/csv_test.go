package generator

import (
	"path/filepath"
	"testing"

	"ones-cfg2md/pkg/model"
)

func TestCreateCatalogEntry(t *testing.T) {
	g := NewCSVGenerator("./out")
	obj := model.MetadataObject{
		Type:    model.ObjectTypeCatalog,
		Name:    "Контрагенты",
		Synonym: "Контрагенты",
	}

	entry := g.createCatalogEntry(obj)
	if entry.FileName != "Справочник_Контрагенты.md" {
		t.Fatalf("unexpected file name: %s", entry.FileName)
	}
	if entry.ObjectType != "Справочник" {
		t.Fatalf("unexpected object type: %s", entry.ObjectType)
	}
	if entry.ObjectName != filepath.Join("Справочник.Контрагенты") {
		t.Fatalf("unexpected object name: %s", entry.ObjectName)
	}
}
