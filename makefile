DB_FILE=./db.db
MIGRATIONS_DIR=./migrations
MIGRATE_BIN=migrate

.PHONY: migrate-create
migrate-create:
	$(MIGRATE_BIN) create  -seq -ext=.sql -dir $(MIGRATIONS_DIR) "$(word 2, $(MAKECMDGOALS))"

.PHONY: migrate-up
migrate-up:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "sqlite3://$(DB_FILE)" up

.PHONY: migrate-down
migrate-down:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "sqlite3://$(DB_FILE)" down 1

.PHONY: run
run:
	go run ./cmd/app

.PHONY: build-prod
build-prod:
	CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -o ./bin/TwitchDonation ./cmd/app

.PHONY: build
build:
	go build -o ./bin/TwitchDonation ./cmd/app