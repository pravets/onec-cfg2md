package parser

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"ones-cfg2md/pkg/model"
)

// CFGParser парсер для CFG формата
type CFGParser struct {
	sourcePath    string
	typeConverter TypeConverter
}

// NewCFGParser создает новый парсер CFG формата
func NewCFGParser(sourcePath string) (*CFGParser, error) {
	return &CFGParser{
		sourcePath:    sourcePath,
		typeConverter: NewTypeConverter(),
	}, nil
}

// CFGDocument структура для парсинга CFG документа
type CFGDocument struct {
    XMLName  xml.Name           `xml:"http://v8.1c.ru/8.3/MDClasses MetaDataObject"`
    Document CFGDocumentContent `xml:"http://v8.1c.ru/8.3/MDClasses Document"`
}

// CFGCatalog структура для парсинга CFG справочника
type CFGCatalog struct {
    XMLName xml.Name         `xml:"http://v8.1c.ru/8.3/MDClasses MetaDataObject"`
    Catalog CFGCatalogContent `xml:"http://v8.1c.ru/8.3/MDClasses Catalog"`
}

// CFGDocumentContent содержимое документа в CFG формате
type CFGDocumentContent struct {
    Properties   CFGProperties   `xml:"http://v8.1c.ru/8.3/MDClasses Properties"`
    ChildObjects CFGChildObjects `xml:"http://v8.1c.ru/8.3/MDClasses ChildObjects"`
}

// CFGCatalogContent содержимое справочника в CFG формате
type CFGCatalogContent struct {
    Properties   CFGProperties   `xml:"http://v8.1c.ru/8.3/MDClasses Properties"`
    ChildObjects CFGChildObjects `xml:"http://v8.1c.ru/8.3/MDClasses ChildObjects"`
}

// CFGProperties свойства документа
type CFGProperties struct {
    Name    string     `xml:"http://v8.1c.ru/8.3/MDClasses Name"`
    Synonym CFGSynonym `xml:"http://v8.1c.ru/8.3/MDClasses Synonym"`
}

// CFGSynonym синоним в CFG формате
type CFGSynonym struct {
	Items []CFGSynonymItem `xml:"http://v8.1c.ru/8.1/data/core item"`
}

// CFGSynonymItem элемент синонима
type CFGSynonymItem struct {
	Lang    string `xml:"http://v8.1c.ru/8.1/data/core lang"`
	Content string `xml:"http://v8.1c.ru/8.1/data/core content"`
}

// CFGChildObjects дочерние объекты (атрибуты, табличные части)
type CFGChildObjects struct {
    Attributes      []CFGAttribute      `xml:"http://v8.1c.ru/8.3/MDClasses Attribute"`
    TabularSections []CFGTabularSection `xml:"http://v8.1c.ru/8.3/MDClasses TabularSection"`
}

// CFGAttribute атрибут в CFG формате
type CFGAttribute struct {
    Properties CFGAttributeProperties `xml:"http://v8.1c.ru/8.3/MDClasses Properties"`
}

// CFGAttributeProperties свойства атрибута
type CFGAttributeProperties struct {
    Name    string     `xml:"http://v8.1c.ru/8.3/MDClasses Name"`
    Synonym CFGSynonym `xml:"http://v8.1c.ru/8.3/MDClasses Synonym"`
    Type    CFGType    `xml:"http://v8.1c.ru/8.3/MDClasses Type"`
}

// CFGType тип в CFG формате
type CFGType struct {
	// В CFG формате типы могут быть заданы как одиночные элементы или как массивы
	Types    []string `xml:"http://v8.1c.ru/8.1/data/core Type"`
	TypeSets []string `xml:"http://v8.1c.ru/8.1/data/core TypeSet"`
    // Квалификаторы даты позволяют различать дату и дату-время
    DateQualifiers []CFGDateQualifiers `xml:"http://v8.1c.ru/8.1/data/core DateQualifiers"`
}

// CFGDateQualifiers квалификаторы для дат
type CFGDateQualifiers struct {
    DateFractions string `xml:"http://v8.1c.ru/8.1/data/core DateFractions"`
}

// CFGTabularSection табличная часть в CFG формате
type CFGTabularSection struct {
    Properties   CFGTabularSectionProperties `xml:"http://v8.1c.ru/8.3/MDClasses Properties"`
    ChildObjects CFGTabularSectionChilds     `xml:"http://v8.1c.ru/8.3/MDClasses ChildObjects"`
}

// CFGTabularSectionProperties свойства табличной части
type CFGTabularSectionProperties struct {
    Name    string     `xml:"http://v8.1c.ru/8.3/MDClasses Name"`
    Synonym CFGSynonym `xml:"http://v8.1c.ru/8.3/MDClasses Synonym"`
}

// CFGTabularSectionChilds дочерние объекты табличной части
type CFGTabularSectionChilds struct {
    Attributes []CFGAttribute `xml:"http://v8.1c.ru/8.3/MDClasses Attribute"`
}

// ParseDocuments парсит все документы в CFG формате
func (p *CFGParser) ParseDocuments() ([]model.MetadataObject, error) {
	documentsPath := filepath.Join(p.sourcePath, "Documents")
	
	// Проверяем существование каталога Documents
	if _, err := os.Stat(documentsPath); os.IsNotExist(err) {
		return []model.MetadataObject{}, nil // Нет документов
	}
	
	var documents []model.MetadataObject
	
	// Сканируем XML файлы документов
	err := filepath.Walk(documentsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".xml") {
			document, parseErr := p.parseDocumentFile(path)
			if parseErr != nil {
				// Логируем ошибку, но продолжаем обработку других документов
				fmt.Printf("Предупреждение: ошибка парсинга документа %s: %v\n", path, parseErr)
				return nil
			}
			documents = append(documents, document)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("ошибка сканирования каталога документов: %w", err)
	}
	
	return documents, nil
}

// parseDocumentFile парсит отдельный XML файл документа в CFG формате
func (p *CFGParser) parseDocumentFile(filePath string) (model.MetadataObject, error) {
	// Читаем файл
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка чтения файла %s: %w", filePath, err)
	}
	
	// Парсим XML
	var cfgDoc CFGDocument
	if err := xml.Unmarshal(data, &cfgDoc); err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка парсинга XML файла %s: %w", filePath, err)
	}
	
	// Преобразуем в нашу модель
	document := model.MetadataObject{
		Type:    model.ObjectTypeDocument,
		Name:    cfgDoc.Document.Properties.Name,
		Synonym: p.extractSynonym(cfgDoc.Document.Properties.Synonym),
	}
	
	// Парсим атрибуты
	for _, attr := range cfgDoc.Document.ChildObjects.Attributes {
		types := p.extractTypes(attr.Properties.Type)
		convertedTypes := p.typeConverter.ConvertTypes(types)
		
		document.Attributes = append(document.Attributes, model.Attribute{
			Name:    attr.Properties.Name,
			Synonym: p.extractSynonym(attr.Properties.Synonym),
			Types:   convertedTypes,
		})
	}
	
	// Парсим табличные части
	for _, ts := range cfgDoc.Document.ChildObjects.TabularSections {
		tabularSection := model.TabularSection{
			Name:    ts.Properties.Name,
			Synonym: p.extractSynonym(ts.Properties.Synonym),
		}
		
		// Парсим атрибуты табличной части
		for _, attr := range ts.ChildObjects.Attributes {
			types := p.extractTypes(attr.Properties.Type)
			convertedTypes := p.typeConverter.ConvertTypes(types)
			
			tabularSection.Attributes = append(tabularSection.Attributes, model.Attribute{
				Name:    attr.Properties.Name,
				Synonym: p.extractSynonym(attr.Properties.Synonym),
				Types:   convertedTypes,
			})
		}
		
		document.TabularSections = append(document.TabularSections, tabularSection)
	}
	
	return document, nil
}

// extractSynonym извлекает русский синоним из структуры CFGSynonym
func (p *CFGParser) extractSynonym(synonym CFGSynonym) string {
	for _, item := range synonym.Items {
		if item.Lang == "ru" {
			return item.Content
		}
	}
	return ""
}

// extractTypes извлекает типы из структуры CFGType
func (p *CFGParser) extractTypes(typeInfo CFGType) []string {
	var types []string
	
	// Обрабатываем Types (могут быть как одиночными элементами, так и массивами)
	for _, typeStr := range typeInfo.Types {
		if strings.TrimSpace(typeStr) != "" {
			// Убираем префиксы cfg: и v8: если они есть
			typeStr = strings.TrimPrefix(typeStr, "cfg:")
			typeStr = strings.TrimPrefix(typeStr, "v8:")
            typeStr = strings.TrimSpace(typeStr)
            // Особая обработка для xs:dateTime с квалификатором Date -> Date
            if typeStr == "xs:dateTime" && p.isDateOnly(typeInfo) {
                types = append(types, "Date")
            } else {
                types = append(types, typeStr)
            }
		}
	}
	
	// Обрабатываем TypeSets
	for _, typeStr := range typeInfo.TypeSets {
		if strings.TrimSpace(typeStr) != "" {
			// Убираем префиксы cfg: и v8: если они есть
			typeStr = strings.TrimPrefix(typeStr, "cfg:")
			typeStr = strings.TrimPrefix(typeStr, "v8:")
			types = append(types, strings.TrimSpace(typeStr))
		}
	}
	
	return types
}

// isDateOnly определяет, что для типа указана только дата (без времени)
func (p *CFGParser) isDateOnly(typeInfo CFGType) bool {
    for _, dq := range typeInfo.DateQualifiers {
        // Если явно задана только дата
        if strings.EqualFold(strings.TrimSpace(dq.DateFractions), "Date") {
            return true
        }
    }
    return false
}

// ParseCatalogs парсит все справочники в CFG формате
func (p *CFGParser) ParseCatalogs() ([]model.MetadataObject, error) {
	catalogsPath := filepath.Join(p.sourcePath, "Catalogs")
	
	// Проверяем существование каталога Catalogs
	if _, err := os.Stat(catalogsPath); os.IsNotExist(err) {
		return []model.MetadataObject{}, nil // Нет справочников
	}
	
	var catalogs []model.MetadataObject
	
	// Сканируем XML файлы справочников
	err := filepath.Walk(catalogsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".xml") {
			catalog, parseErr := p.parseCatalogFile(path)
			if parseErr != nil {
				// Логируем ошибку, но продолжаем обработку других справочников
				fmt.Printf("Предупреждение: ошибка парсинга справочника %s: %v\n", path, parseErr)
				return nil
			}
			catalogs = append(catalogs, catalog)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("ошибка сканирования каталога справочников: %w", err)
	}
	
	return catalogs, nil
}

// parseCatalogFile парсит отдельный XML файл справочника в CFG формате
func (p *CFGParser) parseCatalogFile(filePath string) (model.MetadataObject, error) {
	// Читаем файл
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка чтения файла %s: %w", filePath, err)
	}
	
	// Парсим XML
	var cfgCatalog CFGCatalog
	if err := xml.Unmarshal(data, &cfgCatalog); err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка парсинга XML файла %s: %w", filePath, err)
	}
	
	// Преобразуем в нашу модель
	catalog := model.MetadataObject{
		Type:    model.ObjectTypeCatalog,
		Name:    cfgCatalog.Catalog.Properties.Name,
		Synonym: p.extractSynonym(cfgCatalog.Catalog.Properties.Synonym),
	}
	
	// Парсим атрибуты
	for _, attr := range cfgCatalog.Catalog.ChildObjects.Attributes {
		types := p.extractTypes(attr.Properties.Type)
		convertedTypes := p.typeConverter.ConvertTypes(types)
		
		catalog.Attributes = append(catalog.Attributes, model.Attribute{
			Name:    attr.Properties.Name,
			Synonym: p.extractSynonym(attr.Properties.Synonym),
			Types:   convertedTypes,
		})
	}
	
	// Парсим табличные части
	for _, ts := range cfgCatalog.Catalog.ChildObjects.TabularSections {
		tabularSection := model.TabularSection{
			Name:    ts.Properties.Name,
			Synonym: p.extractSynonym(ts.Properties.Synonym),
		}
		
		// Парсим атрибуты табличной части
		for _, attr := range ts.ChildObjects.Attributes {
			types := p.extractTypes(attr.Properties.Type)
			convertedTypes := p.typeConverter.ConvertTypes(types)
			
			tabularSection.Attributes = append(tabularSection.Attributes, model.Attribute{
				Name:    attr.Properties.Name,
				Synonym: p.extractSynonym(attr.Properties.Synonym),
				Types:   convertedTypes,
			})
		}
		
		catalog.TabularSections = append(catalog.TabularSections, tabularSection)
	}
	
	return catalog, nil
}

// ParseEnums парсит перечисления (заглушка для MVP)
func (p *CFGParser) ParseEnums() ([]model.MetadataObject, error) {
	return []model.MetadataObject{}, nil
}

// ParseChartsOfCharacteristicTypes парсит планы видов характеристик (заглушка для MVP)
func (p *CFGParser) ParseChartsOfCharacteristicTypes() ([]model.MetadataObject, error) {
	return []model.MetadataObject{}, nil
}

// ParseObjectsByType парсит объекты указанных типов
func (p *CFGParser) ParseObjectsByType(objectTypes []model.ObjectType) ([]model.MetadataObject, error) {
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