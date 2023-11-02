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
	@docker build -t mock-api:latest -f ./scripts/mockapiserver/.Dockerfile .

build-docker-mock-api-logged:
	@docker build --no-cache --progress=plain -t mock-savings-api:latest -f ./scripts/mockapiserver/.Dockerfile . &> build.log

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

generate-codegen-client:
	@oapi-codegen \
	-generate types,client \
	-package main \
	${MOCK_API_FILE} > ./cmd/codegenclient/client.gen.go
	@echo "Mock client codebase generated"

start-codegen-client: generate-codegen-client
	@go run ./cmd/codegenclient/*.go -input ${MOCK_API_FILE}

generate-codegen-multi-client:
	@mkdir -p ./cmd/codegenmulticlient/savings && \
	mkdir -p ./cmd/codegenmulticlient/housing
	@oapi-codegen \
	-generate types,client \
	-package savings \
	./docs/savings-api/savings-api.yaml > ./cmd/codegenmulticlient/savings/client.gen.go
	@oapi-codegen \
	-generate types,client \
	-package housing \
	./docs/housing-api/housing-api.yaml > ./cmd/codegenmulticlient/housing/client.gen.go
	@echo "Mock client codebase generated"

start-codegen-multi-client: generate-codegen-multi-client
	@go run ./cmd/codegenmulticlient/*.go -input ./docs/savings-api/savings-api.yaml -input ./docs/housing-api/housing-api.yaml
