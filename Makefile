host = localhost
port = 8080
run:
	HOST=$(host) PORT=$(port) go run ./cmd/pony/main.go
test:
	go test ./...
