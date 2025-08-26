package detector

import (
	"fmt"
	"os"
	"path/filepath"

	"ones-cfg2md/pkg/model"
)

// DetectFormat автоматически определяет формат метаданных в указанном каталоге
func DetectFormat(sourcePath string) (model.SourceFormat, error) {
	// Проверяем существование каталога
	info, err := os.Stat(sourcePath)
	if err != nil {
		return "", fmt.Errorf("не удается получить доступ к каталогу %s: %w", sourcePath, err)
	}
	
	if !info.IsDir() {
		return "", fmt.Errorf("путь %s не является каталогом", sourcePath)
	}

	// Проверяем CFG формат
	if isCFGFormat(sourcePath) {
		return model.FormatCFG, nil
	}

	// Проверяем EDT формат  
	if isEDTFormat(sourcePath) {
		return model.FormatEDT, nil
	}

	return "", fmt.Errorf("не удалось определить формат метаданных в каталоге %s", sourcePath)
}

// isCFGFormat проверяет наличие маркеров CFG формата
func isCFGFormat(sourcePath string) bool {
	configPath := filepath.Join(sourcePath, "Configuration.xml")
	_, err := os.Stat(configPath)
	return err == nil
}

// isEDTFormat проверяет наличие маркеров EDT формата
func isEDTFormat(sourcePath string) bool {
	projectPath := filepath.Join(sourcePath, ".project")
	srcPath := filepath.Join(sourcePath, "src")
	
	_, errProject := os.Stat(projectPath)
	_, errSrc := os.Stat(srcPath)
	
	return errProject == nil && errSrc == nil
}

// ValidateFormat проверяет корректность формата
func ValidateFormat(sourcePath string, format model.SourceFormat) error {
	switch format {
	case model.FormatCFG:
		if !isCFGFormat(sourcePath) {
			return fmt.Errorf("каталог %s не содержит файл Configuration.xml", sourcePath)
		}
	case model.FormatEDT:
		if !isEDTFormat(sourcePath) {
			return fmt.Errorf("каталог %s не содержит файлы .project и src/", sourcePath)
		}
	default:
		return fmt.Errorf("неподдерживаемый формат: %s", format)
	}
	
	return nil
}