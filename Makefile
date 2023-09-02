.PHONY: all
all:
	@echo "COMMAND    DESCRIPTION"
	@echo "docs       Generate groshi API Swagger documentation."
	@echo "fmt        Format source code using `swag fmt` and `go fmt`."

.PHONY: docs
docs:
	swag init .

.PHONY: fmt
fmt:
	swag fmt
	go fmt
