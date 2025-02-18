# Target to run the Go program
run:
	@go run cmd/main.go
migrate:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=postgres://postgres:postgres@localhost:5432/voxpopuli GOOSE_MIGRATION_DIR=./migrations goose up
test:
	go test -timeout 60s -p 1 -race -v -count=1 ./...
