package cmd

import (
	"fmt"
	"os"
	"strings"

	"ones-cfg2md/pkg/detector"
	"ones-cfg2md/pkg/generator"
	"ones-cfg2md/pkg/model"
	"ones-cfg2md/pkg/parser"

	"github.com/spf13/cobra"
)

var (
	// Флаги командной строки
	formatFlag   string
	typesFlag    string
	verboseFlag  bool
)

// rootCmd основная команда
var rootCmd = &cobra.Command{
	Use:   "ones-cfg2md <source_directory> <output_directory>",
	Short: "Конвертер метаданных 1С в Markdown",
	Long: `Программа конвертирует метаданные конфигурации 1С из CFG или EDT формата 
в документацию Markdown для использования в Model Context Protocol (MCP).

Поддерживаемые форматы:
  - CFG (Конфигуратор): маркер Configuration.xml в корне
  - EDT (Eclipse Development Tools): маркеры .project и src/ в корне

Поддерживаемые типы объектов:
  - documents (документы)
  - catalogs (справочники) 
  - accumulationregisters (регистры накопления)
  - informationregisters (регистры сведений)
  - enums (перечисления)
  - chartsofcharacteristictypes (планы видов характеристик)`,
	Args: cobra.ExactArgs(2),
	Run:  runConversion,
}

// Execute выполняет основную команду
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Настройка флагов
	rootCmd.Flags().StringVar(&formatFlag, "format", "", 
		"Принудительное указание формата (cfg/edt), по умолчанию автоопределение")
	
    rootCmd.Flags().StringVar(&typesFlag, "types", "documents,catalogs,accumulationregisters,informationregisters,enums,chartsofcharacteristictypes", 
        "Типы объектов для обработки, разделенные запятыми (documents,catalogs,accumulationregisters,informationregisters,enums,chartsofcharacteristictypes)")
	
	rootCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, 
		"Подробный вывод процесса обработки")
}

// runConversion выполняет конвертацию
func runConversion(cmd *cobra.Command, args []string) {
	sourcePath := args[0]
	outputPath := args[1]
	
	// Создаем опции конвертации
	options := model.ConversionOptions{
		SourcePath:  sourcePath,
		OutputPath:  outputPath,
		Verbose:     verboseFlag,
	}
	
	// Определяем формат
	var err error
	if formatFlag != "" {
		// Используем указанный формат
		switch formatFlag {
		case "cfg":
			options.Format = model.FormatCFG
		case "edt":
			options.Format = model.FormatEDT
		default:
			fmt.Fprintf(os.Stderr, "Ошибка: неподдерживаемый формат '%s'. Используйте 'cfg' или 'edt'\n", formatFlag)
			os.Exit(1)
		}
		
		// Проверяем корректность формата
		if err := detector.ValidateFormat(sourcePath, options.Format); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Автоопределение формата
		options.Format, err = detector.DetectFormat(sourcePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка определения формата: %v\n", err)
			os.Exit(1)
		}
	}
	
	if verboseFlag {
		fmt.Printf("Определен формат: %s\n", options.Format)
	}
	
	// Парсим типы объектов
	options.ObjectTypes, err = parseObjectTypes(typesFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
	
	if verboseFlag {
		fmt.Printf("Типы объектов для обработки: %v\n", options.ObjectTypes)
	}
	
	// Выполняем конвертацию
	if err := performConversion(options); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка конвертации: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Конвертация завершена успешно!\n")
	fmt.Printf("Результаты сохранены в: %s\n", outputPath)
}

// parseObjectTypes парсит строку типов объектов
func parseObjectTypes(typesStr string) ([]model.ObjectType, error) {
	if typesStr == "" {
		return []model.ObjectType{model.ObjectTypeDocument}, nil
	}
	
	typeNames := strings.Split(typesStr, ",")
	var objectTypes []model.ObjectType
	
	for _, typeName := range typeNames {
		typeName = strings.TrimSpace(typeName)
		switch typeName {
		case "documents":
			objectTypes = append(objectTypes, model.ObjectTypeDocument)
		case "catalogs":
			objectTypes = append(objectTypes, model.ObjectTypeCatalog)
        case "accumulationregisters":
            objectTypes = append(objectTypes, model.ObjectTypeAccumulationRegister)
		case "informationregisters":
            objectTypes = append(objectTypes, model.ObjectTypeInformationRegister)
		case "enums":
			objectTypes = append(objectTypes, model.ObjectTypeEnum)
		case "chartsofcharacteristictypes":
			objectTypes = append(objectTypes, model.ObjectTypeChartOfCharacteristicTypes)
		default:
			return nil, fmt.Errorf("неподдерживаемый тип объекта: %s", typeName)
		}
	}
	
	return objectTypes, nil
}

// performConversion выполняет конвертацию
func performConversion(options model.ConversionOptions) error {
	// Создаем парсер
	metadataParser, err := parser.NewParser(options.SourcePath, options.Format)
	if err != nil {
		return fmt.Errorf("ошибка создания парсера: %w", err)
	}
	
	if options.Verbose {
		fmt.Printf("Начинаем парсинг метаданных...\n")
	}
	
	// Парсим объекты
	objects, err := metadataParser.ParseObjectsByType(options.ObjectTypes)
	if err != nil {
		return fmt.Errorf("ошибка парсинга метаданных: %w", err)
	}
	
	if options.Verbose {
		fmt.Printf("Найдено объектов: %d\n", len(objects))
	}
	
	if len(objects) == 0 {
		fmt.Printf("Объекты указанных типов не найдены\n")
		return nil
	}
	
	// Генерируем Markdown файлы
	if options.Verbose {
		fmt.Printf("Генерируем Markdown файлы...\n")
	}
	
	markdownGen := generator.NewMarkdownGenerator(options.OutputPath)
	if err := markdownGen.GenerateFiles(objects); err != nil {
		return fmt.Errorf("ошибка генерации Markdown файлов: %w", err)
	}
	
	// Генерируем CSV каталог
	if options.Verbose {
		fmt.Printf("Генерируем CSV каталог...\n")
	}
	
	csvGen := generator.NewCSVGenerator(options.OutputPath)
	if err := csvGen.GenerateCatalog(objects); err != nil {
		return fmt.Errorf("ошибка генерации CSV каталога: %w", err)
	}
	
	return nil
}