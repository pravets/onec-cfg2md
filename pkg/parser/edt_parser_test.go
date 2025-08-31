package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ones-cfg2md/pkg/generator"
	"ones-cfg2md/pkg/model"
	"ones-cfg2md/pkg/testutil"
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

func TestEDT_ParseCharts_FromFixtures(t *testing.T) {
	fixtures := filepath.Join("..", "..", "fixtures", "input", "edt")
	p, err := NewEDTParser(fixtures)
	if err != nil {
		t.Fatalf("NewEDTParser: %v", err)
	}

	charts, err := p.ParseChartsOfCharacteristicTypes()
	if err != nil {
		t.Fatalf("ParseChartsOfCharacteristicTypes: %v", err)
	}
	if len(charts) == 0 {
		t.Fatalf("expected charts from EDT fixtures, got none")
	}
	if findByName(charts, "ВидыХарактеристик") == nil {
		t.Fatalf("expected chart ВидыХарактеристик in EDT fixtures")
	}
}

func TestEDT_ParseConstants_FromFixtures(t *testing.T) {
	fixtures := filepath.Join("..", "..", "fixtures", "input", "edt")
	p, err := NewEDTParser(fixtures)
	if err != nil {
		t.Fatalf("NewEDTParser: %v", err)
	}

	consts, err := p.ParseConstants()
	if err != nil {
		t.Fatalf("ParseConstants: %v", err)
	}

	if len(consts) == 0 {
		t.Fatalf("expected constants from EDT fixtures, got none")
	}

	// Check at least ВалютаУчета exists
	if findByName(consts, "ВалютаУчета") == nil {
		t.Fatalf("expected constant ВалютаУчета in EDT fixtures")
	}
}

func TestEDT_ParseFilterCriteria_FromFixtures(t *testing.T) {
	fixtures := filepath.Join("..", "..", "fixtures")
	p, err := NewEDTParser(filepath.Join(fixtures, "input", "edt"))
	if err != nil {
		t.Fatalf("NewEDTParser: %v", err)
	}

	fcs, err := p.ParseFilterCriteria()
	if err != nil {
		t.Fatalf("ParseFilterCriteria: %v", err)
	}
	if len(fcs) == 0 {
		t.Fatalf("expected filter criteria from EDT fixtures, got none")
	}

	// Найдём критерий ДокументыКонтрагента
	var target model.MetadataObject
	for _, f := range fcs {
		if f.Name == "ДокументыКонтрагента" {
			target = f
			break
		}
	}
	if target.Name == "" {
		t.Fatalf("expected criterion ДокументыКонтрагента in fixtures")
	}

	if target.Synonym != "Документы контрагента" {
		t.Fatalf("unexpected synonym: %s", target.Synonym)
	}

	// Проверим типы
	if len(target.FilterCriteriaTypes) == 0 {
		t.Fatalf("expected FilterCriteriaTypes to be populated, got none")
	}

	// Проверим состав (content)
	if len(target.FilterCriteriaContents) == 0 {
		t.Fatalf("expected FilterCriteriaContents to be populated, got none")
	}
}

func TestEDT_ParseFilterCriteria_TempFile(t *testing.T) {
	fixtures := filepath.Join("..", "..", "fixtures", "input", "edt")
	p, err := NewEDTParser(fixtures)
	if err != nil {
		t.Fatalf("NewEDTParser: %v", err)
	}

	fcs, err := p.ParseFilterCriteria()
	if err != nil {
		t.Fatalf("ParseFilterCriteria: %v", err)
	}
	if len(fcs) == 0 {
		t.Fatalf("expected filter criteria from EDT fixtures, got none")
	}

	// Найдём критерий ДокументыКонтрагента
	var target model.MetadataObject
	for _, f := range fcs {
		if f.Name == "ДокументыКонтрагента" {
			target = f
			break
		}
	}
	if target.Name == "" {
		t.Fatalf("expected criterion ДокументыКонтрагента in fixtures")
	}

	if target.Synonym != "Документы контрагента" {
		t.Fatalf("unexpected synonym: %s", target.Synonym)
	}

	// Проверим типы
	if len(target.FilterCriteriaTypes) == 0 {
		t.Fatalf("expected FilterCriteriaTypes to be populated, got none")
	}

	// Проверим состав (content)
	if len(target.FilterCriteriaContents) == 0 {
		t.Fatalf("expected FilterCriteriaContents to be populated, got none")
	}

	// Сгенерируем Markdown и сравним с эталоном
	outDir := t.TempDir()
	mg := generator.NewMarkdownGenerator(outDir)
	if err := mg.GenerateFiles([]model.MetadataObject{target}); err != nil {
		t.Fatalf("GenerateFiles: %v", err)
	}
	genPath := filepath.Join(outDir, fmt.Sprintf("КритерийОтбора_%s.md", target.Name))
	genBytes, err := os.ReadFile(genPath)
	if err != nil {
		t.Fatalf("read generated markdown: %v", err)
	}
	genContent := string(genBytes)

	fixturePath := filepath.Join("..", "..", "fixtures", "output", "КритерийОтбора_ДокументыКонтрагента.md")
	expectedBytes, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	expectedContent := string(expectedBytes)

	genNorm := testutil.Normalize(genContent)
	expNorm := testutil.Normalize(expectedContent)

	if !strings.Contains(genNorm, "# Критерий отбора: ДокументыКонтрагента") && !strings.Contains(genNorm, "# КритерийОтбора: ДокументыКонтрагента") {
		t.Fatalf("generated markdown missing expected header for ДокументыКонтрагента")
	}
	if !strings.Contains(genNorm, "## Типы") {
		t.Fatalf("generated markdown missing 'Типы' section")
	}
	if !strings.Contains(genNorm, "## Состав") {
		t.Fatalf("generated markdown missing 'Состав' section")
	}

	// Проверим типы из fixture против разобранных
	var expectedTypes []string
	{
		lines := strings.Split(expNorm, "\n")
		in := false
		for _, ln := range lines {
			if strings.HasPrefix(ln, "## ") && in {
				break
			}
			if in {
				ln = strings.TrimSpace(ln)
				if strings.HasPrefix(ln, "- ") {
					expectedTypes = append(expectedTypes, strings.TrimSpace(strings.TrimPrefix(ln, "- ")))
				}
			}
			if strings.HasPrefix(ln, "## Типы") {
				in = true
			}
		}
	}
	if len(expectedTypes) == 0 {
		t.Fatalf("fixture: expected items in 'Типы' section, found none")
	}
	if len(target.FilterCriteriaTypes) != len(expectedTypes) {
		t.Fatalf("mismatch in number of types: parsed=%d expected=%d", len(target.FilterCriteriaTypes), len(expectedTypes))
	}
	for i, et := range expectedTypes {
		if et != target.FilterCriteriaTypes[i] {
			t.Fatalf("type mismatch at index %d: parsed='%s' expected='%s'", i, target.FilterCriteriaTypes[i], et)
		}
	}

	if genNorm != expNorm {
		t.Fatalf("generated markdown does not exactly match fixture\n--- generated ---\n%s\n--- expected ---\n%s", genNorm, expNorm)
	}
}
