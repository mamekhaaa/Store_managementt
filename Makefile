.PHONY: migrate
migrate:
	go run migrations/migrate.go -dsn="$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?parseTime=true"

.PHONY: migrate-prod
migrate-prod:
	go run migrations/migrate.go -dsn="$(PROD_DB_USER):$(PROD_DB_PASSWORD)@tcp($(PROD_DB_HOST):$(PROD_DB_PORT))/$(PROD_DB_NAME)?parseTime=true"