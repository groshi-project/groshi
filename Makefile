.PHONY: all
all: help

.PHONY: help
help:
	@echo "COMMAND             DESCRIPTION"
	@echo "secrets             create directories files for secrets"

.PHONY: secrets
secrets:
	mkdir -p ./.secrets/groshi/ ./.secrets/postgres
	touch ./.secrets/groshi/jwt_secret_key
	touch ./.secrets/postgres/user ./.secrets/postgres/password
