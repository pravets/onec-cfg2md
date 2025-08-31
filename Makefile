.PHONY: build test clean run-test-edt run-test-cfg run-edt run-cfg init deps

# Сборка основной программы
build:
	go build -o bin/onec-cfg2md .

# Запуск всех unit-тестов
test:
	go test ./...

# Очистка артефактов
clean:
	rm -rf bin/ result/ test_output/

# Запуск тестов парсера EDT (локальные тесты пакета)
run-test-edt:
	go test ./pkg/parser -run TestEDT -v

# Запуск тестов парсера CFG (локальные тесты пакета)
run-test-cfg:
	go test ./pkg/parser -run TestCFG -v

# Пример запуска основного бинаря для CFG/EDT на фикстурах
run-cfg: build
	./bin/onec-cfg2md --format=cfg --verbose ./fixtures/input/cfg ./result/cfg

run-edt: build
	./bin/onec-cfg2md --format=edt --verbose ./fixtures/input/edt ./result/edt

# Установка зависимостей
deps:
	go mod tidy
	go mod download

# Инициализация рабочего окружения
init: deps
	mkdir -p bin result

# Форматирование исходников
fmt:
	gofmt -s -w .

# Lint (использует golangci-lint, если установлен; fallback -> go vet)
lint:
	@command -v golangci-lint >/dev/null 2>&1 || echo "golangci-lint not found: fallback to 'go vet'"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		go vet ./...; \
	fi