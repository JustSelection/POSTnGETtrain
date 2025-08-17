DB_HOST := localhost
DB_PORT := 5432
DB_NAME := postgres
DB_USER := postgres
DB_PASS := "yourpassword"

# Важно: Без пробелов и с правильным экранированием
DB_DSN := "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"
MIGRATE := migrate -path ./migrations -database $(DB_DSN)

# Команды PSQL (без -it для автоматизации)
PSQL := docker exec postgres-container psql -U $(DB_USER) -d $(DB_NAME)

migrate-new:
	@echo "Creating new migration: $(NAME)"
	$(MIGRATE) create -ext sql -dir ./migrations -seq $(NAME)

migrate:
	@echo "Applying migrations to $(DB_HOST):$(DB_PORT)"
	$(MIGRATE) up

migrate-down:
	@echo "Rolling back migrations"
	$(MIGRATE) down

drop-tasks:
	@echo "Dropping tasks table..."
	$(PSQL) -c "DROP TABLE IF EXISTS tasks;"

force-fix:
	@echo "Forcing version fix"
	$(MIGRATE) force 20250813124442
	
run:
	@echo "Starting server..."
	go run cmd/main.go

check-db:
	@echo "Testing DB connection..."
	$(PSQL) -c "\dt"
	
migrate-down-hard:
	$(MIGRATE) down
	$(PSQL) -c "DROP TABLE IF EXISTS tasks;"

gen:
	oapi-codegen -config openapi/.openapi -include-tags tasks -package tasks openapi/openapi.yaml > ./internal/web/tasks/api.gen.go
	
lint:
	golangci-lint run --color=auto
	
clean-linter-cache:
	golangci-lint cache clean