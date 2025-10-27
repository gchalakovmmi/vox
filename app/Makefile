BINARY_NAME := app

.PHONY: compile

generate:
	@echo "Generating Templ Files..."
	templ generate

compile:
	@echo "Compiling Binary..."
	CGO_ENABLED=0 go build -o cmd/app/$(BINARY_NAME) cmd/app/main.go

run:
	@echo "Running Binary"
	export $$(grep -v '^#' .env | xargs) && ./cmd/app/$(BINARY_NAME)

all: generate compile run
