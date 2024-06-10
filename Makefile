run: build
	@./bin/api.stick.it

build:
	@go build -o bin/api.stick.it cmd/api.stick.it/main.go

docker-up:
	# @docker volume create --name=db-data
	@docker compose -f cmd/docker/docker-compose.yml up -d

docker-down:
	@docker compose -f cmd/docker/docker-compose.yml down

build-migration:
	@go build -o bin/migrate cmd/migration/main.go

migrate-up: build-migration
	@./bin/migrate up

migrate-down: build-migration
	@./bin/migrate down

migrate-reset: build-migration
	@./bin/migrate reset

sqlc:
	@sqlc generate -f internal/db/sqlc.yaml

seed:
	@go run cmd/seed/main.go

host:
	# @go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@tar -C /opt -xzf sqlc_1.26.0_linux_amd64.tar.gz
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	$(MAKE) migrate-up
	$(MAKE) sqlc
	$(MAKE) build