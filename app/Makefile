BINARY_NAME := app

.PHONY: compile

compile:
	@echo "Compiling Binary..."
	go build -o cmd/app/$(BINARY_NAME) cmd/app/main.go

run:
	@echo "Running Binary"
	export $$(grep -v '^#' .env | xargs)
	./cmd/app/$(BINARY_NAME)
