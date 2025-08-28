package generator

import (
	"strings"
	"testing"

	"ones-cfg2md/pkg/model"
)

func TestGenerateContent_Enum(t *testing.T) {
	g := NewMarkdownGenerator("./out")
	obj := model.MetadataObject{
		Type:    model.ObjectTypeEnum,
		Name:    "ТестовоеПеречисление",
		Synonym: "Тестовое перечисление",
		EnumValues: []model.EnumValue{
			{Name: "Значение1", Synonym: "Значение 1"},
			{Name: "Значение2", Synonym: "Значение (2)"},
		},
	}

	content := g.generateContent(obj)
	if len(content) == 0 {
		t.Fatalf("generated content is empty")
	}
	if want := "# Перечисление: ТестовоеПеречисление (Тестовое перечисление)"; !strings.Contains(content, want) {
		t.Fatalf("expected header not found: %s", want)
	}
	if !strings.Contains(content, "- Значение1 (Значение 1)") {
		t.Fatalf("expected enum value line not found")
	}
}
