GOBIN = $$(go env GOPATH)/bin
LINTER = $(GOBIN)/golangci-lint
LINTER_VERSION = 1.55.0
CURRENT_DIR=$(shell pwd)
PACKAGE_DIRS=`go list -e ./... | egrep -v "binary_output_dir|.git|mocks"`
LABEL=$(shell git log -1 --format=%h)

.PHONY: path
path:
	export PATH=$PATH:$HOME/go/bin
########################################################################################################################
# LINT #################################################################################################################
########################################################################################################################
.PHONY: linter-check
linter-check:
	@if ! [ -x "$(LINTER)" ] || [ "$$($(LINTER) version | grep -o 'version [0-9\.]*')" != "version $(LINTER_VERSION)" ]; then \
		echo "Installing golangci-lint $(LINTER_VERSION)"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v$(LINTER_VERSION); \
	else \
		echo "golangci-lint $(LINTER_VERSION) is already installed"; \
	fi

.PHONY: install-mockgen
install-mockgen:
	go install github.com/golang/mock/mockgen@v1.6.0
	go get github.com/golang/mock/mockgen/model

.PHONY: lint
lint: linter-check
	$(LINTER) run ./...

.PHONY: lint-fix
lint-fix: linter-check
	$(LINTER) run --fix ./...
	go fmt ./...

########################################################################################################################
#  TEST ################################################################################################################
########################################################################################################################
.PHONY: unit-test
unit-test: test-package-dirs test-report

.PHONY: test-package-dirs
test-package-dirs:
	@echo "Running tests in package directories..."
	go test $(PACKAGE_DIRS) -race -coverprofile=cover.out -covermode=atomic

.PHONY: test-report
test-report:
	@echo 'Generating test coverage report...'
	go tool cover -html=cover.out -o cover.html

	@echo 'Generating overall test coverage ...'
	go tool cover -func cover.out | grep total:

.PHONY: integration-test
integration-test:
	make up-test-pg
	make wait-test-pg
	$(ENABLE_PG) go test $(PACKAGE_DIRS) -race -coverprofile=integration_cover.out -covermode=atomic
	@echo 'Generating test coverage report...'
	go tool cover -html=cover.out -o cover.html
	@echo 'Generating overall test coverage ...'
	go tool cover -func cover.out | grep total:
	make down-test-pg
	@echo "Integration tests completed successfully!"

.PHONY: test
test: clean-test up-test-pg wait-test-pg up-test-rabbitmq wait-test-rabbitmq
	@echo "Running unit tests..."
	go test $(PACKAGE_DIRS) -race -coverprofile=cover.out -covermode=atomic
	@echo "Running integration tests..."
	$(ENABLE_PG) go test $(PACKAGE_DIRS) -race -coverprofile=integration_cover.out -covermode=atomic
	@echo "Merging coverage reports..."
	echo "mode: atomic" > combined_cover.out
	{ grep -h -v "^mode:" cover.out integration_cover.out; } >> combined_cover.out
	go tool cover -html=combined_cover.out -o cover.html
	@echo "Combined test coverage report generated at cover.html"
	@echo 'Generating overall test coverage ...'
	go tool cover -func=combined_cover.out | grep total:
	make down-test-pg
	make down-test-rabbitmq

.PHONY: clean-test
clean-test:
	@rm -f cover.out integration_cover.out combined_cover.out cover.html
	@echo "Cleaned up test artifacts."

########################################################################################################################
# TEST - POSTGRES ######################################################################################################
########################################################################################################################
ENABLE_PG = POSTGRES_DSN=postgres://postgres:password@localhost:8432/postgres?sslmode=disable

.PHONY: up-test-pg
up-test-pg:
	docker run -d --name test-pg -e POSTGRES_PASSWORD=password -p 8432:5432 postgres:latest

.PHONY: wait-test-pg
wait-test-pg:
	@echo "Waiting for PostgreSQL to be ready..."
	@for i in $$(seq 1 30); do \
		if docker exec test-pg pg_isready -U postgres > /dev/null 2>&1; then \
			echo "PostgreSQL is ready!"; \
			exit 0; \
		fi; \
		echo "PostgreSQL not ready, retrying in 5 seconds..."; \
		sleep 5; \
	done; \
	echo "PostgreSQL did not become ready in time."; \
	exit 1

.PHONY: test-pg
test-pg:
	$(ENABLE_PG) make test

.PHONY: down-test-pg
down-test-pg:
	docker stop test-pg && docker rm test-pg


########################################################################################################################
# TEST - RABBITMQ ######################################################################################################
########################################################################################################################
.PHONY: up-test-rabbitmq
up-test-rabbitmq:
	docker run -d --name test-rabbitmq -p 4672:5672 -p 14672:15672 rabbitmq &

.PHONY: wait-test-rabbitmq
wait-test-rabbitmq:
	@echo "Waiting for RabbitMQ to be ready..."
	@for i in $$(seq 1 30); do \
		if nc -z localhost 4672; then \
			echo "RabbitMQ is ready!"; \
			exit 0; \
		fi; \
		echo "RabbitMQ not ready, retrying in 5 seconds..."; \
		sleep 5; \
	done; \
	echo "RabbitMQ did not become ready in time."; \
	docker logs test-rabbitmq; \
	exit 1

.PHONY: down-test-rabbitmq
down-test-rabbitmq:
	docker stop test-rabbitmq && docker rm test-rabbitmq

########################################################################################################################
# MOCKS ################################################################################################################
########################################################################################################################

.PHONY: mock-store
mock-store:
	@mockgen -destination=./mocks/store/mock_store.go -package=store message-app/repository  Store

.PHONY: mock-jwt-service
mock-jwt-service:
	@mockgen -destination=./mocks/services/mock_jwt_service.go -package=services message-app/services JWTService

.PHONY: mock-auth-service
mock-auth-service:
	@mockgen -destination=./mocks/services/mock_auth_service.go -package=services message-app/services AuthenticationService

.PHONY: mock-config-service
mock-config-service:
	@mockgen -destination=./mocks/services/mock_conf_service.go -package=services message-app/services GbeConfigService

.PHONY: mock-message-service
mock-message-service:
	@mockgen -destination=./mocks/services/mock_message_service.go -package=services message-app/services MessageService

.PHONY: mock-producer
mock-producer:
	@mockgen -destination=./mocks/pkg/rabbitmq/mock_producer.go -package=rabbitmq message-app/pkg/rabbitmq Producer

.PHONY: mock-consumer
mock-consumer:
	@mockgen -destination=./mocks/pkg/rabbitmq/mock_consumer.go -package=rabbitmq message-app/pkg/rabbitmq Consumer

.PHONY: mocks
mocks:
	@echo "Generating mocks..."
	$(MAKE) mock-store
	$(MAKE) mock-jwt-service
	$(MAKE) mock-config-service
	$(MAKE) mock-auth-service
	$(MAKE) mock-message-service
	$(MAKE) mock-producer
	$(MAKE) mock-consumer