GROSHI_TEST_SOCKET := http://127.0.0.1:8080

.PHONY: all
all: help

.PHONY: help
help:
	@echo "COMMAND             DESCRIPTION"
	@echo "secrets             create directories files for secrets"
	@echo "docs                generate groshi API Swagger documentation"
	@echo "fmt                 format source code using 'swag fmt' and 'go fmt'"
	@echo "test-integration    run integration tests"

.PHONY: secrets
secrets:
	mkdir -p ./.secrets/groshi/ ./.secrets/groshi-mongo
	touch ./.secrets/groshi/exchangerates_api_key ./.secrets/groshi/jwt_secret_key
	touch ./.secrets/groshi-mongo/username ./.secrets/groshi-mongo/password ./.secrets/groshi-mongo/database

.PHONY: docs
docs:
	swag init .

.PHONY: fmt
fmt:
	swag fmt
	go fmt

.PHONY: test-integration
test-integration:
	GROSHI_TEST_SOCKET=$(GROSHI_TEST_SOCKET) go test -count=1 ./tests/integration

