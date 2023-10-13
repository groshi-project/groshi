GROSHI_TEST_SOCKET := http://127.0.0.1:80

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
	mkdir ./secrets/ ./secrets/app/ ./secrets/mongo
	touch ./secrets/app/exchangerates_api_key ./secrets/app/jwt_secret_key
	touch ./secrets/mongo/username ./secrets/mongo/password ./secrets/mongo/database

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

