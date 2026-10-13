run: 
    @./bin/api/main

build:
    @go build -o bin/api/main ./cmd/api/main.go
