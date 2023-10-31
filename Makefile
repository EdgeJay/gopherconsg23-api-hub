include .env

generate-mock-api-data:
	@go run cmd/mockgenerator/*.go \
	-input ${MOCK_API_FILE} \
	-method GET \
	-path /savings-records \
	-status 200

generate-mock-api-server:
	@oapi-codegen \
	-templates ./templates/mockapiserver/ \
	-generate types,server \
	-package main \
	${MOCK_API_FILE} > ./cmd/mockapiserver/server.gen.go
	@echo "Mock server codebase generated"

start-mock-api-without-generation:
	@go run cmd/mockapiserver/*.go -input ${MOCK_API_FILE}

start-mock-api: generate-mock-api-server
	@go run cmd/mockapiserver/*.go -input ${MOCK_API_FILE}

build-docker-mock-api:
	@docker build -t mock-api:latest -f ./scripts/docker/mockapiserver.Dockerfile .

build-docker-mock-api-logged:
	@docker build --no-cache --progress=plain -t mock-savings-api:latest -f ./scripts/docker/mockapiserver.Dockerfile . &> build.log

build-mock-api:
	@go build -o ./build/mockapiserver ./cmd/mockapiserver/*.go

run-docker-mock-savings-api:
	@docker run -v `pwd`/docs/savings-api:/app/docs/savings-api \
	-p 1337:1337 -it --name mock-savings-api mock-api:latest \
	/app/docs/savings-api/savings-api.yaml

run-docker-mock-housing-api:
	@docker run -v `pwd`/docs/housing-api:/app/docs/housing-api \
	-p 1338:1337 -it --name mock-housing-api mock-api:latest \
	/app/docs/housing-api/housing-api.yaml

start-docker-mock-savings-api:
	@docker start -i mock-savings-api

start-docker-mock-housing-api:
	@docker start -i mock-housing-api

stop-docker-mock-savings-api:
	@docker stop mock-savings-api

stop-docker-mock-housing-api:
	@docker stop mock-housing-api

generate-mock-codegen-client:
	@oapi-codegen \
	-generate types,client \
	-package main \
	${MOCK_API_FILE} > ./cmd/mockcodegenclient/client.gen.go
	@echo "Mock client codebase generated"

.PHONY: generate-mock-api-data start-mock-api start-mock-api-without-generation generate-mock-api-server build-docker-mock-api build-docker-mock-api-logged build-mock-api run-docker-mock-savings-api start-docker-mock-savings-api stop-docker-mock-savings-api generate-mock-codegen-client