package parser

import "strings"

// NormalizeFilterContentItem преобразует элементы состава критерия отбора в читабельную русскую форму
func NormalizeFilterContentItem(item string) string {
	if strings.TrimSpace(item) == "" {
		return item
	}
	parts := strings.Split(item, ".")
	for i, tok := range parts {
		switch tok {
		case "Document":
			parts[i] = "Документ"
		case "Attribute":
			parts[i] = "Реквизит"
		case "Catalog":
			parts[i] = "Справочник"
		case "Enum":
			parts[i] = "Перечисление"
		case "ChartOfCharacteristicTypes":
			parts[i] = "ПланВидовХарактеристик"
		case "Constant":
			parts[i] = "Константа"
		}
	}
	return strings.Join(parts, ".")
}
