init:
	@echo "📡 Initialize requirement"
	@.dev/initialize.sh

lint:
	@which golangci-lint 2>/dev/null || echo "please install golangci-lint"
	@golangci-lint run ./... -D errcheck -E depguard -E gofmt -E goimports -E nakedret -E goconst -D gosimple

test:
	@which go test 2>/dev/null || echo "please install github.com/rakyll/gotest"
	@GOSUMDB=off go test -short -gcflags="-l" -race ./...

test-coverage:
	@echo "🔎 Starting check unit test coverage"
	@GOSUMDB=off go test -short -coverprofile coverage.out ./... && go tool cover -func=coverage.out
	@echo "🎉 Checking unit test completed"

integration-test:
	@echo "🔎 Starting check integration test"
	@GOSUMDB=off go test -v -count=1 ./integration_test/...
	@echo "🎉 Checking integration test completed"

check-buildable:
	@echo "🔎 Check this code is buildable"
	@GOSUMDB=off go build -o /dev/null ./...

check-imports-newline:
	@echo "🔎 Check this code have proper new line"
	@.dev/check-imports-newline.sh

pre-commit:
	@echo "🔮 Running magic pre commit"
	@.dev/run-pre-commit.sh

pre-push:
	@echo "🔮 Running magic pre push"
	@.dev/run-pre-push.sh

mocks:
	@echo "🎩 Running mock generating"
	@which mockgen 2>/dev/null || echo "please install github.com/golang/mock/mockgen"
	@go generate $$(go list ./... | grep -v vendor)
	@echo "👌 Mocks generated"

schema-generate:
	@echo "🎩 Generating schema application"
	@which gqlgen 2>/dev/null || echo "please install github.com/99designs/gqlgen"
	@for f in $(shell ls ./graph/*.yml); do gqlgen generate --config $${f}; done
	@echo "👌 Schema generated"

services-up:
	@echo "🎬 Starting the service - $$(date)"
	@sudo chown -R $$(id -u):$$(id -g) .docker
	@echo "🎉 Service is going to UP!"
	@docker-compose -f docker-compose.yml up

services-down:
	@docker-compose -f docker-compose.yaml down

clear-postgres-docker-data:
	@sudo rm -rf .docker/postgres-data

download:
	@echo "📡 Starting download package dependencies"
	@go mod download -x
	@sleep 1
	@echo "📁 Setup vendor directory"
	@go mod vendor
	@sleep 1
	@echo "👌 Download package completed"

migrate:
	@.dev/run-migration.sh $(DRIVER) $(DB_USERNAME) $(DB_PASSWORD) $(DB_HOST) $(DB_NAME) $(SSL_MODE)

create-migration:
	@.dev/create-migration.sh ${DRIVER} $(NAME)
