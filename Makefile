build:
	@go build -o bin/main cmd/app/main.go

run:build
	@./bin/main