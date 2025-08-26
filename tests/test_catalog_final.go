package main

import (
	"fmt"
	"os"
	"ones-cfg2md/pkg/generator"
	"ones-cfg2md/pkg/model"
	"ones-cfg2md/pkg/parser"
)

func main() {
	fmt.Println("=== Финальный тест парсинга справочников ===")
	
	// Тестируем оба формата
	testResults := make(map[string]int)
	
	// CFG формат
	fmt.Println("1. Тестирование CFG формата...")
	cfgParser, err := parser.NewParser("../examples/cfg", model.FormatCFG)
	if err != nil {
		fmt.Printf("❌ CFG: ошибка создания парсера: %v\n", err)
	} else {
		cfgCatalogs, err := cfgParser.ParseCatalogs()
		if err != nil {
			fmt.Printf("❌ CFG: ошибка парсинга: %v\n", err)
		} else {
			testResults["CFG"] = len(cfgCatalogs)
			fmt.Printf("✅ CFG: найдено справочников: %d\n", len(cfgCatalogs))
			
			// Показываем детали первого справочника
			if len(cfgCatalogs) > 0 {
				cat := cfgCatalogs[0]
				fmt.Printf("   Пример: %s (%s) - %d атрибутов, %d ТЧ\n", 
					cat.Name, cat.Synonym, len(cat.Attributes), len(cat.TabularSections))
			}
		}
	}
	
	// EDT формат  
	fmt.Println("2. Тестирование EDT формата...")
	edtParser, err := parser.NewParser("../examples/edt", model.FormatEDT)
	if err != nil {
		fmt.Printf("❌ EDT: ошибка создания парсера: %v\n", err)
	} else {
		edtCatalogs, err := edtParser.ParseCatalogs()
		if err != nil {
			fmt.Printf("❌ EDT: ошибка парсинга: %v\n", err)
		} else {
			testResults["EDT"] = len(edtCatalogs)
			fmt.Printf("✅ EDT: найдено справочников: %d\n", len(edtCatalogs))
			
			// Показываем детали первого справочника
			if len(edtCatalogs) > 0 {
				cat := edtCatalogs[0]
				fmt.Printf("   Пример: %s (%s) - %d атрибутов, %d ТЧ\n", 
					cat.Name, cat.Synonym, len(cat.Attributes), len(cat.TabularSections))
			}
		}
	}
	
	// Тестируем интеграцию с генераторами
	fmt.Println("3. Тестирование генерации файлов...")
	if cfgParser != nil {
		catalogs, err := cfgParser.ParseCatalogs()
		if err == nil && len(catalogs) > 0 {
			// Создаем временную директорию
			outputDir := "../test_final_output"
			os.MkdirAll(outputDir, 0755)
			
			// Генерируем markdown
			mdGen := generator.NewMarkdownGenerator(outputDir)
			if err := mdGen.GenerateFiles(catalogs); err == nil {
				fmt.Println("✅ Markdown файлы созданы")
			} else {
				fmt.Printf("❌ Ошибка генерации Markdown: %v\n", err)
			}
			
			// Генерируем CSV
			csvGen := generator.NewCSVGenerator(outputDir)
			if err := csvGen.GenerateCatalog(catalogs); err == nil {
				fmt.Println("✅ CSV каталог создан")
			} else {
				fmt.Printf("❌ Ошибка генерации CSV: %v\n", err)
			}
			
			// Показываем созданные файлы
			if entries, err := os.ReadDir(outputDir); err == nil {
				fmt.Printf("   Создано файлов: %d\n", len(entries))
			}
		}
	}
	
	// Итоговый отчет
	fmt.Println("\n=== Результаты тестирования ===")
	totalCatalogs := 0
	for format, count := range testResults {
		fmt.Printf("%s формат: %d справочников\n", format, count)
		totalCatalogs += count
	}
	fmt.Printf("Всего обработано: %d справочников\n", totalCatalogs)
	
	if totalCatalogs > 0 {
		fmt.Println("🎉 Тест успешно завершен!")
	} else {
		fmt.Println("⚠️  Справочники не найдены или произошли ошибки")
	}
}