package parser

import (
    "path/filepath"
    "testing"

    "ones-cfg2md/pkg/model"
)

func TestEDT_ParseDocumentsAndCatalogs_FromFixtures(t *testing.T) {
    fixtures := filepath.Join("..", "..", "fixtures", "input", "edt")
    p, err := NewEDTParser(fixtures)
    if err != nil {
        t.Fatalf("NewEDTParser: %v", err)
    }

    docs, err := p.ParseDocuments()
    if err != nil {
        t.Fatalf("ParseDocuments: %v", err)
    }
    if len(docs) == 0 {
        t.Fatalf("expected documents from EDT fixtures, got none")
    }
    if findByName(docs, "Заказ") == nil {
        t.Fatalf("expected document Заказ in EDT fixtures")
    }

    cats, err := p.ParseCatalogs()
    if err != nil {
        t.Fatalf("ParseCatalogs: %v", err)
    }
    if len(cats) == 0 {
        t.Fatalf("expected catalogs from EDT fixtures, got none")
    }
}

func TestEDT_ParseEnumsAndRegisters_FromFixtures(t *testing.T) {
    fixtures := filepath.Join("..", "..", "fixtures", "input", "edt")
    p, err := NewEDTParser(fixtures)
    if err != nil {
        t.Fatalf("NewEDTParser: %v", err)
    }

    enums, err := p.ParseEnums()
    if err != nil {
        t.Fatalf("ParseEnums: %v", err)
    }
    if len(enums) == 0 {
        t.Fatalf("expected enums from EDT fixtures, got none")
    }

    aregs, err := p.ParseAccumulationRegisters()
    if err != nil {
        t.Fatalf("ParseAccumulationRegisters: %v", err)
    }
    if len(aregs) == 0 {
        t.Fatalf("expected accumulation registers from EDT fixtures, got none")
    }

    iregs, err := p.ParseInformationRegisters()
    if err != nil {
        t.Fatalf("ParseInformationRegisters: %v", err)
    }
    if len(iregs) == 0 {
        t.Fatalf("expected information registers from EDT fixtures, got none")
    }
}

func TestEDT_ParseObjectsByType_Aggregation(t *testing.T) {
    fixtures := filepath.Join("..", "..", "fixtures", "input", "edt")
    p, err := NewEDTParser(fixtures)
    if err != nil {
        t.Fatalf("NewEDTParser: %v", err)
    }

    types := []model.ObjectType{
        model.ObjectTypeDocument,
        model.ObjectTypeCatalog,
        model.ObjectTypeEnum,
    }
    objs, err := p.ParseObjectsByType(types)
    if err != nil {
        t.Fatalf("ParseObjectsByType: %v", err)
    }
    if len(objs) == 0 {
        t.Fatalf("expected aggregated objects from EDT fixtures, got none")
    }
}
