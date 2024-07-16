GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

# Load environment variables for local development
include .env
export

.PHONY: dev-db-migrate
# Run database migrations
dev-db-migrate:
	@go run cmd/migrate/main.go

.PHONY: dev-db-rollback
# Rollback database migrations
dev-db-rollback:
	@go run cmd/migrate/main.go -rollback

.PHONY: docker-build
# Build docker image
docker-build:
	@docker buildx build --pull \
		--platform=linux/amd64,linux/arm64 \
		--build-arg GOLANG_VERSION=${GOLANG_VERSION} \
		-f Dockerfile -t ${REGISTRY}:${IMAGE_TAG} --push .

.PHONY: docker-compose-run
# Run docker-compose
docker-compose-run:
	@docker-compose up --build

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
