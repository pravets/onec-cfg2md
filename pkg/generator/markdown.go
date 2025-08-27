package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ones-cfg2md/pkg/model"
)

// MarkdownGenerator генератор Markdown файлов
type MarkdownGenerator struct {
	outputPath string
}

// NewMarkdownGenerator создает новый генератор Markdown
func NewMarkdownGenerator(outputPath string) *MarkdownGenerator {
	return &MarkdownGenerator{
		outputPath: outputPath,
	}
}

// GenerateFiles генерирует Markdown файлы для всех объектов
func (g *MarkdownGenerator) GenerateFiles(objects []model.MetadataObject) error {
	// Создаем выходной каталог, если он не существует
	if err := os.MkdirAll(g.outputPath, 0755); err != nil {
		return fmt.Errorf("ошибка создания выходного каталога %s: %w", g.outputPath, err)
	}
	
	for _, obj := range objects {
		if err := g.generateFile(obj); err != nil {
			return fmt.Errorf("ошибка генерации файла для объекта %s: %w", obj.Name, err)
		}
	}
	
	return nil
}

// generateFile генерирует Markdown файл для одного объекта
func (g *MarkdownGenerator) generateFile(obj model.MetadataObject) error {
	// Формируем имя файла
	fileName := g.getFileName(obj)
	filePath := filepath.Join(g.outputPath, fileName)
	
	// Генерируем содержимое
	content := g.generateContent(obj)
	
	// Записываем файл
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("ошибка записи файла %s: %w", filePath, err)
	}
	
	return nil
}

// getFileName формирует имя файла для объекта
func (g *MarkdownGenerator) getFileName(obj model.MetadataObject) string {
	return fmt.Sprintf("%s_%s.md", g.getObjectTypeRussian(obj.Type), obj.Name)
}

// getObjectTypeRussian возвращает русское название типа объекта
func (g *MarkdownGenerator) getObjectTypeRussian(objType model.ObjectType) string {
	switch objType {
	case model.ObjectTypeDocument:
		return "Документ"
	case model.ObjectTypeCatalog:
		return "Справочник"
    case model.ObjectTypeAccumulationRegister:
        return "Регистр накопления"
	case model.ObjectTypeInformationRegister:
        return "Регистр сведений"
	case model.ObjectTypeEnum:
		return "Перечисление"
	case model.ObjectTypeChartOfCharacteristicTypes:
		return "ПланВидовХарактеристик"
	default:
		return string(objType)
	}
}

// generateContent генерирует содержимое Markdown файла
func (g *MarkdownGenerator) generateContent(obj model.MetadataObject) string {
	var content strings.Builder
	
	// Заголовок
	typeRussian := g.getObjectTypeRussian(obj.Type)
	content.WriteString(fmt.Sprintf("# %s: %s", typeRussian, obj.Name))
	if obj.Synonym != "" {
		content.WriteString(fmt.Sprintf(" (%s)", obj.Synonym))
	}
	content.WriteString("\n\n")
	
    // Для перечислений: печать значений
    if obj.Type == model.ObjectTypeEnum {
        if len(obj.EnumValues) > 0 {
            content.WriteString("## Значения\n\n")
            for _, v := range obj.EnumValues {
                if v.Synonym != "" {
                    content.WriteString(fmt.Sprintf("- %s (%s)\n", v.Name, v.Synonym))
                } else {
                    content.WriteString(fmt.Sprintf("- %s\n", v.Name))
                }
            }
            content.WriteString("\n")
        }
        return content.String()
    }

    // Для регистров накопления: Измерения, Ресурсы, Реквизиты
    if obj.Type == model.ObjectTypeAccumulationRegister {
        if len(obj.Dimensions) > 0 {
            content.WriteString("## Измерения\n\n")
            for _, d := range obj.Dimensions {
                typesStr := strings.Join(d.Types, ", ")
                content.WriteString(fmt.Sprintf("- %s (%s)\n", d.Name, typesStr))
            }
            content.WriteString("\n")
        }
        if len(obj.Resources) > 0 {
            content.WriteString("## Ресурсы\n\n")
            for _, r := range obj.Resources {
                typesStr := strings.Join(r.Types, ", ")
                content.WriteString(fmt.Sprintf("- %s (%s)\n", r.Name, typesStr))
            }
            content.WriteString("\n")
        }
        // Реквизиты регистра
        if len(obj.Attributes) > 0 {
            content.WriteString("## Реквизиты\n\n")
            for _, a := range obj.Attributes {
                typesStr := strings.Join(a.Types, ", ")
                content.WriteString(fmt.Sprintf("- %s (%s)\n", a.Name, typesStr))
            }
            content.WriteString("\n")
        }
        return content.String()
    }

    // Для регистров измерений: Измерения, Ресурсы, Реквизиты
    if obj.Type == model.ObjectTypeInformationRegister {
        if len(obj.Dimensions) > 0 {
            content.WriteString("## Измерения\n\n")
            for _, d := range obj.Dimensions {
                typesStr := strings.Join(d.Types, ", ")
                content.WriteString(fmt.Sprintf("- %s (%s)\n", d.Name, typesStr))
            }
            content.WriteString("\n")
        }
        if len(obj.Resources) > 0 {
            content.WriteString("## Ресурсы\n\n")
            for _, r := range obj.Resources {
                typesStr := strings.Join(r.Types, ", ")
                content.WriteString(fmt.Sprintf("- %s (%s)\n", r.Name, typesStr))
            }
            content.WriteString("\n")
        }
        // Реквизиты регистра
        if len(obj.Attributes) > 0 {
            content.WriteString("## Реквизиты\n\n")
            for _, a := range obj.Attributes {
                typesStr := strings.Join(a.Types, ", ")
                content.WriteString(fmt.Sprintf("- %s (%s)\n", a.Name, typesStr))
            }
            content.WriteString("\n")
        }
        return content.String()
    }

    // Реквизиты / Реквизиты шапки
    if len(obj.Attributes) > 0 {
        if obj.Type == model.ObjectTypeCatalog {
            content.WriteString("## Реквизиты\n\n")
        } else {
            content.WriteString("## Реквизиты шапки\n\n")
        }
        for _, attr := range obj.Attributes {
            typesStr := strings.Join(attr.Types, ", ")
            content.WriteString(fmt.Sprintf("- %s (%s)\n", attr.Name, typesStr))
        }
        content.WriteString("\n")
    }
	
	// Табличные части
	if len(obj.TabularSections) > 0 {
		content.WriteString("## Табличные части\n\n")
		for _, ts := range obj.TabularSections {
			// Заголовок табличной части
			content.WriteString(fmt.Sprintf("### %s", ts.Name))
			if ts.Synonym != "" {
				content.WriteString(fmt.Sprintf(" (%s)", ts.Synonym))
			}
			content.WriteString("\n\n")
			
			// Атрибуты табличной части
			for _, attr := range ts.Attributes {
				typesStr := strings.Join(attr.Types, ", ")
				content.WriteString(fmt.Sprintf("- %s (%s)\n", attr.Name, typesStr))
			}
			content.WriteString("\n")
		}
	}
	
	return content.String()
}