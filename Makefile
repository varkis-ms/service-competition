# Commands
.PHONY: env
env:
	@$(eval SHELL:=/bin/bash)
	@cp .env.sample .env

.PHONY: build
build:
	go build -o ./cmd/service

.PHONY: run
run:
	go run ./cmd/service

.PHONY: compose-up
compose-up:
	docker-compose -f docker-compose.yml up -d --remove-orphans

.PHONY: compose-down
compose-down:
	docker-compose down --remove-orphans

.PHONY: unit-test
unit-test:
	go test ./...

.PHONY: cover
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm coverage.out

.PHONY: linter
linter:
	golangci-lint run

.PHONY: swagger
swagger:
	swag init -g internal/app/app.go --parseInternal --parseDependency

.PHONY: proto
proto:
	rm -f internal/pkg/pb/*.go

	protoc -I proto proto/*.proto --go_out=./internal/pkg/pb --go_opt=paths=source_relative \
	--go-grpc_out=./internal/pkg/pb --go-grpc_opt=paths=source_relative

	protoc -I proto proto/clients/*.proto --go_out=./internal/pkg/pb --go_opt=paths=source_relative \
    	--go-grpc_out=./internal/pkg/pb --go-grpc_opt=paths=source_relative