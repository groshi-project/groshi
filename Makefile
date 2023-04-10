SOURCES := ./groshi.go ./internal/

.PHONY: all
all: help

.PHONY: help
help:
	@echo "COMMAND      DESCRIPTION                                 "
	@echo "---------------------------------------------------------"
	@echo "make gofmt   format source code using gofmt              "
	@echo "make todo    grep TODOs                                  "

.PHONY: gofmt
gofmt:
	gofmt -w $(SOURCES)

.PHONY: todo
todo:
	grep -irn todo $(SOURCES)
