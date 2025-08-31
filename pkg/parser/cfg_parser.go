package parser

import (
	"encoding/xml"
	"fmt"
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
	XMLName xml.Name          `xml:"http://v8.1c.ru/8.3/MDClasses MetaDataObject"`
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

// CFGChildObjects дочерние объекты (атрибуты, табличные части, измерения, ресурсы)
type CFGChildObjects struct {
	Attributes      []CFGAttribute      `xml:"http://v8.1c.ru/8.3/MDClasses Attribute"`
	TabularSections []CFGTabularSection `xml:"http://v8.1c.ru/8.3/MDClasses TabularSection"`
	// Для регистров накопления и сведений
	Dimensions []CFGAttribute `xml:"http://v8.1c.ru/8.3/MDClasses Dimension"`
	Resources  []CFGAttribute `xml:"http://v8.1c.ru/8.3/MDClasses Resource"`
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
	data, err := os.ReadFile(filePath)
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
	data, err := os.ReadFile(filePath)
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
	enumsPath := filepath.Join(p.sourcePath, "Enums")
	if _, err := os.Stat(enumsPath); os.IsNotExist(err) {
		return []model.MetadataObject{}, nil
	}

	var enums []model.MetadataObject

	err := filepath.Walk(enumsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(info.Name()), ".xml") {
			return nil
		}

		data, rerr := os.ReadFile(path)
		if rerr != nil {
			fmt.Printf("Предупреждение: ошибка чтения файла перечисления %s: %v\n", path, rerr)
			return nil
		}

		// Структура для разбора перечисления
		type cfgEnum struct {
			XMLName xml.Name `xml:"http://v8.1c.ru/8.3/MDClasses MetaDataObject"`
			Enum    struct {
				Properties   CFGProperties `xml:"http://v8.1c.ru/8.3/MDClasses Properties"`
				ChildObjects struct {
					EnumValues []struct {
						Properties struct {
							Name    string     `xml:"http://v8.1c.ru/8.3/MDClasses Name"`
							Synonym CFGSynonym `xml:"http://v8.1c.ru/8.3/MDClasses Synonym"`
						} `xml:"http://v8.1c.ru/8.3/MDClasses Properties"`
					} `xml:"http://v8.1c.ru/8.3/MDClasses EnumValue"`
				} `xml:"http://v8.1c.ru/8.3/MDClasses ChildObjects"`
			} `xml:"http://v8.1c.ru/8.3/MDClasses Enum"`
		}

		var ce cfgEnum
		if err := xml.Unmarshal(data, &ce); err != nil {
			fmt.Printf("Предупреждение: ошибка парсинга перечисления %s: %v\n", path, err)
			return nil
		}

		name := ce.Enum.Properties.Name
		obj := model.MetadataObject{
			Type:    model.ObjectTypeEnum,
			Name:    name,
			Synonym: p.extractSynonym(ce.Enum.Properties.Synonym),
		}

		for _, v := range ce.Enum.ChildObjects.EnumValues {
			evName := v.Properties.Name
			evSyn := p.extractSynonym(v.Properties.Synonym)
			obj.EnumValues = append(obj.EnumValues, model.EnumValue{Name: evName, Synonym: evSyn})
		}

		enums = append(enums, obj)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка сканирования каталога перечислений: %w", err)
	}

	return enums, nil
}

// ParseChartsOfCharacteristicTypes парсит планы видов характеристик (заглушка для MVP)
func (p *CFGParser) ParseChartsOfCharacteristicTypes() ([]model.MetadataObject, error) {
	chartsPath := filepath.Join(p.sourcePath, "ChartsOfCharacteristicTypes")
	if _, err := os.Stat(chartsPath); os.IsNotExist(err) {
		return []model.MetadataObject{}, nil
	}

	var charts []model.MetadataObject

	err := filepath.Walk(chartsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(info.Name()), ".xml") {
			return nil
		}

		data, rerr := os.ReadFile(path)
		if rerr != nil {
			fmt.Printf("Предупреждение: ошибка чтения файла плана %s: %v\n", path, rerr)
			return nil
		}

		type cfgChart struct {
			XMLName xml.Name `xml:"http://v8.1c.ru/8.3/MDClasses MetaDataObject"`
			Chart   struct {
				Properties   CFGProperties `xml:"http://v8.1c.ru/8.3/MDClasses Properties"`
				ChildObjects struct {
					Attributes      []CFGAttribute `xml:"http://v8.1c.ru/8.3/MDClasses Attribute"`
					TabularSections []struct {
						Properties   CFGTabularSectionProperties `xml:"http://v8.1c.ru/8.3/MDClasses Properties"`
						ChildObjects CFGTabularSectionChilds     `xml:"http://v8.1c.ru/8.3/MDClasses ChildObjects"`
					} `xml:"http://v8.1c.ru/8.3/MDClasses TabularSection"`
				} `xml:"http://v8.1c.ru/8.3/MDClasses ChildObjects"`
			} `xml:"http://v8.1c.ru/8.3/MDClasses ChartOfCharacteristicTypes"`
		}

		var cc cfgChart
		if err := xml.Unmarshal(data, &cc); err != nil {
			fmt.Printf("Предупреждение: ошибка парсинга плана %s: %v\n", path, err)
			return nil
		}

		obj := model.MetadataObject{
			Type:    model.ObjectTypeChartOfCharacteristicTypes,
			Name:    cc.Chart.Properties.Name,
			Synonym: p.extractSynonym(cc.Chart.Properties.Synonym),
		}

		for _, a := range cc.Chart.ChildObjects.Attributes {
			types := p.extractTypes(a.Properties.Type)
			converted := p.typeConverter.ConvertTypes(types)
			obj.Attributes = append(obj.Attributes, model.Attribute{
				Name:    a.Properties.Name,
				Synonym: p.extractSynonym(a.Properties.Synonym),
				Types:   converted,
			})
		}

		for _, ts := range cc.Chart.ChildObjects.TabularSections {
			tab := model.TabularSection{
				Name:    ts.Properties.Name,
				Synonym: p.extractSynonym(ts.Properties.Synonym),
			}
			for _, a := range ts.ChildObjects.Attributes {
				types := p.extractTypes(a.Properties.Type)
				converted := p.typeConverter.ConvertTypes(types)
				tab.Attributes = append(tab.Attributes, model.Attribute{
					Name:    a.Properties.Name,
					Synonym: p.extractSynonym(a.Properties.Synonym),
					Types:   converted,
				})
			}
			obj.TabularSections = append(obj.TabularSections, tab)
		}

		charts = append(charts, obj)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка сканирования каталога планов видов характеристик: %w", err)
	}

	return charts, nil
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
		case model.ObjectTypeAccumulationRegister:
			regs, err := p.ParseAccumulationRegisters()
			if err != nil {
				return nil, err
			}
			allObjects = append(allObjects, regs...)

		case model.ObjectTypeInformationRegister:
			regs, err := p.ParseInformationRegisters()
			if err != nil {
				return nil, err
			}
			allObjects = append(allObjects, regs...)

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

		case model.ObjectTypeConstant:
			consts, err := p.ParseConstants()
			if err != nil {
				return nil, err
			}
			allObjects = append(allObjects, consts...)
		case model.ObjectTypeFilterCriteria:
			fcs, err := p.ParseFilterCriteria()
			if err != nil {
				return nil, err
			}
			allObjects = append(allObjects, fcs...)
		}
	}

	return allObjects, nil
}

// ParseAccumulationRegisters парсит регистры накопления в CFG формате
func (p *CFGParser) ParseAccumulationRegisters() ([]model.MetadataObject, error) {
	regsPath := filepath.Join(p.sourcePath, "AccumulationRegisters")
	if _, err := os.Stat(regsPath); os.IsNotExist(err) {
		return []model.MetadataObject{}, nil
	}
	var regs []model.MetadataObject
	err := filepath.Walk(regsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".xml") {
			reg, perr := p.parseAccumulationRegisterFile(path)
			if perr != nil {
				fmt.Printf("Предупреждение: ошибка парсинга регистра %s: %v\n", path, perr)
				return nil
			}
			regs = append(regs, reg)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка сканирования регистров накопления: %w", err)
	}
	return regs, nil
}

// ParseConstants парсит константы в CFG формате
func (p *CFGParser) ParseConstants() ([]model.MetadataObject, error) {
	constsPath := filepath.Join(p.sourcePath, "Constants")
	if _, err := os.Stat(constsPath); os.IsNotExist(err) {
		return []model.MetadataObject{}, nil
	}

	var consts []model.MetadataObject
	err := filepath.Walk(constsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(info.Name()), ".xml") {
			return nil
		}

		c, perr := p.parseConstantFile(path)
		if perr != nil {
			fmt.Printf("Предупреждение: ошибка парсинга константы %s: %v\n", path, perr)
			return nil
		}
		consts = append(consts, c)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка сканирования каталога констант: %w", err)
	}
	return consts, nil
}

// parseConstantFile парсит отдельный XML файл константы
func (p *CFGParser) parseConstantFile(filePath string) (model.MetadataObject, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка чтения файла %s: %w", filePath, err)
	}

	type cfgConstant struct {
		XMLName  xml.Name `xml:"http://v8.1c.ru/8.3/MDClasses MetaDataObject"`
		Constant struct {
			Properties CFGProperties `xml:"http://v8.1c.ru/8.3/MDClasses Properties"`
		} `xml:"http://v8.1c.ru/8.3/MDClasses Constant"`
	}

	var cc cfgConstant
	if err := xml.Unmarshal(data, &cc); err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка парсинга XML %s: %w", filePath, err)
	}

	name := cc.Constant.Properties.Name
	syn := p.extractSynonym(cc.Constant.Properties.Synonym)

	// Попробуем извлечь тип, если он есть
	var types []string
	// В CFG структура Properties может содержать Type, но CFGProperties не содержит Type напрямую
	// Попробуем распарсить вспомогательную структуру для получения Type
	type cfgConstantWithType struct {
		Constant struct {
			Properties struct {
				Name    string     `xml:"http://v8.1c.ru/8.3/MDClasses Name"`
				Synonym CFGSynonym `xml:"http://v8.1c.ru/8.3/MDClasses Synonym"`
				Type    CFGType    `xml:"http://v8.1c.ru/8.3/MDClasses Type"`
			} `xml:"http://v8.1c.ru/8.3/MDClasses Properties"`
		} `xml:"http://v8.1c.ru/8.3/MDClasses Constant"`
	}
	var ct cfgConstantWithType
	if err := xml.Unmarshal(data, &ct); err == nil {
		types = p.extractTypes(ct.Constant.Properties.Type)
	}

	converted := p.typeConverter.ConvertTypes(types)

	obj := model.MetadataObject{
		Type:    model.ObjectTypeConstant,
		Name:    name,
		Synonym: syn,
	}

	// Поместим информацию о значении константы как атрибут "Значение"
	obj.Attributes = append(obj.Attributes, model.Attribute{
		Name:    "Значение",
		Synonym: "",
		Types:   converted,
	})

	return obj, nil
}

// ParseFilterCriteria парсит критерии отбора в CFG формате
func (p *CFGParser) ParseFilterCriteria() ([]model.MetadataObject, error) {
	fcPath := filepath.Join(p.sourcePath, "FilterCriteria")
	if _, err := os.Stat(fcPath); os.IsNotExist(err) {
		return []model.MetadataObject{}, nil
	}

	var result []model.MetadataObject
	err := filepath.Walk(fcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(info.Name()), ".xml") {
			return nil
		}
		obj, perr := p.parseFilterCriteriaFile(path)
		if perr != nil {
			fmt.Printf("Предупреждение: ошибка парсинга критерия отбора %s: %v\n", path, perr)
			return nil
		}
		result = append(result, obj)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка сканирования каталога критериев отбора: %w", err)
	}
	return result, nil
}

// parseFilterCriteriaFile парсит отдельный XML файл критерия отбора
func (p *CFGParser) parseFilterCriteriaFile(filePath string) (model.MetadataObject, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка чтения файла %s: %w", filePath, err)
	}

	type cfgFilter struct {
		XMLName xml.Name `xml:"http://v8.1c.ru/8.3/MDClasses MetaDataObject"`
		Filter  struct {
			Properties   CFGProperties   `xml:"http://v8.1c.ru/8.3/MDClasses Properties"`
			ChildObjects CFGChildObjects `xml:"http://v8.1c.ru/8.3/MDClasses ChildObjects"`
		} `xml:"http://v8.1c.ru/8.3/MDClasses FilterCriteria"`
	}

	var cf cfgFilter
	if err := xml.Unmarshal(data, &cf); err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка парсинга XML %s: %w", filePath, err)
	}

	name := cf.Filter.Properties.Name
	syn := p.extractSynonym(cf.Filter.Properties.Synonym)

	// If the file uses singular FilterCriterion element (fixture), try parsing that too
	if name == "" {
		type cfgFilterSingular struct {
			XMLName xml.Name `xml:"http://v8.1c.ru/8.3/MDClasses MetaDataObject"`
			Filter  struct {
				Properties CFGProperties `xml:"http://v8.1c.ru/8.3/MDClasses Properties"`
			} `xml:"http://v8.1c.ru/8.3/MDClasses FilterCriterion"`
		}
		var cfs cfgFilterSingular
		if err := xml.Unmarshal(data, &cfs); err == nil {
			name = cfs.Filter.Properties.Name
			syn = p.extractSynonym(cfs.Filter.Properties.Synonym)
		}
	}

	obj := model.MetadataObject{
		Type:    model.ObjectTypeFilterCriteria,
		Name:    name,
		Synonym: syn,
	}

	// Критерии отбора не имеют реквизитов в нашей модели — пропускаем ChildObjects

	// Попытка извлечь Type из Properties (если он присутствует)
	// Для CFG мы попробуем распарсить Type внутри Properties при необходимости
	type probe struct {
		Filter struct {
			Properties struct {
				Type CFGType `xml:"http://v8.1c.ru/8.3/MDClasses Type"`
			} `xml:"http://v8.1c.ru/8.3/MDClasses Properties"`
		} `xml:"http://v8.1c.ru/8.3/MDClasses FilterCriteria"`
	}
	var pr probe
	if err := xml.Unmarshal(data, &pr); err == nil {
		if len(pr.Filter.Properties.Type.Types) > 0 || len(pr.Filter.Properties.Type.TypeSets) > 0 {
			types := p.extractTypes(pr.Filter.Properties.Type)
			obj.FilterCriteriaTypes = p.typeConverter.ConvertTypes(types)
		}
	}

	// Если probe не нашёл типов, попробуем более гибко извлечь содержимое элемента <Type> в любом пространстве имён
	if len(obj.FilterCriteriaTypes) == 0 {
		if vals := p.extractTypeValuesFromXML(data); len(vals) > 0 {
			// Убираем префиксы и конвертируем
			for _, v := range vals {
				v = strings.TrimSpace(v)
				v = strings.TrimPrefix(v, "cfg:")
				v = strings.TrimPrefix(v, "v8:")
				if v != "" {
					obj.FilterCriteriaTypes = append(obj.FilterCriteriaTypes, v)
				}
			}
			if len(obj.FilterCriteriaTypes) > 0 {
				obj.FilterCriteriaTypes = p.typeConverter.ConvertTypes(obj.FilterCriteriaTypes)
			}
		}
	}

	// Попытка извлечь Content элементов, если они присутствуют в файле
	// В CFG это редкость, но оставим поддержку
	// Попытка извлечь Content элементов, поддерживаем разные вложенные имена и пространства имён
	type contentAny struct {
		Filter struct {
			Content struct {
				Items []struct {
					Value string `xml:",chardata"`
				} `xml:",any"`
			} `xml:"Content"`
		} `xml:"http://v8.1c.ru/8.3/MDClasses FilterCriteria"`
	}
	var ca contentAny
	if err := xml.Unmarshal(data, &ca); err == nil && len(ca.Filter.Content.Items) > 0 {
		for _, it := range ca.Filter.Content.Items {
			if it.Value != "" {
				obj.FilterCriteriaContents = append(obj.FilterCriteriaContents, NormalizeFilterContentItem(it.Value))
			}
		}
	} else {
		// try singular FilterCriterion path
		type contentAnySingular struct {
			Filter struct {
				Content struct {
					Items []struct {
						Value string `xml:",chardata"`
					} `xml:",any"`
				} `xml:"Content"`
			} `xml:"http://v8.1c.ru/8.3/MDClasses FilterCriterion"`
		}
		var cas contentAnySingular
		if err := xml.Unmarshal(data, &cas); err == nil && len(cas.Filter.Content.Items) > 0 {
			for _, it := range cas.Filter.Content.Items {
				if it.Value != "" {
					obj.FilterCriteriaContents = append(obj.FilterCriteriaContents, NormalizeFilterContentItem(it.Value))
				}
			}
		}
	}

	// If still empty, try a generic XML token walk to extract Content/Item char data
	if len(obj.FilterCriteriaContents) == 0 {
		if vals := p.extractContentValuesFromXML(data); len(vals) > 0 {
			for _, v := range vals {
				obj.FilterCriteriaContents = append(obj.FilterCriteriaContents, NormalizeFilterContentItem(v))
			}
		}
	}

	return obj, nil
}

// extractTypeValuesFromXML tries to find <Type> element text anywhere in the XML (any namespace)
func (p *CFGParser) extractTypeValuesFromXML(data []byte) []string {
	dec := xml.NewDecoder(strings.NewReader(string(data)))
	var results []string

	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch se := tok.(type) {
		case xml.StartElement:
			if se.Name.Local == "Type" {
				// Collect any character data inside this <Type> element.
				// It may contain a nested <v8:Type> element, so we iterate tokens
				// until we reach the corresponding EndElement.
				var sb strings.Builder
				for {
					innerTok, ierr := dec.Token()
					if ierr != nil {
						break
					}
					switch it := innerTok.(type) {
					case xml.CharData:
						s := strings.TrimSpace(string(it))
						if s != "" {
							if sb.Len() > 0 {
								sb.WriteByte(' ')
							}
							sb.WriteString(s)
						}
					case xml.StartElement:
						// If nested element (like v8:Type) contains char data, decode it directly
						if it.Name.Local == "Type" {
							var inner string
							if err2 := dec.DecodeElement(&inner, &it); err2 == nil {
								inner = strings.TrimSpace(inner)
								if inner != "" {
									if sb.Len() > 0 {
										sb.WriteByte(' ')
									}
									sb.WriteString(inner)
								}
							}
						}
					case xml.EndElement:
						if it.Name.Local == "Type" {
							// finished this Type element
							goto FLUSH
						}
					}
				}
			FLUSH:
				v := strings.TrimSpace(sb.String())
				if v != "" {
					// split by whitespace in case of multiple values
					for _, part := range strings.Fields(v) {
						if strings.TrimSpace(part) != "" {
							results = append(results, part)
						}
					}
				}
			}
		}
	}

	return results
}

// ParseInformationRegisters парсит регистры сведений в CFG формате
func (p *CFGParser) ParseInformationRegisters() ([]model.MetadataObject, error) {
	regsPath := filepath.Join(p.sourcePath, "InformationRegisters")
	if _, err := os.Stat(regsPath); os.IsNotExist(err) {
		return []model.MetadataObject{}, nil
	}

	var regs []model.MetadataObject

	err := filepath.Walk(regsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".xml") {
			reg, parseErr := p.parseInformationRegisterFile(path)
			if parseErr != nil {
				fmt.Printf("Предупреждение: ошибка парсинга регистра %s: %v\n", path, parseErr)
				return nil
			}
			regs = append(regs, reg)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка сканирования регистров сведений: %w", err)
	}

	return regs, nil
}

// extractContentValuesFromXML walks the XML and collects character data inside Content->Item elements
func (p *CFGParser) extractContentValuesFromXML(data []byte) []string {
	dec := xml.NewDecoder(strings.NewReader(string(data)))
	var results []string
	var inContent bool
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == "Content" {
				inContent = true
			} else if inContent && (t.Name.Local == "Item" || t.Name.Local == "xr:Item") {
				// decode inner char data of Item
				var v string
				if err := dec.DecodeElement(&v, &t); err == nil {
					v = strings.TrimSpace(v)
					if v != "" {
						results = append(results, v)
					}
				}
			}
		case xml.EndElement:
			if t.Name.Local == "Content" {
				inContent = false
			}
		}
	}
	return results
}

// use NormalizeFilterContentItem from normalize.go

// parseInformationRegisterFile парсит один XML файл регистра сведений
func (p *CFGParser) parseInformationRegisterFile(filePath string) (model.MetadataObject, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка чтения файла %s: %w", filePath, err)
	}

	// CFGRegister структура для разбора регистра из XML
	type CFGRegister struct {
		XMLName             xml.Name `xml:"http://v8.1c.ru/8.3/MDClasses MetaDataObject"`
		InformationRegister struct {
			Properties struct {
				Name    string     `xml:"http://v8.1c.ru/8.3/MDClasses Name"`
				Synonym CFGSynonym `xml:"http://v8.1c.ru/8.3/MDClasses Synonym"`
			} `xml:"http://v8.1c.ru/8.3/MDClasses Properties"`
			ChildObjects struct {
				Dimensions []CFGAttribute `xml:"http://v8.1c.ru/8.3/MDClasses Dimension"`
				Resources  []CFGAttribute `xml:"http://v8.1c.ru/8.3/MDClasses Resource"`
				Attributes []CFGAttribute `xml:"http://v8.1c.ru/8.3/MDClasses Attribute"`
			} `xml:"http://v8.1c.ru/8.3/MDClasses ChildObjects"`
		} `xml:"http://v8.1c.ru/8.3/MDClasses InformationRegister"`
	}

	var reg CFGRegister
	if err := xml.Unmarshal(data, &reg); err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка парсинга XML файла %s: %w", filePath, err)
	}

	result := model.MetadataObject{
		Type:    model.ObjectTypeInformationRegister,
		Name:    reg.InformationRegister.Properties.Name,
		Synonym: p.extractSynonym(reg.InformationRegister.Properties.Synonym),
	}

	// Измерения
	for _, d := range reg.InformationRegister.ChildObjects.Dimensions {
		types := p.extractTypes(d.Properties.Type)
		converted := p.typeConverter.ConvertTypes(types)
		result.Dimensions = append(result.Dimensions, model.Attribute{
			Name:    d.Properties.Name,
			Synonym: p.extractSynonym(d.Properties.Synonym),
			Types:   converted,
		})
	}

	// Ресурсы
	for _, r := range reg.InformationRegister.ChildObjects.Resources {
		types := p.extractTypes(r.Properties.Type)
		converted := p.typeConverter.ConvertTypes(types)
		result.Resources = append(result.Resources, model.Attribute{
			Name:    r.Properties.Name,
			Synonym: p.extractSynonym(r.Properties.Synonym),
			Types:   converted,
		})
	}

	// Реквизиты
	for _, a := range reg.InformationRegister.ChildObjects.Attributes {
		types := p.extractTypes(a.Properties.Type)
		converted := p.typeConverter.ConvertTypes(types)
		result.Attributes = append(result.Attributes, model.Attribute{
			Name:    a.Properties.Name,
			Synonym: p.extractSynonym(a.Properties.Synonym),
			Types:   converted,
		})
	}

	return result, nil
}

// parseAccumulationRegisterFile парсит один XML файл регистра накопления
func (p *CFGParser) parseAccumulationRegisterFile(filePath string) (model.MetadataObject, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка чтения файла %s: %w", filePath, err)
	}

	// Определяем корень по пространству имен так же, как и для других объектов
	type cfgRegContent struct {
		Properties   CFGProperties   `xml:"http://v8.1c.ru/8.3/MDClasses Properties"`
		ChildObjects CFGChildObjects `xml:"http://v8.1c.ru/8.3/MDClasses ChildObjects"`
	}
	type cfgReg struct {
		XMLName  xml.Name      `xml:"http://v8.1c.ru/8.3/MDClasses MetaDataObject"`
		Register cfgRegContent `xml:"http://v8.1c.ru/8.3/MDClasses AccumulationRegister"`
	}

	var reg cfgReg
	if err := xml.Unmarshal(data, &reg); err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка парсинга XML файла %s: %w", filePath, err)
	}

	result := model.MetadataObject{
		Type:    model.ObjectTypeAccumulationRegister,
		Name:    reg.Register.Properties.Name,
		Synonym: p.extractSynonym(reg.Register.Properties.Synonym),
	}

	// Измерения
	for _, d := range reg.Register.ChildObjects.Dimensions {
		types := p.extractTypes(d.Properties.Type)
		converted := p.typeConverter.ConvertTypes(types)
		result.Dimensions = append(result.Dimensions, model.Attribute{
			Name:    d.Properties.Name,
			Synonym: p.extractSynonym(d.Properties.Synonym),
			Types:   converted,
		})
	}
	// Ресурсы
	for _, r := range reg.Register.ChildObjects.Resources {
		types := p.extractTypes(r.Properties.Type)
		converted := p.typeConverter.ConvertTypes(types)
		result.Resources = append(result.Resources, model.Attribute{
			Name:    r.Properties.Name,
			Synonym: p.extractSynonym(r.Properties.Synonym),
			Types:   converted,
		})
	}
	// Реквизиты
	for _, a := range reg.Register.ChildObjects.Attributes {
		types := p.extractTypes(a.Properties.Type)
		converted := p.typeConverter.ConvertTypes(types)
		result.Attributes = append(result.Attributes, model.Attribute{
			Name:    a.Properties.Name,
			Synonym: p.extractSynonym(a.Properties.Synonym),
			Types:   converted,
		})
	}
	return result, nil
}
