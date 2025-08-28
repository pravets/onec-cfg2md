package parser

import (
	"regexp"
	"strings"
)

// DefaultTypeConverter реализация преобразователя типов
type DefaultTypeConverter struct{}

// NewTypeConverter создает новый преобразователь типов
func NewTypeConverter() TypeConverter {
	return &DefaultTypeConverter{}
}

// ConvertType преобразует тип из метаданных 1С в читаемый формат
func (c *DefaultTypeConverter) ConvertType(metadataType string) string {
	// Убираем лишние пробелы
	metadataType = strings.TrimSpace(metadataType)

	// Если тип пустой, возвращаем пустую строку
	if metadataType == "" {
		return ""
	}

	// Регулярные выражения для различных типов
	patterns := map[string]string{
		`^CatalogRef\.(.+)$`:                    "Справочник.$1",
		`^DocumentRef\.(.+)$`:                   "Документ.$1",
		`^EnumRef\.(.+)$`:                       "Перечисление.$1",
		`^ChartOfCharacteristicTypesRef\.(.+)$`: "ПланВидовХарактеристик.$1",
		`^DefinedType\.(.+)$`:                   "ОпределяемыйТип.$1",
		`^String$`:                              "Строка",
		`^Boolean$`:                             "Булево",
		`^Date$`:                                "Дата",
		`^Number$`:                              "Число",
		`^Characteristic\.(.+)$`:                "Характеристика.$1",
		// XML Schema типы (часто встречаются в CFG формате)
		`^xs:string$`:   "Строка",
		`^xs:boolean$`:  "Булево",
		`^xs:date$`:     "Дата",
		`^xs:dateTime$`: "ДатаВремя",
		`^xs:decimal$`:  "Число",
		`^xs:double$`:   "Число",
		`^xs:int$`:      "Число",
		`^xs:integer$`:  "Число",
		// Типы 1С без префиксов
		`^Type$`:         "Тип",
		`^ValueStorage$`: "ХранилищеЗначения",
		`^UUID$`:         "УникальныйИдентификатор",
		`^AnyRef$`:       "ЛюбаяСсылка",
	}

	// Применяем паттерны
	for pattern, replacement := range patterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(metadataType) {
			return re.ReplaceAllString(metadataType, replacement)
		}
	}

	// Если тип не распознан, возвращаем как есть
	return metadataType
}

// ConvertTypes преобразует массив типов
func (c *DefaultTypeConverter) ConvertTypes(metadataTypes []string) []string {
	converted := make([]string, len(metadataTypes))
	for i, metadataType := range metadataTypes {
		converted[i] = c.ConvertType(metadataType)
	}
	return converted
}
