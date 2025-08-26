package main

import (
	"fmt"
	"log"
	"ones-cfg2md/pkg/parser"
)

func main() {
	fmt.Println("=== Тест парсера справочников ===")
	
	// Создаем парсер для CFG формата
	cfgParser, err := parser.NewCFGParser("../examples/cfg")
	if err != nil {
		log.Fatal("Ошибка создания парсера:", err)
	}

	// Парсим справочники
	catalogs, err := cfgParser.ParseCatalogs()
	if err != nil {
		log.Fatal("Ошибка парсинга справочников:", err)
	}

	fmt.Printf("✓ Найдено справочников: %d\n", len(catalogs))
	
	for _, catalog := range catalogs {
		fmt.Printf("\n=== Справочник: %s ===\n", catalog.Name)
		fmt.Printf("Синоним: %s\n", catalog.Synonym)
		fmt.Printf("Тип: %s\n", catalog.Type)
		fmt.Printf("Атрибутов: %d\n", len(catalog.Attributes))
		fmt.Printf("Табличных частей: %d\n", len(catalog.TabularSections))
		
		// Выводим первые 5 атрибутов
		if len(catalog.Attributes) > 0 {
			fmt.Println("Первые атрибуты:")
			for i, attr := range catalog.Attributes {
				if i >= 5 { // ограничиваем вывод
					fmt.Println("  ... (и другие)")
					break
				}
				fmt.Printf("  - %s (%s): %v\n", attr.Name, attr.Synonym, attr.Types)
			}
		}
		
		// Выводим табличные части
		if len(catalog.TabularSections) > 0 {
			fmt.Println("Табличные части:")
			for _, ts := range catalog.TabularSections {
				fmt.Printf("  - %s (%s): %d атрибутов\n", ts.Name, ts.Synonym, len(ts.Attributes))
				
				// Выводим первые атрибуты табличной части
				for i, attr := range ts.Attributes {
					if i >= 3 { // ограничиваем вывод
						fmt.Println("      ... (и другие)")
						break
					}
					fmt.Printf("      - %s (%s): %v\n", attr.Name, attr.Synonym, attr.Types)
				}
			}
		}
	}
	
	// Дополнительная проверка: тестируем парсинг через общий интерфейс
	fmt.Println("\n=== Тест через общий интерфейс ===")
	
	generalParser, err := parser.NewParser("../examples/cfg", "cfg")
	if err != nil {
		log.Fatal("Ошибка создания общего парсера:", err)
	}
	
	allCatalogs, err := generalParser.ParseCatalogs()
	if err != nil {
		log.Fatal("Ошибка парсинга через общий интерфейс:", err)
	}
	
	fmt.Printf("✓ Через общий интерфейс найдено справочников: %d\n", len(allCatalogs))
	
	// Проверяем что результаты совпадают
	if len(catalogs) == len(allCatalogs) {
		fmt.Println("✓ Результаты парсинга совпадают")
	} else {
		fmt.Printf("⚠ Несоответствие: прямой парсер нашел %d, общий - %d\n", len(catalogs), len(allCatalogs))
	}
	
	fmt.Println("=== Тест завершен ===")
}