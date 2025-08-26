package main

import (
	"fmt"
	"ones-cfg2md/pkg/parser"
)

func main() {
	fmt.Println("=== Тест преобразователя типов ===")
	
	converter := parser.NewTypeConverter()
	
	testTypes := []string{
		"xs:string",
		"xs:boolean", 
		"CatalogRef.Организации",
		"EnumRef.СтатусыДокумента",
		"DefinedType.ДенежнаяСуммаЛюбогоЗнака",
		"DocumentRef.ЗаказКлиента",
		"",
		"SomeUnknownType",
	}
	
	for _, testType := range testTypes {
		converted := converter.ConvertType(testType)
		fmt.Printf("'%s' -> '%s'\n", testType, converted)
	}
	
	fmt.Println("=== Тест завершен ===")
}