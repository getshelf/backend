.PHONY: run test test-unit test-integration vet build migration-create docker-up docker-down docker-reset

MIGRATIONS_DIR=internal/infrastructure/persistence/postgres/migrations
MIGRATE_PACKAGE=github.com/golang-migrate/migrate/v4/cmd/migrate

run:
	go run ./cmd/api

test:
	go test ./...

test-unit:
	go test ./test/unit/...

test-integration:
	go test ./test/integration/...

vet:
	go vet ./...

build:
	mkdir -p bin
	go build -o bin/api ./cmd/api

migration-create:
	@test -n "$(name)" || (echo "usage: make migration-create name=create_users" && exit 1)
	go run $(MIGRATE_PACKAGE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

docker-up:
	docker compose --env-file .env -f deployments/compose.yml up --build -d

docker-down:
	docker compose --env-file .env -f deployments/compose.yml down

docker-reset:
	docker compose --env-file .env -f deployments/compose.yml down -v
	docker compose --env-file .env -f deployments/compose.yml up --build -d
