DB_FILE=./db.db
MIGRATIONS_DIR=./migrations
MIGRATE_BIN=migrate


.PHONY: migrate-up
migrate-up:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "sqlite3://$(DB_FILE)" up

.PHONY: migrate-down
migrate-down:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "sqlite3://$(DB_FILE)" down 1

.PHONY: run
run:
	go run ./cmd

.PHONY: build
build:
	go build -o ./bin/TwitchDonation ./cmd