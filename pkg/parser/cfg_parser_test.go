package parser

import (
	"path/filepath"
	"testing"

	"ones-cfg2md/pkg/model"
)

func findByName(objs []model.MetadataObject, name string) *model.MetadataObject {
	for _, o := range objs {
		if o.Name == name {
			return &o
		}
	}
	return nil
}

func TestParseDocumentFile_ParseAndTypes_UsingFixtures(t *testing.T) {
	fixtures := filepath.Join("..", "..", "fixtures", "input", "cfg")

	p, err := NewCFGParser(fixtures)
	if err != nil {
		t.Fatalf("NewCFGParser: %v", err)
	}

	docs, err := p.ParseDocuments()
	if err != nil {
		t.Fatalf("ParseDocuments: %v", err)
	}

	d := findByName(docs, "Заказ")
	if d == nil {
		t.Fatalf("document Заказ not found among %d documents", len(docs))
	}
	if d.Synonym == "" {
		t.Fatalf("expected synonym for Заказ, got empty")
	}

	// Атрибуты и табличные части должны присутствовать (fixture содержит Товары и Дата)
	if len(d.TabularSections) == 0 && len(d.Attributes) == 0 {
		t.Fatalf("expected attributes or tabular sections in Заказ, got none")
	}
	// Проверим, что хотя бы один атрибут или табличная часть имеет ненулевой набор типов (парсинг типов)
	hasTypes := false
	for _, a := range d.Attributes {
		if len(a.Types) > 0 {
			hasTypes = true
			break
		}
	}
	if !hasTypes {
		for _, ts := range d.TabularSections {
			for _, a := range ts.Attributes {
				if len(a.Types) > 0 {
					hasTypes = true
					break
				}
			}
			if hasTypes {
				break
			}
		}
	}
	if !hasTypes {
		t.Fatalf("expected some attribute types parsed for Заказ")
	}
}

func TestParseCatalogsEnumsRegistersCharts_UsingFixtures(t *testing.T) {
	fixtures := filepath.Join("..", "..", "fixtures", "input", "cfg")

	p, err := NewCFGParser(fixtures)
	if err != nil {
		t.Fatalf("NewCFGParser: %v", err)
	}

	cats, err := p.ParseCatalogs()
	if err != nil {
		t.Fatalf("ParseCatalogs: %v", err)
	}
	if findByName(cats, "Контрагенты") == nil {
		t.Fatalf("expected catalog Контрагенты among %d catalogs", len(cats))
	}

	enums, err := p.ParseEnums()
	if err != nil {
		t.Fatalf("ParseEnums: %v", err)
	}
	if findByName(enums, "СостоянияЗаказов") == nil {
		t.Fatalf("expected enum СостоянияЗаказов among %d enums", len(enums))
	}

	regs, err := p.ParseAccumulationRegisters()
	if err != nil {
		t.Fatalf("ParseAccumulationRegisters: %v", err)
	}
	// fixtures contain Взаиморасчеты and Продажи — check at least one
	if len(regs) == 0 {
		t.Fatalf("expected accumulation registers, got none")
	}

	iregs, err := p.ParseInformationRegisters()
	if err != nil {
		t.Fatalf("ParseInformationRegisters: %v", err)
	}
	// fixtures include КурсыВалют and МобильныеОтчеты
	if len(iregs) == 0 {
		t.Fatalf("expected information registers, got none")
	}

	charts, err := p.ParseChartsOfCharacteristicTypes()
	if err != nil {
		t.Fatalf("ParseChartsOfCharacteristicTypes: %v", err)
	}
	if findByName(charts, "ВидыХарактеристик") == nil {
		t.Fatalf("expected chart ВидыХарактеристик among %d charts", len(charts))
	}
}

func TestParseObjectsByType_All_UsingFixtures(t *testing.T) {
	fixtures := filepath.Join("..", "..", "fixtures", "input", "cfg")
	p, err := NewCFGParser(fixtures)
	if err != nil {
		t.Fatalf("NewCFGParser: %v", err)
	}

	types := []model.ObjectType{
		model.ObjectTypeDocument,
		model.ObjectTypeCatalog,
		model.ObjectTypeEnum,
		model.ObjectTypeAccumulationRegister,
		model.ObjectTypeInformationRegister,
		model.ObjectTypeChartOfCharacteristicTypes,
	}
	objs, err := p.ParseObjectsByType(types)
	if err != nil {
		t.Fatalf("ParseObjectsByType: %v", err)
	}
	if len(objs) == 0 {
		t.Fatalf("expected aggregated objects from fixtures, got none")
	}
	if findByName(objs, "Заказ") == nil {
		t.Fatalf("expected document Заказ in aggregated objects")
	}
}

func TestIsDateOnly(t *testing.T) {
	p, err := NewCFGParser(".")
	if err != nil {
		t.Fatalf("NewCFGParser: %v", err)
	}

	// тип с DateQualifiers -> Date
	ct := CFGType{
		DateQualifiers: []CFGDateQualifiers{{DateFractions: "Date"}},
	}
	if !p.isDateOnly(ct) {
		t.Fatalf("expected isDateOnly to return true when DateFractions=Date")
	}

	// тип без DateQualifiers -> false
	ct2 := CFGType{}
	if p.isDateOnly(ct2) {
		t.Fatalf("expected isDateOnly to return false when no DateQualifiers present")
	}
}
