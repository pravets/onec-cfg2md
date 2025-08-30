package model

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestConstants(t *testing.T) {
	if ObjectTypeDocument != "Document" {
		t.Fatalf("ObjectTypeDocument value changed: %q", ObjectTypeDocument)
	}
	if ObjectTypeCatalog != "Catalog" {
		t.Fatalf("ObjectTypeCatalog value changed: %q", ObjectTypeCatalog)
	}
	if ObjectTypeEnum != "Enum" {
		t.Fatalf("ObjectTypeEnum value changed: %q", ObjectTypeEnum)
	}
	if FormatCFG != "cfg" {
		t.Fatalf("FormatCFG value changed: %q", FormatCFG)
	}
	if FormatEDT != "edt" {
		t.Fatalf("FormatEDT value changed: %q", FormatEDT)
	}
}

func TestMetadataObjectJSONRoundtrip(t *testing.T) {
	src := MetadataObject{
		Type:    ObjectTypeCatalog,
		Name:    "Контрагенты",
		Synonym: "Контрагенты син",
		Attributes: []Attribute{
			{Name: "Code", Synonym: "Код", Types: []string{"String"}, Required: true},
		},
		TabularSections: []TabularSection{
			{Name: "Items", Synonym: "Товары", Attributes: []Attribute{{Name: "Qty"}}},
		},
		Dimensions: []Attribute{{Name: "Dim1"}},
		Resources:  []Attribute{{Name: "Res1"}},
		EnumValues: []EnumValue{{Name: "Val1", Synonym: "Знач1"}},
	}

	b, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var dst MetadataObject
	if err := json.Unmarshal(b, &dst); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if !reflect.DeepEqual(src, dst) {
		t.Fatalf("roundtrip mismatch\nsrc=%#v\ndst=%#v", src, dst)
	}
}

func TestConversionOptionsJSONRoundtrip(t *testing.T) {
	src := ConversionOptions{
		SourcePath:  "in/path",
		OutputPath:  "out/path",
		Format:      FormatCFG,
		ObjectTypes: []ObjectType{ObjectTypeDocument, ObjectTypeCatalog},
		Verbose:     true,
	}

	b, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var dst ConversionOptions
	if err := json.Unmarshal(b, &dst); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if !reflect.DeepEqual(src, dst) {
		t.Fatalf("roundtrip mismatch\nsrc=%#v\ndst=%#v", src, dst)
	}
}

func TestCatalogEntryJSONRoundtrip(t *testing.T) {
	src := CatalogEntry{
		ObjectName: "Заказ",
		ObjectType: string(ObjectTypeDocument),
		Synonym:    "Заказ (син)",
		FileName:   "Документ_Заказ.md",
	}

	b, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var dst CatalogEntry
	if err := json.Unmarshal(b, &dst); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if !reflect.DeepEqual(src, dst) {
		t.Fatalf("roundtrip mismatch\nsrc=%#v\ndst=%#v", src, dst)
	}
}
