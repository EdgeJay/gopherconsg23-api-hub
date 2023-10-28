generate-mock-api-server:
	@oapi-codegen \
	-templates ./templates/mockapiserver/ \
	-generate types,server \
	-package main \
	./docs/savings-api/savings-api.yaml > ./cmd/mockapiserver/server.go
	@echo "Mock server codebase generated"

start-mock-savings-api: generate-mock-api-server
	@go run cmd/mockapiserver/*.go -input ./docs/savings-api/savings-api.yaml

build-docker-mock-savings-api:
	@docker build -t mock-savings-api:latest -f ./scripts/docker/mockapiserver.Dockerfile .

build-docker-mock-savings-api-logged:
	@docker build --no-cache --progress=plain -t mock-savings-api:latest -f ./scripts/docker/mockapiserver.Dockerfile . &> build.log

build-mock-savings-api:
	@go build -o ./build/mockapiserver ./cmd/mockapiserver/*.go

run-docker-mock-savings-api:
	@docker run -v `pwd`/docs/savings-api:/app/docs/savings-api \
	-p 1337:1337 -it --name mock-savings-api mock-savings-api:latest \
	/app/docs/savings-api/savings-api.yaml

stop-docker-mock-savings-api:
	@docker stop mock-savings-api

.PHONY: start-mock-savings-api generate-mock-api-server build-mock-savings-api build-mock-savings-api-logged run-mock-savings-api