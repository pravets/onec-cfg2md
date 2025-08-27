package generator

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"ones-cfg2md/pkg/model"
)

// CSVGenerator генератор CSV каталога объектов
type CSVGenerator struct {
	outputPath string
}

// NewCSVGenerator создает новый генератор CSV
func NewCSVGenerator(outputPath string) *CSVGenerator {
	return &CSVGenerator{
		outputPath: outputPath,
	}
}

// GenerateCatalog генерирует CSV каталог объектов
func (g *CSVGenerator) GenerateCatalog(objects []model.MetadataObject) error {
	// Создаем выходной каталог, если он не существует
	if err := os.MkdirAll(g.outputPath, 0755); err != nil {
		return fmt.Errorf("ошибка создания выходного каталога %s: %w", g.outputPath, err)
	}
	
	// Создаем CSV файл
	csvPath := filepath.Join(g.outputPath, "objects.csv")
	file, err := os.Create(csvPath)
	if err != nil {
		return fmt.Errorf("ошибка создания CSV файла %s: %w", csvPath, err)
	}
	defer file.Close()
	
	writer := csv.NewWriter(file)
	writer.Comma = ';' // Используем точку с запятой как разделитель
	defer writer.Flush()
	
	// Записываем заголовок
	header := []string{"Имя объекта", "Тип объекта", "Синоним", "Файл"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("ошибка записи заголовка CSV: %w", err)
	}
	
	// Записываем данные объектов
	for _, obj := range objects {
		entry := g.createCatalogEntry(obj)
		record := []string{
			entry.ObjectName,
			entry.ObjectType,
			entry.Synonym,
			entry.FileName,
		}
		
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("ошибка записи записи CSV для объекта %s: %w", obj.Name, err)
		}
	}
	
	return nil
}

// createCatalogEntry создает запись каталога для объекта
func (g *CSVGenerator) createCatalogEntry(obj model.MetadataObject) model.CatalogEntry {
	typeRussian := g.getObjectTypeRussian(obj.Type)
	objectName := fmt.Sprintf("%s.%s", typeRussian, obj.Name)
	fileName := fmt.Sprintf("%s_%s.md", typeRussian, obj.Name)
	
	return model.CatalogEntry{
		ObjectName: objectName,
		ObjectType: typeRussian,
		Synonym:    obj.Synonym,
		FileName:   fileName,
	}
}

// getObjectTypeRussian возвращает русское название типа объекта
func (g *CSVGenerator) getObjectTypeRussian(objType model.ObjectType) string {
	switch objType {
	case model.ObjectTypeDocument:
		return "Документ"
	case model.ObjectTypeCatalog:
		return "Справочник"
    case model.ObjectTypeAccumulationRegister:
        return "РегистрНакопления"
	case model.ObjectTypeInformationRegister:
        return "РегистрСведений"
	case model.ObjectTypeEnum:
		return "Перечисление"
	case model.ObjectTypeChartOfCharacteristicTypes:
		return "ПланВидовХарактеристик"
	default:
		return string(objType)
	}
}