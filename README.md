# onec-cfg2md

![CodeRabbit Pull Request Reviews](https://img.shields.io/coderabbit/prs/github/pravets/onec-cfg2md?utm_source=oss&utm_medium=github&utm_campaign=pravets%2Fonec-cfg2md&labelColor=171717&color=FF570A&link=https%3A%2F%2Fcoderabbit.ai&label=CodeRabbit+Reviews)
![License](https://img.shields.io/github/license/pravets/onec-cfg2md)
[![Telegram](https://telegram-badge.vercel.app/api/telegram-badge?channelId=@pravets_IT)](https://t.me/pravets_it)

Конвертер метаданных конфигурации 1С в формат Markdown для использования в Model Context Protocol (MCP).

## Описание

Программа автоматически определяет формат метаданных 1С (CFG или EDT) и конвертирует их в удобную для чтения документацию в формате Markdown. Также создается CSV каталог всех объектов.

## Поддерживаемые форматы

- **CFG формат** (Конфигуратор) - файлы XML с Configuration.xml в корне
- **EDT формат** (Eclipse Development Tools) - файлы MDO с .project и src/ в корне

## Поддерживаемые типы метаданных

Внизу приведена таблица поддерживаемых типов метаданных: русское имя, англоязычный идентификатор (используется в коде) и ключ, который принимает опция `--types`.

| Русское имя | Английский идентификатор | CLI-ключ (`--types`) |
|---|---:|---:|
| Документ | `Document` | `documents` |
| Справочник | `Catalog` | `catalogs` |
| Перечисление | `Enum` | `enums` |
| План видов характеристик | `ChartOfCharacteristicTypes` | `chartsofcharacteristictypes` |
| Регистр накопления | `AccumulationRegister` | `accumulationregisters` |
| Регистр сведений | `InformationRegister` | `informationregisters` |
| Константа | `Constant` | `constants` |
| Критерий отбора | `FilterCriteria` | `filtercriterias` |

Опция `--types` принимает перечисление ключей через запятую. Пример валидного значения:

```
documents,catalogs,accumulationregisters,informationregisters,enums,chartsofcharacteristictypes,constants,filtercriterias
```

Шаблон имени Markdown-файла: `Тип_Имя.md`, где `Тип` — русское название типа (например, `Документ`, `Справочник`), а `Имя` — системное имя объекта.


## Установка

```bash
# Клонируем репозиторий
git clone <repository-url>
cd onec-cfg2md

# Устанавливаем зависимости
go mod tidy

# Собираем программу
go build -o bin/onec-cfg2md .
```

## Использование

### Основные команды

```bash
# Автоопределение формата и обработка всех документов
./bin/onec-cfg2md ./src ./output

# Принудительное указание формата
./bin/onec-cfg2md --format=edt ./edt-project ./docs

# Обработка только документов в подробном режиме
./bin/onec-cfg2md --format=cfg --types=documents --verbose ./cfg ./docs
```

### Параметры

- `--format` - принудительное указание формата (cfg/edt)
- `--types` - типы объектов для обработки (documents,catalogs,enums,charts)
- `--verbose` - подробный вывод процесса обработки

### Примеры

```bash
# Обработка примера EDT формата
onec-cfg2md ./fixtures/input/edt ./result/edt

# Обработка примера CFG формата  
onec-cfg2md ./fixtures/input/cfg ./result/cfg

# Только документы с подробным выводом
onec-cfg2md --types=documents --verbose ./fixtures/input/edt ./docs
```

## Структура выходных файлов

### Markdown файлы

Для каждого объекта создается файл вида `Тип_Имя.md`:

```markdown
# Документ: АвансовыйОтчет (Авансовый отчет)

## Реквизиты шапки

- Организация (Справочник.Организации)
- ПодотчетноеЛицо (Справочник.ФизическиеЛица)
- Валюта (Справочник.Валюты)

## Табличные части

### ПрочиеРасходы (Расходы)

- Сумма (Число)
- Комментарий (Строка)
```

### CSV каталог

Файл `objects.csv` содержит сводную информацию:

```csv
Имя объекта;Тип объекта;Синоним;Файл
Документ.АвансовыйОтчет;Документ;Авансовый отчет;Документ_АвансовыйОтчет.md
```

## Разработка

### Сборка

```bash
make build       # Основная программа
```

### Тестирование

```bash
make test           # Unit тесты
make run-test-edt   # Тест EDT формата
make run-test-cfg   # Тест CFG формата
```

### Структура проекта

```
├── main.go              # точка входа (вызов CLI)
├── cmd/                 # реализация CLI (cobra-команды)
├── pkg/
│   ├── detector/        # определение формата (CFG/EDT)
│   ├── generator/       # генераторы Markdown и CSV
│   ├── model/           # модель метаданных (MetadataObject и пр.)
│   ├── parser/          # парсеры CFG и EDT (cfg_parser.go, edt_parser.go) 
│   └── testutil/        # вспомогательные модули для тестов
├── fixtures/            # тестовые фикстуры (input/ и output/)
└── docs/                # техническая документация
```

## Техническое задание

Подробное техническое задание доступно в файле [docs/TECHNICAL_SPECIFICATION.md](docs/TECHNICAL_SPECIFICATION.md).

## Лицензия

[LICENSE](LICENSE)
