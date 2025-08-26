package parser

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"ones-cfg2md/pkg/model"
)

// EDTParser парсер для EDT формата
type EDTParser struct {
	sourcePath    string
	typeConverter TypeConverter
}

// NewEDTParser создает новый парсер EDT формата
func NewEDTParser(sourcePath string) (*EDTParser, error) {
	return &EDTParser{
		sourcePath:    sourcePath,
		typeConverter: NewTypeConverter(),
	}, nil
}

// EDTDocument структура для парсинга EDT документа
type EDTDocument struct {
	XMLName           xml.Name           `xml:"http://g5.1c.ru/v8/dt/metadata/mdclass Document"`
	Name              string             `xml:"name"`
	Synonym           EDTSynonym         `xml:"synonym"`
	Attributes        []EDTAttribute     `xml:"attributes"`
	TabularSections   []EDTTabularSection `xml:"tabularSections"`
}

// EDTCatalog структура для парсинга EDT справочника
type EDTCatalog struct {
	XMLName           xml.Name           `xml:"http://g5.1c.ru/v8/dt/metadata/mdclass Catalog"`
	Name              string             `xml:"name"`
	Synonym           EDTSynonym         `xml:"synonym"`
	Attributes        []EDTAttribute     `xml:"attributes"`
	TabularSections   []EDTTabularSection `xml:"tabularSections"`
}

// EDTSynonym синоним в EDT формате
type EDTSynonym struct {
	Key   string `xml:"key"`
	Value string `xml:"value"`
}

// EDTAttribute атрибут в EDT формате
type EDTAttribute struct {
	Name    string     `xml:"name"`
	Synonym EDTSynonym `xml:"synonym"`
	Type    EDTType    `xml:"type"`
}

// EDTType тип атрибута в EDT формате
type EDTType struct {
	Types []string `xml:"types"`
}

// EDTTabularSection табличная часть в EDT формате
type EDTTabularSection struct {
	Name       string         `xml:"name"`
	Synonym    EDTSynonym     `xml:"synonym"`
	Attributes []EDTAttribute `xml:"attributes"`
}

// ParseDocuments парсит все документы в EDT формате
func (p *EDTParser) ParseDocuments() ([]model.MetadataObject, error) {
	documentsPath := filepath.Join(p.sourcePath, "src", "Documents")
	
	// Проверяем существование каталога Documents
	if _, err := os.Stat(documentsPath); os.IsNotExist(err) {
		return []model.MetadataObject{}, nil // Нет документов
	}
	
	var documents []model.MetadataObject
	
	// Сканируем каталоги документов
	entries, err := ioutil.ReadDir(documentsPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения каталога документов: %w", err)
	}
	
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		documentName := entry.Name()
		mdoFile := filepath.Join(documentsPath, documentName, documentName+".mdo")
		
		// Проверяем существование MDO файла
		if _, err := os.Stat(mdoFile); os.IsNotExist(err) {
			continue
		}
		
		document, err := p.parseDocumentFile(mdoFile)
		if err != nil {
			// Логируем ошибку, но продолжаем обработку других документов
			fmt.Printf("Предупреждение: ошибка парсинга документа %s: %v\n", documentName, err)
			continue
		}
		
		documents = append(documents, document)
	}
	
	return documents, nil
}

// parseDocumentFile парсит отдельный MDO файл документа
func (p *EDTParser) parseDocumentFile(filePath string) (model.MetadataObject, error) {
	// Читаем файл
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка чтения файла %s: %w", filePath, err)
	}
	
	// Парсим XML
	var edtDoc EDTDocument
	if err := xml.Unmarshal(data, &edtDoc); err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка парсинга XML файла %s: %w", filePath, err)
	}
	
	// Преобразуем в нашу модель
	document := model.MetadataObject{
		Type:    model.ObjectTypeDocument,
		Name:    edtDoc.Name,
		Synonym: edtDoc.Synonym.Value,
	}
	
	// Парсим атрибуты
	for _, attr := range edtDoc.Attributes {
		convertedTypes := p.typeConverter.ConvertTypes(attr.Type.Types)
		document.Attributes = append(document.Attributes, model.Attribute{
			Name:    attr.Name,
			Synonym: attr.Synonym.Value,
			Types:   convertedTypes,
		})
	}
	
	// Парсим табличные части
	for _, ts := range edtDoc.TabularSections {
		tabularSection := model.TabularSection{
			Name:    ts.Name,
			Synonym: ts.Synonym.Value,
		}
		
		// Парсим атрибуты табличной части
		for _, attr := range ts.Attributes {
			convertedTypes := p.typeConverter.ConvertTypes(attr.Type.Types)
			tabularSection.Attributes = append(tabularSection.Attributes, model.Attribute{
				Name:    attr.Name,
				Synonym: attr.Synonym.Value,
				Types:   convertedTypes,
			})
		}
		
		document.TabularSections = append(document.TabularSections, tabularSection)
	}
	
	return document, nil
}

// ParseCatalogs парсит все справочники в EDT формате
func (p *EDTParser) ParseCatalogs() ([]model.MetadataObject, error) {
	catalogsPath := filepath.Join(p.sourcePath, "src", "Catalogs")
	
	// Проверяем существование каталога Catalogs
	if _, err := os.Stat(catalogsPath); os.IsNotExist(err) {
		return []model.MetadataObject{}, nil // Нет справочников
	}
	
	var catalogs []model.MetadataObject
	
	// Сканируем каталоги справочников
	entries, err := ioutil.ReadDir(catalogsPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения каталога справочников: %w", err)
	}
	
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		catalogName := entry.Name()
		mdoFile := filepath.Join(catalogsPath, catalogName, catalogName+".mdo")
		
		// Проверяем существование MDO файла
		if _, err := os.Stat(mdoFile); os.IsNotExist(err) {
			continue
		}
		
		catalog, err := p.parseCatalogFile(mdoFile)
		if err != nil {
			// Логируем ошибку, но продолжаем обработку других справочников
			fmt.Printf("Предупреждение: ошибка парсинга справочника %s: %v\n", catalogName, err)
			continue
		}
		
		catalogs = append(catalogs, catalog)
	}
	
	return catalogs, nil
}

// parseCatalogFile парсит отдельный MDO файл справочника
func (p *EDTParser) parseCatalogFile(filePath string) (model.MetadataObject, error) {
	// Читаем файл
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка чтения файла %s: %w", filePath, err)
	}
	
	// Парсим XML
	var edtCatalog EDTCatalog
	if err := xml.Unmarshal(data, &edtCatalog); err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка парсинга XML файла %s: %w", filePath, err)
	}
	
	// Преобразуем в нашу модель
	catalog := model.MetadataObject{
		Type:    model.ObjectTypeCatalog,
		Name:    edtCatalog.Name,
		Synonym: edtCatalog.Synonym.Value,
	}
	
	// Парсим атрибуты
	for _, attr := range edtCatalog.Attributes {
		convertedTypes := p.typeConverter.ConvertTypes(attr.Type.Types)
		catalog.Attributes = append(catalog.Attributes, model.Attribute{
			Name:    attr.Name,
			Synonym: attr.Synonym.Value,
			Types:   convertedTypes,
		})
	}
	
	// Парсим табличные части
	for _, ts := range edtCatalog.TabularSections {
		tabularSection := model.TabularSection{
			Name:    ts.Name,
			Synonym: ts.Synonym.Value,
		}
		
		// Парсим атрибуты табличной части
		for _, attr := range ts.Attributes {
			convertedTypes := p.typeConverter.ConvertTypes(attr.Type.Types)
			tabularSection.Attributes = append(tabularSection.Attributes, model.Attribute{
				Name:    attr.Name,
				Synonym: attr.Synonym.Value,
				Types:   convertedTypes,
			})
		}
		
		catalog.TabularSections = append(catalog.TabularSections, tabularSection)
	}
	
	return catalog, nil
}

// ParseEnums парсит перечисления (заглушка для MVP)
func (p *EDTParser) ParseEnums() ([]model.MetadataObject, error) {
	return []model.MetadataObject{}, nil
}

// ParseChartsOfCharacteristicTypes парсит планы видов характеристик (заглушка для MVP)
func (p *EDTParser) ParseChartsOfCharacteristicTypes() ([]model.MetadataObject, error) {
	return []model.MetadataObject{}, nil
}

// ParseObjectsByType парсит объекты указанных типов
func (p *EDTParser) ParseObjectsByType(objectTypes []model.ObjectType) ([]model.MetadataObject, error) {
	var allObjects []model.MetadataObject
	
	for _, objType := range objectTypes {
		switch objType {
		case model.ObjectTypeDocument:
			docs, err := p.ParseDocuments()
			if err != nil {
				return nil, err
			}
			allObjects = append(allObjects, docs...)
			
		case model.ObjectTypeCatalog:
			catalogs, err := p.ParseCatalogs()
			if err != nil {
				return nil, err
			}
			allObjects = append(allObjects, catalogs...)
			
		case model.ObjectTypeEnum:
			enums, err := p.ParseEnums()
			if err != nil {
				return nil, err
			}
			allObjects = append(allObjects, enums...)
			
		case model.ObjectTypeChartOfCharacteristicTypes:
			charts, err := p.ParseChartsOfCharacteristicTypes()
			if err != nil {
				return nil, err
			}
			allObjects = append(allObjects, charts...)
		}
	}
	
	return allObjects, nil
}