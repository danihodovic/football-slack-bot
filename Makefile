GO_FILES = $(shell find ! -name '*_test.go' -name '*.go')

# Starts the container
.PHONY: start
start:
	docker-compose run app go run $(GO_FILES) --config config.json

.PHONY: test
test:
	docker-compose run app go test

.PHONY: interactive
interactive:
	docker-compose run app bash

# Run the app in a container. Used for debugging purposes
.PHONY: run
run:
	go run $(GO_FILES) --config config.json
