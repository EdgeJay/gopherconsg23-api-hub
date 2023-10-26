generate-mock-savings-api-server:
	@cd cmd/mock-savings-api && mkdir -p server && oapi-codegen \
	-templates templates/ \
	-generate types,server \
	-package server \
	../../docs/savings-api/savings-api.yaml > ./server/server.go

start-mock-savings-api:
	@go run cmd/mock-savings-api/main.go

.PHONY: start-mock-savings-api generate-mock-savings-api-server