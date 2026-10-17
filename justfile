run: build
    @./bin/api/main

build:
    @go build -o bin/api/main ./cmd/api/main.go

build_migrate:
    @mkdir -p bin/migration
    @go build -o bin/migration/main ./cmd/migration

migrate-up:
    @./bin/migration/main up

migrate-down:
    @./bin/migration/main down
