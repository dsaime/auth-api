`.PHONY: test vet lint check run build

CORES := $(shell echo $$(( $$(nproc) - 1 )))

test:
	echo "Запуск go test..."
	go test ./...

vet:
	echo "Запуск go vet..."
	go vet ./...

lint:
	echo "Запуск golangci-lint на $(CORES) ядрах..."
# Параметр -j для golangci-lint ограничивает использование количество используемых ядер до "всего ядер - 1" (чтобы избежать перегрузки системы)
	golangci-lint run -v -j $(CORES)

check: vet lint

run:
	echo "Запуск сервера..."
	go run ./cmd/auth-api/main.go

build:
	echo "Компиляция сервера без запуска..."
	CGO_ENABLED=0 GOOS=linux go build -o bin/auth-api cmd/auth-api/main.go