.PHONY: all
all:
	@echo "COMMAND    DESCRIPTION"
	@echo "docs       Generate groshi API Swagger documentation."

.PHONY: docs
docs:
	swag init .
