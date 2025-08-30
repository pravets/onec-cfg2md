package generator

import (
	"fmt"
	"ones-cfg2md/pkg/model"
	"ones-cfg2md/pkg/parser"
	"ones-cfg2md/pkg/testutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateContent_GoldenFiles(t *testing.T) {
	fixtureDir := filepath.Join("..", "..", "fixtures")
	cfgRoot := filepath.Join(fixtureDir, "input", "cfg")

	// 1. Parse all objects from the CFG fixture at once
	p, err := parser.NewCFGParser(cfgRoot)
	if err != nil {
		t.Fatalf("failed to create CFG parser: %v", err)
	}
	// Parse all supported types
	allObjectTypes := []model.ObjectType{
		model.ObjectTypeDocument,
		model.ObjectTypeCatalog,
		model.ObjectTypeAccumulationRegister,
		model.ObjectTypeInformationRegister,
		model.ObjectTypeEnum,
		model.ObjectTypeChartOfCharacteristicTypes,
	}
	parsedObjects, err := p.ParseObjectsByType(allObjectTypes)
	if err != nil {
		t.Fatalf("ParseObjectsByType CFG error: %v", err)
	}

	// Create a map for easy lookup
	objectsMap := make(map[string]model.MetadataObject)
	for _, obj := range parsedObjects {
		objectsMap[obj.Name] = obj
	}

	testCases := []struct {
		name           string
		objectName     string
		goldenFileName string
	}{
		{"Enum", "СостоянияЗаказов", "Перечисление_СостоянияЗаказов.md"},
		{"Document", "Заказ", "Документ_Заказ.md"},
		{"Catalog", "Контрагенты", "Справочник_Контрагенты.md"},
		{"AccumulationRegister", "Взаиморасчеты", "РегистрНакопления_Взаиморасчеты.md"},
		{"InformationRegister", "КурсыВалют", "РегистрСведений_КурсыВалют.md"},
		{"ChartOfCharacteristicTypes", "ВидыХарактеристик", "ПланВидовХарактеристик_ВидыХарактеристик.md"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			targetObject, ok := objectsMap[tc.objectName]
			if !ok {
				t.Fatalf("target object `%s` not found in parsed fixtures", tc.objectName)
			}

			// 2. Generate the markdown content
			g := NewMarkdownGenerator("") // output dir doesn't matter for content generation
			got := testutil.Normalize(g.generateContent(targetObject))

			// 3. Load the reference ("golden") reference file
			refPath := filepath.Join(fixtureDir, "output", tc.goldenFileName)
			refBytes, err := os.ReadFile(refPath)
			if err != nil {
				t.Fatalf("failed to read reference file '%s': %v", refPath, err)
			}
			want := testutil.Normalize(string(refBytes))

			// 4. Compare generated content with the reference
			if got != want {
				t.Fatalf("generated markdown does not match reference\n---\ngot---\n%s\n--- want ---\n%s", got, want)
			}
		})
	}
}

func TestGenerateFiles(t *testing.T) {
	// 1. Prepare mock data
	objects := []model.MetadataObject{
		{Type: model.ObjectTypeDocument, Name: "РасходнаяНакладная"},
		{Type: model.ObjectTypeCatalog, Name: "Номенклатура"},
	}

	// 2. Create a temporary output directory
	outputDir, err := os.MkdirTemp("", "markdown-generator-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(outputDir)

	// 3. Run the generator
	g := NewMarkdownGenerator(outputDir)
	if err := g.GenerateFiles(objects); err != nil {
		t.Fatalf("GenerateFiles() returned an error: %v", err)
	}

	// 4. Check if files were created
	expectedFiles := []string{
		"Документ_РасходнаяНакладная.md",
		"Справочник_Номенклатура.md",
	}
	for _, fName := range expectedFiles {
		fPath := filepath.Join(outputDir, fName)
		if _, err := os.Stat(fPath); os.IsNotExist(err) {
			t.Errorf("Expected file was not created: %s", fPath)
		}
	}

	// 5. Check content of one file
	content, err := os.ReadFile(filepath.Join(outputDir, expectedFiles[0]))
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}
	expectedContent := "# Документ: РасходнаяНакладная\n\n"
	if testutil.Normalize(string(content)) != testutil.Normalize(expectedContent) {
		t.Errorf("File content mismatch. Got '%s', want '%s'", string(content), expectedContent)
	}
}

func TestGenerateFiles_EmptyInput(t *testing.T) {
	outputDir, err := os.MkdirTemp("", "markdown-empty-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(outputDir)

	g := NewMarkdownGenerator(outputDir)
	if err := g.GenerateFiles([]model.MetadataObject{}); err != nil {
		t.Fatalf("GenerateFiles() with empty slice returned an error: %v", err)
	}

	// Check that dir is empty
	files, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatalf("Could not read output dir: %v", err)
	}
	if len(files) > 0 {
		t.Errorf("Output directory should be empty, but contains %d files", len(files))
	}
}

func TestGenerateFiles_MkdirError(t *testing.T) {
	// Create a temporary file to cause MkdirAll to fail
	tmpFile, err := os.CreateTemp("", "test-file-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	g := NewMarkdownGenerator(tmpFile.Name())
	err = g.GenerateFiles([]model.MetadataObject{{Type: "test", Name: "test"}})
	if err == nil {
		t.Fatalf("Expected an error when output path is a file, but got nil")
	}
}

func TestGenerateFiles_WriteFileError(t *testing.T) {
	// Create a temporary directory
	outputDir, err := os.MkdirTemp("", "markdown-readonly-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(outputDir)

	// Make the directory read-only
	if err := os.Chmod(outputDir, 0555); err != nil {
		t.Fatalf("Failed to change directory permissions: %v", err)
	}

	g := NewMarkdownGenerator(outputDir)
	objects := []model.MetadataObject{
		{Type: model.ObjectTypeDocument, Name: "SomeDoc"},
	}
	err = g.GenerateFiles(objects)

	// Check if we got an error as expected
	if err == nil {
		// If there was no error, check if the file was unexpectedly created.
		// This can happen if the test is run as root or on a filesystem
		// that doesn't enforce read-only permissions strictly.
		filePath := filepath.Join(outputDir, g.getFileName(objects[0]))
		if _, statErr := os.Stat(filePath); statErr == nil {
			t.Logf("Warning: TestGenerateFiles_WriteFileError is unreliable on this system. A file was created in a directory that should be read-only.")
		} else {
			t.Fatalf("Expected a file write error, but got nil and the file was not created.")
		}
	}
}

func TestGetObjectTypeRussian(t *testing.T) {
	testCases := []struct {
		objType model.ObjectType
		want    string
	}{
		{model.ObjectTypeDocument, "Документ"},
		{model.ObjectTypeCatalog, "Справочник"},
		{model.ObjectTypeAccumulationRegister, "РегистрНакопления"},
		{model.ObjectTypeInformationRegister, "РегистрСведений"},
		{model.ObjectTypeEnum, "Перечисление"},
		{model.ObjectTypeChartOfCharacteristicTypes, "ПланВидовХарактеристик"},
		{"UnknownType", "UnknownType"},
	}

	g := NewMarkdownGenerator("")
	for _, tc := range testCases {
		t.Run(string(tc.objType), func(t *testing.T) {
			got := g.getObjectTypeRussian(tc.objType)
			if got != tc.want {
				t.Errorf("getObjectTypeRussian(%s) got %s, want %s", tc.objType, got, tc.want)
			}
		})
	}
}

func TestGetFileName(t *testing.T) {
	obj := model.MetadataObject{
		Type: model.ObjectTypeDocument,
		Name: "MyTestDocument",
	}
	g := NewMarkdownGenerator("")
	want := "Документ_MyTestDocument.md"
	got := g.getFileName(obj)
	if got != want {
		t.Errorf("getFileName() got %s, want %s", got, want)
	}
}

func TestGenerateContent_EdgeCases(t *testing.T) {
	g := NewMarkdownGenerator("")

	testCases := []struct {
		name string
		obj  model.MetadataObject
		want string
	}{
		{
			name: "Enum with no values",
			obj: model.MetadataObject{
				Type:    model.ObjectTypeEnum,
				Name:    "EmptyEnum",
				Synonym: "ПустоеПеречисление",
			},
			want: "# Перечисление: EmptyEnum (ПустоеПеречисление)\n\n",
		},
		{
			name: "Document with no attributes or tabular sections",
			obj: model.MetadataObject{
				Type: model.ObjectTypeDocument,
				Name: "EmptyDoc",
			},
			want: "# Документ: EmptyDoc\n\n",
		},
		{
			name: "Catalog with no attributes",
			obj: model.MetadataObject{
				Type: model.ObjectTypeCatalog,
				Name: "EmptyCatalog",
			},
			want: "# Справочник: EmptyCatalog\n\n",
		},
		{
			name: "Accumulation Register with only one field type",
			obj: model.MetadataObject{
				Type: model.ObjectTypeAccumulationRegister,
				Name: "PartialAccumReg",
				Dimensions: []model.Attribute{
					{Name: "Dim1", Types: []string{"String"}},
				},
			},
			want: "# РегистрНакопления: PartialAccumReg\n\n## Измерения\n\n- Dim1 (String)\n\n",
		},
		{
			name: "Information Register with only one field type",
			obj: model.MetadataObject{
				Type: model.ObjectTypeInformationRegister,
				Name: "PartialInfoReg",
				Resources: []model.Attribute{
					{Name: "Res1", Types: []string{"Number"}},
				},
			},
			want: "# РегистрСведений: PartialInfoReg\n\n## Ресурсы\n\n- Res1 (Number)\n\n",
		},
		{
			name: "Document with tabular section without synonym",
			obj: model.MetadataObject{
				Type: model.ObjectTypeDocument,
				Name: "DocWithTabSection",
				TabularSections: []model.TabularSection{
					{
						Name: "Товары",
						Attributes: []model.Attribute{
							{Name: "Номенклатура", Types: []string{"Catalog.Номенклатура"}},
						},
					},
				},
			},
			want: "# Документ: DocWithTabSection\n\n## Табличные части\n\n### Товары\n\n- Номенклатура (Catalog.Номенклатура)\n\n",
		},
		{
			name: "Enum value without synonym",
			obj: model.MetadataObject{
				Type: model.ObjectTypeEnum,
				Name: "EnumWithNoSynonymValue",
				EnumValues: []model.EnumValue{
					{Name: "Value1"},
				},
			},
			want: "# Перечисление: EnumWithNoSynonymValue\n\n## Значения\n\n- Value1\n\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := g.generateContent(tc.obj)
			if testutil.Normalize(got) != testutil.Normalize(tc.want) {
				t.Errorf("generateContent() for '%s' failed\n---\ngot---\n%s\n--- want ---\n%s", tc.name, got, tc.want)
			}
		})
	}
}

func TestGenerateContent_Constant(t *testing.T) {
	g := NewMarkdownGenerator("")
	obj := model.MetadataObject{
		Type:       model.ObjectTypeConstant,
		Name:       "ВалютаУчета",
		Synonym:    "Валюта учета",
		Attributes: []model.Attribute{{Name: "Значение", Types: []string{"CatalogRef.Валюты"}}},
	}
	got := g.generateContent(obj)
	if !strings.Contains(got, "# Константа: ВалютаУчета") {
		t.Fatalf("expected header for constant, got: %s", got)
	}
	if !strings.Contains(got, "Значение") {
		t.Fatalf("expected attribute name 'Значение' in content, got: %s", got)
	}
}

func TestGenerateContent_Constants_GoldenFiles(t *testing.T) {
	fixtureDir := filepath.Join("..", "..", "fixtures")
	cfgRoot := filepath.Join(fixtureDir, "input", "cfg")

	p, err := parser.NewCFGParser(cfgRoot)
	if err != nil {
		t.Fatalf("failed to create CFG parser: %v", err)
	}

	consts, err := p.ParseConstants()
	if err != nil {
		t.Fatalf("ParseConstants: %v", err)
	}

	if len(consts) == 0 {
		t.Fatalf("expected constants in fixtures, got none")
	}

	g := NewMarkdownGenerator("")

	for _, c := range consts {
		t.Run(c.Name, func(t *testing.T) {
			got := testutil.Normalize(g.generateContent(c))

			golden := filepath.Join(fixtureDir, "output", fmt.Sprintf("Константа_%s.md", c.Name))
			refBytes, err := os.ReadFile(golden)
			if err != nil {
				t.Fatalf("failed to read golden file '%s': %v", golden, err)
			}
			want := testutil.Normalize(string(refBytes))

			if got != want {
				t.Fatalf("constant markdown mismatch for %s\n---\ngot---\n%s\n--- want ---\n%s", c.Name, got, want)
			}
		})
	}
}
