.PHONY: run-server run-cli test test-cover build clean

запуск сервера:
	go run cmd/server/main.go

запуск клиента:
	go run cmd/cli/main.go -from USD -to RUB -amount 100

запуск юнит теста:
	go test -race -v ./...

запуск теста на покрытие:
	go test -race -coverprofile=coverage.out ./...

сборка бинаря:
	go build -o bin/server cmd/server/main.go
	go build -o bin/cli cmd/cli/main.go

чистка файлов бинаря и теста покрытия:
	rm -rf bin/
	rm -f coverage.out coverage.html
