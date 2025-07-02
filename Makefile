.PHONY: test vet lint check run build

CORES := $(shell echo $$(( $$(nproc) - 1 )))

test:
	go test ./...

vet:
	go vet ./...

lint:
# Параметр -j для golangci-lint ограничивает использование количество используемых ядер до "всего ядер - 1" (чтобы избежать перегрузки системы)
	golangci-lint run -v -j $(CORES)

check: vet lint

run:
	go run ./cmd/auth-api/main.go

build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/auth-api cmd/auth-api/main.go