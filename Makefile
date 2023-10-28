generate-mock-api-server:
	@oapi-codegen \
	-templates ./templates/mock-api-server/ \
	-generate types,server \
	-package main \
	./docs/savings-api/savings-api.yaml > ./cmd/mock-api-server/server.go
	@echo "Mock server codebase generated"

start-mock-savings-api: generate-mock-api-server
	@go run cmd/mock-api-server/*.go -input ./docs/savings-api/savings-api.yaml

build-mock-savings-api:
	@docker build -t mock-savings-api:latest -f ./scripts/docker/mock_api_server.Dockerfile .

build-mock-savings-api-logged:
	@docker build --no-cache --progress=plain -t mock-savings-api:latest -f ./scripts/docker/mock_api_server.Dockerfile . &> build.log

run-mock-savings-api:
	@docker run -v `pwd`/docs/savings-api:/app/docs/savings-api \
	-p 1337:1337 -it --name mock-savings-api mock-savings-api:latest \
	/app/docs/savings-api/savings-api.yaml

stop-mock-savings-api:
	@docker stop mock-savings-api

.PHONY: start-mock-savings-api generate-mock-api-server build-mock-savings-api build-mock-savings-api-logged run-mock-savings-api