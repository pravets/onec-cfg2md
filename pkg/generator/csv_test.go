package generator

import (
	"testing"

	"onec-cfg2md/pkg/model"
)

func TestCreateCatalogEntry(t *testing.T) {
	g := NewCSVGenerator("./out")
	cases := []struct {
		name           string
		in             model.MetadataObject
		wantObjectName string
		wantObjectType string
		wantSynonym    string
		wantFileName   string
	}{
		{
			name: "Standard catalog",
			in: model.MetadataObject{
				Type:    model.ObjectTypeCatalog,
				Name:    "Контрагенты",
				Synonym: "Наши Контрагенты",
			},
			wantObjectName: "Справочник.Контрагенты",
			wantObjectType: "Справочник",
			wantSynonym:    "Наши Контрагенты",
			wantFileName:   "Справочник_Контрагенты.md",
		},
		{
			name: "Document with empty synonym",
			in: model.MetadataObject{
				Type: model.ObjectTypeDocument,
				Name: "ЗаказКлиента",
			},
			wantObjectName: "Документ.ЗаказКлиента",
			wantObjectType: "Документ",
			wantSynonym:    "",
			wantFileName:   "Документ_ЗаказКлиента.md",
		},
		{
			name: "Enum with special characters in name",
			in: model.MetadataObject{
				Type:    model.ObjectTypeEnum,
				Name:    "Виды_Цен",
				Synonym: "Виды цен",
			},
			wantObjectName: "Перечисление.Виды_Цен",
			wantObjectType: "Перечисление",
			wantSynonym:    "Виды цен",
			wantFileName:   "Перечисление_Виды_Цен.md",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := g.createCatalogEntry(tc.in)

			wantObjectName := tc.wantObjectType + "." + tc.in.Name

			if got.FileName != tc.wantFileName {
				t.Errorf("FileName: got %q, want %q", got.FileName, tc.wantFileName)
			}
			if got.ObjectType != tc.wantObjectType {
				t.Errorf("ObjectType: got %q, want %q", got.ObjectType, tc.wantObjectType)
			}
			if got.ObjectName != wantObjectName {
				t.Errorf("ObjectName: got %q, want %q", got.ObjectName, wantObjectName)
			}
			if got.Synonym != tc.wantSynonym {
				t.Errorf("Synonym: got %q, want %q", got.Synonym, tc.wantSynonym)
			}
		})
	}
}
