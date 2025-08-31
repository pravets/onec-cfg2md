package parser

import (
	"onec-cfg2md/pkg/model"
	"path/filepath"
	"runtime"
	"testing"
)

func TestParseEnums_FromFixture(t *testing.T) {
	verify := func(t *testing.T, enums []model.MetadataObject) {
		if len(enums) == 0 {
			t.Fatalf("no enums parsed from fixtures")
		}

		var targetEnum model.MetadataObject
		for _, e := range enums {
			if e.Name == "СостоянияЗаказов" {
				targetEnum = e
				break
			}
		}

		if targetEnum.Name == "" {
			t.Fatalf("target enum `СостоянияЗаказов` not found in fixtures")
		}

		if targetEnum.Type != model.ObjectTypeEnum {
			t.Errorf("expected Type to be '%s', but got '%s'", model.ObjectTypeEnum, targetEnum.Type)
		}
		if targetEnum.Synonym != "Состояния заказов" {
			t.Errorf("expected Synonym to be 'Состояния заказов', but got '%s'", targetEnum.Synonym)
		}
		if len(targetEnum.EnumValues) != 4 {
			t.Fatalf("expected 4 enum values, but got %d", len(targetEnum.EnumValues))
		}

		expectedValues := map[string]string{
			"Открыт":   "Открыт",
			"ВРаботе":  "В работе",
			"Выполнен": "Выполнен",
			"Закрыт":   "Закрыт",
		}

		for _, val := range targetEnum.EnumValues {
			expectedSynonym, ok := expectedValues[val.Name]
			if !ok {
				t.Errorf("unexpected enum value name: %s", val.Name)
				continue
			}
			if val.Synonym != expectedSynonym {
				t.Errorf("for enum value %s, expected synonym '%s', but got '%s'", val.Name, expectedSynonym, val.Synonym)
			}
		}
	}

	t.Run("EDT", func(t *testing.T) {
		_, thisFile, _, _ := runtime.Caller(0)
		testRoot := filepath.Join(filepath.Dir(thisFile), "..", "..", "fixtures", "input", "edt")

		p, err := NewEDTParser(testRoot)
		if err != nil {
			t.Fatalf("failed to create EDT parser: %v", err)
		}

		enums, err := p.ParseEnums()
		if err != nil {
			t.Fatalf("ParseEnums EDT error: %v", err)
		}
		verify(t, enums)
	})

	t.Run("CFG", func(t *testing.T) {
		_, thisFile, _, _ := runtime.Caller(0)
		testRoot := filepath.Join(filepath.Dir(thisFile), "..", "..", "fixtures", "input", "cfg")

		p, err := NewCFGParser(testRoot)
		if err != nil {
			t.Fatalf("failed to create CFG parser: %v", err)
		}

		enums, err := p.ParseEnums()
		if err != nil {
			t.Fatalf("ParseEnums CFG error: %v", err)
		}
		verify(t, enums)
	})
}
