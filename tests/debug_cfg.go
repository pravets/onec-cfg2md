package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
)

// Простая структура для тестирования парсинга CFG
type TestCFGDocument struct {
	XMLName  xml.Name         `xml:"MetaDataObject"`
	Document TestCFGContent   `xml:"Document"`
}

type TestCFGContent struct {
	Properties TestCFGProperties `xml:"Properties"`
}

type TestCFGProperties struct {
	Name    string         `xml:"Name"`
	Synonym TestCFGSynonym `xml:"Synonym"`
}

type TestCFGSynonym struct {
	Items []TestCFGSynonymItem `xml:"item"`
}

type TestCFGSynonymItem struct {
	Lang    string `xml:"lang"`
	Content string `xml:"content"`
}

func main() {
	fmt.Println("=== Отладка CFG парсера ===")
	
	// Читаем файл
	data, err := ioutil.ReadFile("../examples/cfg/Documents/АвансовыйОтчет.xml")
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}
	
	fmt.Printf("Размер файла: %d байт\n", len(data))
	
	// Пробуем парсить
	var doc TestCFGDocument
	if err := xml.Unmarshal(data, &doc); err != nil {
		log.Fatalf("Ошибка парсинга XML: %v", err)
	}
	
	fmt.Printf("Имя документа: %s\n", doc.Document.Properties.Name)
	fmt.Printf("Количество элементов синонима: %d\n", len(doc.Document.Properties.Synonym.Items))
	
	for i, item := range doc.Document.Properties.Synonym.Items {
		fmt.Printf("Синоним %d: lang=%s, content=%s\n", i, item.Lang, item.Content)
	}
	
	fmt.Println("=== Парсинг завершен ===")
}