package parser

import (
    "fmt"
    "ones-cfg2md/pkg/model"
    "ones-cfg2md/pkg/detector"
)

// MetadataParser интерфейс для парсинга метаданных
type MetadataParser interface {
	// ParseDocuments парсит все документы в конфигурации
	ParseDocuments() ([]model.MetadataObject, error)
	
	// ParseCatalogs парсит все справочники в конфигурации
	ParseCatalogs() ([]model.MetadataObject, error)
	
	// ParseEnums парсит все перечисления в конфигурации
	ParseEnums() ([]model.MetadataObject, error)
	
	// ParseChartsOfCharacteristicTypes парсит все планы видов характеристик
	ParseChartsOfCharacteristicTypes() ([]model.MetadataObject, error)
	
	// ParseObjectsByType парсит объекты указанных типов
	ParseObjectsByType(objectTypes []model.ObjectType) ([]model.MetadataObject, error)
}

// NewParser создает парсер для указанного формата
func NewParser(sourcePath string, format model.SourceFormat) (MetadataParser, error) {
	switch format {
	case model.FormatEDT:
		return NewEDTParser(sourcePath)
	case model.FormatCFG:
		return NewCFGParser(sourcePath)
	default:
        // Пытаемся автоопределить формат, если пришло пустое значение
        if string(format) == "" {
            // Ленивая загрузка, чтобы не тянуть detector циклом
            // Вызывается только если формат пустой
            detected, err := detector.DetectFormat(sourcePath)
            if err != nil {
                return nil, fmt.Errorf("не удалось определить формат: %w", err)
            }
            return NewParser(sourcePath, detected)
        }
        return nil, fmt.Errorf("неподдерживаемый формат: %s", format)
	}
}

// TypeConverter интерфейс для преобразования типов
type TypeConverter interface {
	// ConvertType преобразует тип из метаданных в читаемый формат
	ConvertType(metadataType string) string
	
	// ConvertTypes преобразует массив типов
	ConvertTypes(metadataTypes []string) []string
}