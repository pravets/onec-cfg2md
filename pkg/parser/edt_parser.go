package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"onec-cfg2md/pkg/model"
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
	XMLName         xml.Name            `xml:"http://g5.1c.ru/v8/dt/metadata/mdclass Document"`
	Name            string              `xml:"name"`
	Synonym         EDTSynonym          `xml:"synonym"`
	Attributes      []EDTAttribute      `xml:"attributes"`
	TabularSections []EDTTabularSection `xml:"tabularSections"`
}

// EDTCatalog структура для парсинга EDT справочника
type EDTCatalog struct {
	XMLName         xml.Name            `xml:"http://g5.1c.ru/v8/dt/metadata/mdclass Catalog"`
	Name            string              `xml:"name"`
	Synonym         EDTSynonym          `xml:"synonym"`
	Attributes      []EDTAttribute      `xml:"attributes"`
	TabularSections []EDTTabularSection `xml:"tabularSections"`
}

// EDTAccumulationRegister структура для парсинга EDT регистра накопления
type EDTAccumulationRegister struct {
	XMLName    xml.Name       `xml:"http://g5.1c.ru/v8/dt/metadata/mdclass AccumulationRegister"`
	Name       string         `xml:"name"`
	Synonym    EDTSynonym     `xml:"synonym"`
	Dimensions []EDTAttribute `xml:"dimensions"`
	Resources  []EDTAttribute `xml:"resources"`
	Attributes []EDTAttribute `xml:"attributes"`
}

// EDTInformationRegister структура для парсинга EDT регистра сведений
type EDTInformationRegister struct {
	XMLName    xml.Name       `xml:"http://g5.1c.ru/v8/dt/metadata/mdclass InformationRegister"`
	Name       string         `xml:"name"`
	Synonym    EDTSynonym     `xml:"synonym"`
	Dimensions []EDTAttribute `xml:"dimensions"`
	Resources  []EDTAttribute `xml:"resources"`
	Attributes []EDTAttribute `xml:"attributes"`
} // EDTSynonym синоним в EDT формате
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
	entries, err := os.ReadDir(documentsPath)
	if err != nil {
		return nil, err
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
	data, err := os.ReadFile(filePath)
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
	entries, err := os.ReadDir(catalogsPath)
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

// ParseAccumulationRegisters парсит все регистры накопления в EDT формате
func (p *EDTParser) ParseAccumulationRegisters() ([]model.MetadataObject, error) {
	regsPath := filepath.Join(p.sourcePath, "src", "AccumulationRegisters")
	if _, err := os.Stat(regsPath); os.IsNotExist(err) {
		return []model.MetadataObject{}, nil
	}
	entries, err := os.ReadDir(regsPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения каталога регистров накопления: %w", err)
	}
	var regs []model.MetadataObject
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		mdoFile := filepath.Join(regsPath, name, name+".mdo")
		if _, err := os.Stat(mdoFile); os.IsNotExist(err) {
			continue
		}
		reg, perr := p.parseAccumulationRegisterFile(mdoFile)
		if perr != nil {
			fmt.Printf("Предупреждение: ошибка парсинга регистра %s: %v\n", name, perr)
			continue
		}
		regs = append(regs, reg)
	}
	return regs, nil
}

// parseAccumulationRegisterFile парсит отдельный MDO файл регистра накопления
func (p *EDTParser) parseAccumulationRegisterFile(filePath string) (model.MetadataObject, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка чтения файла %s: %w", filePath, err)
	}
	var edtReg EDTAccumulationRegister
	if err := xml.Unmarshal(data, &edtReg); err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка парсинга XML файла %s: %w", filePath, err)
	}
	reg := model.MetadataObject{
		Type:    model.ObjectTypeAccumulationRegister,
		Name:    edtReg.Name,
		Synonym: edtReg.Synonym.Value,
	}
	// Измерения
	for _, d := range edtReg.Dimensions {
		converted := p.typeConverter.ConvertTypes(d.Type.Types)
		reg.Dimensions = append(reg.Dimensions, model.Attribute{
			Name:    d.Name,
			Synonym: d.Synonym.Value,
			Types:   converted,
		})
	}
	// Ресурсы
	for _, r := range edtReg.Resources {
		converted := p.typeConverter.ConvertTypes(r.Type.Types)
		reg.Resources = append(reg.Resources, model.Attribute{
			Name:    r.Name,
			Synonym: r.Synonym.Value,
			Types:   converted,
		})
	}
	// Реквизиты
	for _, a := range edtReg.Attributes {
		converted := p.typeConverter.ConvertTypes(a.Type.Types)
		reg.Attributes = append(reg.Attributes, model.Attribute{
			Name:    a.Name,
			Synonym: a.Synonym.Value,
			Types:   converted,
		})
	}
	return reg, nil
}

// parseCatalogFile парсит отдельный MDO файл справочника
func (p *EDTParser) parseCatalogFile(filePath string) (model.MetadataObject, error) {
	// Читаем файл
	data, err := os.ReadFile(filePath)
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
	enumsPath := filepath.Join(p.sourcePath, "src", "Enums")
	if _, err := os.Stat(enumsPath); os.IsNotExist(err) {
		return []model.MetadataObject{}, nil
	}

	entries, err := os.ReadDir(enumsPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения каталога перечислений: %w", err)
	}

	var enums []model.MetadataObject
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		mdoFile := filepath.Join(enumsPath, name, name+".mdo")
		if _, err := os.Stat(mdoFile); os.IsNotExist(err) {
			continue
		}
		data, rerr := os.ReadFile(mdoFile)
		if rerr != nil {
			fmt.Printf("Предупреждение: ошибка чтения файла перечисления %s: %v\n", mdoFile, rerr)
			continue
		}

		// Структура для парсинга EDT перечисления
		type edtEnumValue struct {
			Name    string     `xml:"name"`
			Synonym EDTSynonym `xml:"synonym"`
		}
		type edtEnum struct {
			XMLName    xml.Name       `xml:"http://g5.1c.ru/v8/dt/metadata/mdclass Enum"`
			Name       string         `xml:"name"`
			Synonym    EDTSynonym     `xml:"synonym"`
			EnumValues []edtEnumValue `xml:"enumValues"`
		}

		var ee edtEnum
		if err := xml.Unmarshal(data, &ee); err != nil {
			fmt.Printf("Предупреждение: ошибка парсинга перечисления %s: %v\n", mdoFile, err)
			continue
		}

		obj := model.MetadataObject{
			Type:    model.ObjectTypeEnum,
			Name:    ee.Name,
			Synonym: ee.Synonym.Value,
		}

		for _, v := range ee.EnumValues {
			obj.EnumValues = append(obj.EnumValues, model.EnumValue{Name: v.Name, Synonym: v.Synonym.Value})
		}

		enums = append(enums, obj)
	}

	return enums, nil
}

// ParseChartsOfCharacteristicTypes парсит планы видов характеристик (заглушка для MVP)
func (p *EDTParser) ParseChartsOfCharacteristicTypes() ([]model.MetadataObject, error) {
	chartsPath := filepath.Join(p.sourcePath, "src", "ChartsOfCharacteristicTypes")
	if _, err := os.Stat(chartsPath); os.IsNotExist(err) {
		return []model.MetadataObject{}, nil
	}

	entries, err := os.ReadDir(chartsPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения каталога планов видов характеристик: %w", err)
	}

	var charts []model.MetadataObject
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		mdoFile := filepath.Join(chartsPath, name, name+".mdo")
		if _, err := os.Stat(mdoFile); os.IsNotExist(err) {
			continue
		}
		data, rerr := os.ReadFile(mdoFile)
		if rerr != nil {
			fmt.Printf("Предупреждение: ошибка чтения файла плана %s: %v\n", mdoFile, rerr)
			continue
		}

		type EDTChart struct {
			XMLName         xml.Name            `xml:"http://g5.1c.ru/v8/dt/metadata/mdclass ChartOfCharacteristicTypes"`
			Name            string              `xml:"name"`
			Synonym         EDTSynonym          `xml:"synonym"`
			Attributes      []EDTAttribute      `xml:"attributes"`
			TabularSections []EDTTabularSection `xml:"tabularSections"`
		}

		var ec EDTChart
		if err := xml.Unmarshal(data, &ec); err != nil {
			fmt.Printf("Предупреждение: ошибка парсинга плана %s: %v\n", mdoFile, err)
			continue
		}

		obj := model.MetadataObject{
			Type:    model.ObjectTypeChartOfCharacteristicTypes,
			Name:    ec.Name,
			Synonym: ec.Synonym.Value,
		}

		for _, attr := range ec.Attributes {
			converted := p.typeConverter.ConvertTypes(attr.Type.Types)
			obj.Attributes = append(obj.Attributes, model.Attribute{
				Name:    attr.Name,
				Synonym: attr.Synonym.Value,
				Types:   converted,
			})
		}

		for _, ts := range ec.TabularSections {
			tab := model.TabularSection{
				Name:    ts.Name,
				Synonym: ts.Synonym.Value,
			}
			for _, a := range ts.Attributes {
				converted := p.typeConverter.ConvertTypes(a.Type.Types)
				tab.Attributes = append(tab.Attributes, model.Attribute{
					Name:    a.Name,
					Synonym: a.Synonym.Value,
					Types:   converted,
				})
			}
			obj.TabularSections = append(obj.TabularSections, tab)
		}

		charts = append(charts, obj)
	}

	return charts, nil
}

// ParseInformationRegisters парсит все регистры сведений в EDT формате
func (p *EDTParser) ParseInformationRegisters() ([]model.MetadataObject, error) {
	regsPath := filepath.Join(p.sourcePath, "src", "InformationRegisters")
	if _, err := os.Stat(regsPath); os.IsNotExist(err) {
		return []model.MetadataObject{}, nil
	}

	entries, err := os.ReadDir(regsPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения каталога регистров сведений: %w", err)
	}

	var regs []model.MetadataObject
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		mdoFile := filepath.Join(regsPath, name, name+".mdo")

		if _, err := os.Stat(mdoFile); os.IsNotExist(err) {
			continue
		}

		reg, perr := p.parseInformationRegisterFile(mdoFile)
		if perr != nil {
			fmt.Printf("Предупреждение: ошибка парсинга регистра %s: %v\n", name, perr)
			continue
		}

		regs = append(regs, reg)
	}

	return regs, nil
}

// ParseConstants парсит константы в EDT (MDO) структуре
func (p *EDTParser) ParseConstants() ([]model.MetadataObject, error) {
	constsPath := filepath.Join(p.sourcePath, "src", "Constants")
	if _, err := os.Stat(constsPath); os.IsNotExist(err) {
		return []model.MetadataObject{}, nil
	}

	entries, err := os.ReadDir(constsPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения каталога констант: %w", err)
	}

	var consts []model.MetadataObject
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		mdoFile := filepath.Join(constsPath, name, name+".mdo")
		if _, err := os.Stat(mdoFile); os.IsNotExist(err) {
			continue
		}
		c, perr := p.parseConstantFile(mdoFile)
		if perr != nil {
			fmt.Printf("Предупреждение: ошибка парсинга константы %s: %v\n", name, perr)
			continue
		}
		consts = append(consts, c)
	}
	return consts, nil
}

// parseConstantFile парсит MDO файл константы
func (p *EDTParser) parseConstantFile(filePath string) (model.MetadataObject, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка чтения файла %s: %w", filePath, err)
	}

	type edtConst struct {
		XMLName xml.Name   `xml:"http://g5.1c.ru/v8/dt/metadata/mdclass Constant"`
		Name    string     `xml:"name"`
		Synonym EDTSynonym `xml:"synonym"`
		Type    EDTType    `xml:"type"`
	}

	var ec edtConst
	if err := xml.Unmarshal(data, &ec); err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка парсинга XML %s: %w", filePath, err)
	}

	obj := model.MetadataObject{
		Type:    model.ObjectTypeConstant,
		Name:    ec.Name,
		Synonym: ec.Synonym.Value,
	}

	converted := p.typeConverter.ConvertTypes(ec.Type.Types)
	obj.Attributes = append(obj.Attributes, model.Attribute{
		Name:    "Значение",
		Synonym: "",
		Types:   converted,
	})

	return obj, nil
}

// ParseFilterCriteria парсит критерии отбора в EDT (MDO) структуре
func (p *EDTParser) ParseFilterCriteria() ([]model.MetadataObject, error) {
	fcPath := filepath.Join(p.sourcePath, "src", "FilterCriteria")
	if _, err := os.Stat(fcPath); os.IsNotExist(err) {
		return []model.MetadataObject{}, nil
	}

	entries, err := os.ReadDir(fcPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения каталога критериев отбора: %w", err)
	}

	var result []model.MetadataObject
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		mdoFile := filepath.Join(fcPath, name, name+".mdo")
		if _, err := os.Stat(mdoFile); os.IsNotExist(err) {
			continue
		}
		obj, perr := p.parseFilterCriteriaFile(mdoFile)
		if perr != nil {
			fmt.Printf("Предупреждение: ошибка парсинга критерия отбора %s: %v\n", name, perr)
			continue
		}
		result = append(result, obj)
	}
	return result, nil
}

// parseFilterCriteriaFile парсит MDO файл критерия отбора
func (p *EDTParser) parseFilterCriteriaFile(filePath string) (model.MetadataObject, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка чтения файла %s: %w", filePath, err)
	}

	// Структура для парсинга критерия отбора с учётом полей type и content
	type edtFilter struct {
		XMLName xml.Name   `xml:"http://g5.1c.ru/v8/dt/metadata/mdclass FilterCriterion"`
		Name    string     `xml:"name"`
		Synonym EDTSynonym `xml:"synonym"`
		Type    struct {
			Types []string `xml:"types"`
		} `xml:"type"`
		Content    []string       `xml:"content"`
		Attributes []EDTAttribute `xml:"attributes"`
	}

	var ef edtFilter
	if err := xml.Unmarshal(data, &ef); err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка парсинга XML %s: %w", filePath, err)
	}

	obj := model.MetadataObject{
		Type:    model.ObjectTypeFilterCriteria,
		Name:    ef.Name,
		Synonym: ef.Synonym.Value,
	}

	// Типы критерия
	if len(ef.Type.Types) > 0 {
		obj.FilterCriteriaTypes = p.typeConverter.ConvertTypes(ef.Type.Types)
	}

	// Состав (content)
	if len(ef.Content) > 0 {
		for _, c := range ef.Content {
			obj.FilterCriteriaContents = append(obj.FilterCriteriaContents, NormalizeFilterContentItem(c))
		}
	}

	// В EDT критерии отбора не содержат реквизитов/атрибутов, пропускаем

	return obj, nil
}

// parseInformationRegisterFile парсит отдельный MDO файл регистра сведений
func (p *EDTParser) parseInformationRegisterFile(filePath string) (model.MetadataObject, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка чтения файла %s: %w", filePath, err)
	}

	var edtReg EDTInformationRegister
	if err := xml.Unmarshal(data, &edtReg); err != nil {
		return model.MetadataObject{}, fmt.Errorf("ошибка парсинга XML файла %s: %w", filePath, err)
	}

	reg := model.MetadataObject{
		Type:    model.ObjectTypeInformationRegister,
		Name:    edtReg.Name,
		Synonym: edtReg.Synonym.Value,
	}

	// Измерения
	for _, d := range edtReg.Dimensions {
		converted := p.typeConverter.ConvertTypes(d.Type.Types)
		reg.Dimensions = append(reg.Dimensions, model.Attribute{
			Name:    d.Name,
			Synonym: d.Synonym.Value,
			Types:   converted,
		})
	}

	// Ресурсы
	for _, r := range edtReg.Resources {
		converted := p.typeConverter.ConvertTypes(r.Type.Types)
		reg.Resources = append(reg.Resources, model.Attribute{
			Name:    r.Name,
			Synonym: r.Synonym.Value,
			Types:   converted,
		})
	}

	// Реквизиты
	for _, a := range edtReg.Attributes {
		converted := p.typeConverter.ConvertTypes(a.Type.Types)
		reg.Attributes = append(reg.Attributes, model.Attribute{
			Name:    a.Name,
			Synonym: a.Synonym.Value,
			Types:   converted,
		})
	}

	return reg, nil
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
