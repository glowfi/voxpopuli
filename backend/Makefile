# Target to run the Go program
run:
	go run cmd/main.go
migrate-up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=postgres://postgres:postgres@localhost:5432/voxpopuli GOOSE_MIGRATION_DIR=./migrations goose up
migrate-down:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=postgres://postgres:postgres@localhost:5432/voxpopuli GOOSE_MIGRATION_DIR=./migrations goose down
test:
	go test -timeout 300s -p 1 -race -v -count=1 ./...
lint:
	golangci-lint run &&misspell .
