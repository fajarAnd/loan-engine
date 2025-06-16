SERVICE_NAME=loan-service


test:
	@echo "$(GREEN)Running all tests...$(NC)"
	go test ./... -cover

run:
	@echo "$(GREEN)Running application...$(NC)"
	go run main.go

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
	migrate create -ext sql -dir migrations $$name
# Install migration tool
install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Database seeding
seed:
	go run cmd/seed/main.go