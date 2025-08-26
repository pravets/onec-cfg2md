.PHONY: build test clean run-test-edt run-test-cfg

# Сборка основной программы
build:
	go build -o bin/ones-cfg2md .

# Сборка тестовых программ
build-test:
	go build -o bin/test-simple tests/test_simple.go
	go build -o bin/test-edt tests/test_edt_parser.go

# Запуск тестов
test:
	go test ./...

# Очистка
clean:
	rm -rf bin/ test_output/

# Тест детектора формата
run-test-simple: build-test
	./bin/test-simple

# Тест EDT формата
run-test-edt: build-test
	./bin/test-edt

# Тест CFG формата
run-test-cfg: build
	./bin/ones-cfg2md --format=cfg --verbose ./examples/cfg ./test_output_cfg

# Тест EDT формата через основную программу
run-edt: build
	./bin/ones-cfg2md --format=edt --verbose ./examples/edt ./test_output_edt

# Установка зависимостей
deps:
	go mod tidy
	go mod download

# Инициализация
init: deps
	mkdir -p bin