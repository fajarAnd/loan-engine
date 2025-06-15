SERVICE_NAME=loan-service


test:
	go test ./... -cover


wire:
	wire gen github.com/go-bonapp/splitbillservice/infrastructure/http/router

# Database migration targets
migrate-up:
	go run cmd/migrate/main.go -action=up

migrate-down:
	go run cmd/migrate/main.go -action=down

migrate-force:
	@read -p "Enter version number: " version; \
	go run cmd/migrate/main.go -action=force -version=$$version

migrate-version:
	go run cmd/migrate/main.go -action=version

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Database setup
db-setup: migrate-up
	@echo "Database setup completed"

db-reset: migrate-down migrate-up
	@echo "Database reset completed"

# Install migration tool
install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@lates
