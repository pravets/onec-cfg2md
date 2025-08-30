package model

import (
	"testing"
)

func TestMetadataObjectCreation(t *testing.T) {
	obj := MetadataObject{
		Type:    ObjectTypeDocument,
		Name:    "TestDocument",
		Synonym: "Тестовый документ",
		Attributes: []Attribute{
			{Name: "Attribute1", Synonym: "Реквизит1", Types: []string{"String"}},
		},
		TabularSections: []TabularSection{
			{
				Name:    "TabularSection1",
				Synonym: "ТабличнаяЧасть1",
				Attributes: []Attribute{
					{Name: "TSAttribute1", Synonym: "ТС_Реквизит1", Types: []string{"Number"}},
				},
			},
		},
	}

	if obj.Name != "TestDocument" {
		t.Errorf("Expected Name to be 'TestDocument', but got '%s'", obj.Name)
	}

	if obj.Type != ObjectTypeDocument {
		t.Errorf("Expected Type to be 'Document', but got '%s'", obj.Type)
	}

	if len(obj.Attributes) != 1 {
		t.Errorf("Expected 1 attribute, but got %d", len(obj.Attributes))
	}

	if len(obj.TabularSections) != 1 {
		t.Errorf("Expected 1 tabular section, but got %d", len(obj.TabularSections))
	}
}

func TestAttributeCreation(t *testing.T) {
	attr := Attribute{
		Name:    "TestAttribute",
		Synonym: "Тестовый реквизит",
		Types:   []string{"String", "Number"},
	}

	if attr.Name != "TestAttribute" {
		t.Errorf("Expected Name to be 'TestAttribute', but got '%s'", attr.Name)
	}

	if len(attr.Types) != 2 {
		t.Errorf("Expected 2 types, but got %d", len(attr.Types))
	}
}

func TestTabularSectionCreation(t *testing.T) {
	ts := TabularSection{
		Name:    "TestTabularSection",
		Synonym: "Тестовая табличная часть",
		Attributes: []Attribute{
			{Name: "TSAttribute1", Synonym: "ТС_Реквизит1"},
		},
	}

	if ts.Name != "TestTabularSection" {
		t.Errorf("Expected Name to be 'TestTabularSection', but got '%s'", ts.Name)
	}

	if len(ts.Attributes) != 1 {
		t.Errorf("Expected 1 attribute, but got %d", len(ts.Attributes))
	}
}
