.PHONY: all
all:
	@echo "COMMAND    DESCRIPTION"
	@echo "secrets    Create directories files for secrets."
	@echo "docs       Generate groshi API Swagger documentation."
	@echo "fmt        Format source code using `swag fmt` and `go fmt`."

.PHONY: docs
docs:
	swag init .

.PHONY: fmt
fmt:
	swag fmt
	go fmt

.PHONY: secrets
secrets:
	mkdir ./secrets/ ./secrets/app/ ./secrets/mongo
	touch ./secrets/app/exchangerats_api_key ./secrets/app/jwt_secret_key
	touch ./secrets/mongo/username ./secrets/mongo/password ./secrets/mongo/database
