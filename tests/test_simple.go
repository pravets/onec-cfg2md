package main

import (
	"fmt"
	"log"
	"ones-cfg2md/pkg/detector"
)

func main() {
	fmt.Println("=== Простой тест детектора ===")
	
	// Тестируем детектор на примерах
	testPaths := []string{
		"../examples/edt",
		"../examples/cfg",
	}
	
	for _, path := range testPaths {
		fmt.Printf("Тестируем: %s\n", path)
		format, err := detector.DetectFormat(path)
		if err != nil {
			log.Printf("Ошибка: %v", err)
			continue
		}
		fmt.Printf("  Формат: %s\n", format)
	}
	
	fmt.Println("=== Тест завершен ===")
}