package model

// MetadataObject представляет объект метаданных 1С
type MetadataObject struct {
	Type       ObjectType      `json:"type"`
	Name       string          `json:"name"`
	Synonym    string          `json:"synonym"`
    // Для документов и справочников
    Attributes []Attribute     `json:"attributes"`
    TabularSections []TabularSection `json:"tabular_sections"`
    // Для регистров накопления
    Dimensions []Attribute     `json:"dimensions"`
    Resources  []Attribute     `json:"resources"`
}

// ObjectType определяет тип объекта метаданных
type ObjectType string

const (
	ObjectTypeDocument                ObjectType = "Document"
	ObjectTypeCatalog                ObjectType = "Catalog"
	ObjectTypeEnum                   ObjectType = "Enum"
	ObjectTypeChartOfCharacteristicTypes ObjectType = "ChartOfCharacteristicTypes"
    ObjectTypeAccumulationRegister   ObjectType = "AccumulationRegister"
)

// Attribute представляет реквизит объекта
type Attribute struct {
	Name     string   `json:"name"`
	Synonym  string   `json:"synonym"`
	Types    []string `json:"types"`
	Required bool     `json:"required"`
}

// TabularSection представляет табличную часть
type TabularSection struct {
	Name       string      `json:"name"`
	Synonym    string      `json:"synonym"`
	Attributes []Attribute `json:"attributes"`
}

// SourceFormat определяет формат исходных данных
type SourceFormat string

const (
	FormatCFG SourceFormat = "cfg"
	FormatEDT SourceFormat = "edt"
)

// ConversionOptions опции конвертации
type ConversionOptions struct {
	SourcePath   string       `json:"source_path"`
	OutputPath   string       `json:"output_path"`
	Format       SourceFormat `json:"format"`
	ObjectTypes  []ObjectType `json:"object_types"`
	Verbose      bool         `json:"verbose"`
}

// CatalogEntry запись в каталоге объектов
type CatalogEntry struct {
	ObjectName string `json:"object_name"`
	ObjectType string `json:"object_type"`
	Synonym    string `json:"synonym"`
	FileName   string `json:"file_name"`
}