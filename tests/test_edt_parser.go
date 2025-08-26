package main

import (
	"fmt"
	"log"

	"ones-cfg2md/pkg/detector"
	"ones-cfg2md/pkg/generator"
	"ones-cfg2md/pkg/parser"
)

func main() {
	// Тестируем на примере EDT формата
	sourcePath := "../examples/edt"
	outputPath := "../test_output"

	fmt.Println("=== Тест конвертера метаданных 1С ===")

	// 1. Тестируем детектор формата
	fmt.Println("1. Определяем формат...")
	format, err := detector.DetectFormat(sourcePath)
	if err != nil {
		log.Fatalf("Ошибка определения формата: %v", err)
	}
	fmt.Printf("   Определен формат: %s\n", format)

	// 2. Создаем парсер
	fmt.Println("2. Создаем парсер...")
	metadataParser, err := parser.NewParser(sourcePath, format)
	if err != nil {
		log.Fatalf("Ошибка создания парсера: %v", err)
	}
	fmt.Println("   Парсер создан успешно")

	// 3. Парсим документы
	fmt.Println("3. Парсим документы...")
	documents, err := metadataParser.ParseDocuments()
	if err != nil {
		log.Fatalf("Ошибка парсинга документов: %v", err)
	}
	fmt.Printf("   Найдено документов: %d\n", len(documents))

	// Выводим информацию о найденных документах
	for i, doc := range documents {
		fmt.Printf("   Документ %d: %s (%s)\n", i+1, doc.Name, doc.Synonym)
		fmt.Printf("     Атрибутов: %d\n", len(doc.Attributes))
		fmt.Printf("     Табличных частей: %d\n", len(doc.TabularSections))
	}

	// 4. Генерируем Markdown
	fmt.Println("4. Генерируем Markdown файлы...")
	markdownGen := generator.NewMarkdownGenerator(outputPath)
	if err := markdownGen.GenerateFiles(documents); err != nil {
		log.Fatalf("Ошибка генерации Markdown: %v", err)
	}
	fmt.Println("   Markdown файлы созданы")

	// 5. Генерируем CSV каталог
	fmt.Println("5. Генерируем CSV каталог...")
	csvGen := generator.NewCSVGenerator(outputPath)
	if err := csvGen.GenerateCatalog(documents); err != nil {
		log.Fatalf("Ошибка генерации CSV: %v", err)
	}
	fmt.Println("   CSV каталог создан")

	fmt.Println("=== Тест завершен успешно! ===")
	fmt.Printf("Результаты сохранены в: %s\n", outputPath)
}