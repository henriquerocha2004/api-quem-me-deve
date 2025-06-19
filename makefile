# Configurações
DB_URL=postgres://devuser:devpass@postgres:5432/devdb?sslmode=disable
MIGRATIONS_DIR=./internal/database/migrations
MIGRATE_BIN=migrate

# Comando para criar uma nova migration: make new name=create_users_table
migrate-new:
	@read -p "Migration name: " name; \
	$(MIGRATE_BIN) create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name

# Rodar as migrations
migrate-up:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

# Reverter última migration
migrate-down:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

# Dropar todo o schema e rodar tudo novamente
migrate-redo:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" drop -f
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

# Mostrar versão atual
migrate-version:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

refresh-schema:
	pg_dump --schema-only --no-owner --file=./internal/database/schema/schema.sql -d $(DB_URL)

sqlc-generate:
	sqlc generate --file=./sqlc.yaml

build-db:
	make migrate-up
	make refresh-schema
	make sqlc-generate

# Mostrar ajuda
help:
	@echo "Comandos disponíveis:"
	@echo "  make new            Cria uma nova migration (interativo)"
	@echo "  make up             Aplica as migrations"
	@echo "  make down           Desfaz a última migration"
	@echo "  make redo           Refaz todas as migrations (drop + up)"
	@echo "  make version        Mostra a versão atual"