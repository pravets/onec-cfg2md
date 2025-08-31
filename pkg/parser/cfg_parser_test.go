package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"onec-cfg2md/pkg/generator"
	"onec-cfg2md/pkg/model"
	"onec-cfg2md/pkg/testutil"
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

func TestParseConstants_FromFixtures(t *testing.T) {
	fixtures := filepath.Join("..", "..", "fixtures", "input", "cfg")
	p, err := NewCFGParser(fixtures)
	if err != nil {
		t.Fatalf("NewCFGParser: %v", err)
	}

	consts, err := p.ParseConstants()
	if err != nil {
		t.Fatalf("ParseConstants: %v", err)
	}

	// Expecting fixtures to contain ВалютаУчета and УчетПоСкладам
	c1 := findByName(consts, "ВалютаУчета")
	if c1 == nil {
		t.Fatalf("expected constant ВалютаУчета among %d constants", len(consts))
	}
	if len(c1.Attributes) == 0 {
		t.Fatalf("expected attributes for constant ВалютаУчета, got none")
	}
	if c1.Attributes[0].Name != "Значение" {
		t.Fatalf("expected attribute name 'Значение' for constant, got '%s'", c1.Attributes[0].Name)
	}
	if len(c1.Attributes[0].Types) == 0 {
		t.Fatalf("expected types for constant ВалютаУчета attribute, got none")
	}

	if findByName(consts, "УчетПоСкладам") == nil {
		t.Fatalf("expected constant УчетПоСкладам among %d constants", len(consts))
	}
}

func TestCFG_ParseFilterCriteria_TempFile(t *testing.T) {
	// Use real fixture from fixtures/input/cfg/FilterCriteria
	fixtures := filepath.Join("..", "..", "fixtures", "input", "cfg")
	p, err := NewCFGParser(fixtures)
	if err != nil {
		t.Fatalf("NewCFGParser: %v", err)
	}

	fcs, err := p.ParseFilterCriteria()
	if err != nil {
		t.Fatalf("ParseFilterCriteria: %v", err)
	}
	if len(fcs) == 0 {
		t.Fatalf("expected filter criteria from CFG fixtures, got none")
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

	// Сгенерируем Markdown и сравним с эталоном в fixtures/output
	// Используем MarkdownGenerator как в реальном коде
	outDir := t.TempDir()
	mg := generator.NewMarkdownGenerator(outDir)
	if err := mg.GenerateFiles([]model.MetadataObject{target}); err != nil {
		t.Fatalf("GenerateFiles: %v", err)
	}

	// Открываем сгенерированный файл
	// getFileName неэкспортирован — в тесте формируем имя как в генераторе для FilterCriteria
	genPath := filepath.Join(outDir, fmt.Sprintf("КритерийОтбора_%s.md", target.Name))
	genBytes, err := os.ReadFile(genPath)
	if err != nil {
		t.Fatalf("read generated markdown: %v", err)
	}
	genContent := string(genBytes)

	// Читаем эталонный fixture
	fixturePath := filepath.Join("..", "..", "fixtures", "output", "КритерийОтбора_ДокументыКонтрагента.md")
	expectedBytes, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	expectedContent := string(expectedBytes)

	// Нормализуем содержимое (удаление BOM/CRLF/лишних пустых строк)
	genNorm := testutil.Normalize(genContent)
	expNorm := testutil.Normalize(expectedContent)

	// Проверим заголовок и секции
	if !strings.Contains(genNorm, "# Критерий отбора: ДокументыКонтрагента") && !strings.Contains(genNorm, "# КритерийОтбора: ДокументыКонтрагента") {
		t.Fatalf("generated markdown missing expected header for ДокументыКонтрагента")
	}
	if !strings.Contains(genNorm, "## Типы") {
		t.Fatalf("generated markdown missing 'Типы' section")
	}
	if !strings.Contains(genNorm, "## Состав") {
		t.Fatalf("generated markdown missing 'Состав' section")
	}

	// Проверим содержимое секции "Типы": извлечём типы из эталона и сравним с target.FilterCriteriaTypes
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

	// Сравним ожидаемые типы с разобранными типами в объекте
	if len(target.FilterCriteriaTypes) != len(expectedTypes) {
		t.Fatalf("mismatch in number of types: parsed=%d expected=%d", len(target.FilterCriteriaTypes), len(expectedTypes))
	}
	for i, et := range expectedTypes {
		if et != target.FilterCriteriaTypes[i] {
			t.Fatalf("type mismatch at index %d: parsed='%s' expected='%s'", i, target.FilterCriteriaTypes[i], et)
		}
	}

	// Строгое сравнение нормализованного сгенерированного содержимого с эталоном
	if genNorm != expNorm {
		t.Fatalf("generated markdown does not exactly match fixture\n--- generated ---\n%s\n--- expected ---\n%s", genNorm, expNorm)
	}

	// Извлечём строки состава из эталона
	var expectedItems []string
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
					expectedItems = append(expectedItems, strings.TrimSpace(strings.TrimPrefix(ln, "- ")))
				}
			}
			if strings.HasPrefix(ln, "## Состав") {
				in = true
			}
		}
	}

	if len(expectedItems) == 0 {
		t.Fatalf("fixture: expected items in 'Состав' section, found none")
	}

	// Применим простую нормализацию/маппинг терминов в сгенерированном выводе,
	// чтобы сравнить виды имен (Document -> Документ, Attribute -> Реквизит)
	mapPairs := []struct{ from, to string }{
		{"Document.", "Документ."},
		{"Attribute.", "Реквизит."},
	}
	mappedGen := genNorm
	for _, mp := range mapPairs {
		mappedGen = strings.ReplaceAll(mappedGen, mp.from, mp.to)
	}

	// Проверим, что все ожидаемые элементы присутствуют в сгенерированном (после маппинга)
	for _, item := range expectedItems {
		if !strings.Contains(mappedGen, item) {
			t.Fatalf("generated markdown missing expected content line after mapping: %s\n--- generated mapped ---\n%s", item, mappedGen)
		}
	}
}
