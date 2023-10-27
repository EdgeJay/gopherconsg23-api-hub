FROM golang:1.21 AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
RUN make generate-mock-api-server
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/mock-api-server ./cmd/mock-api-server/*.go
RUN ls -al ./build

FROM alpine AS runner

COPY --from=builder /app/build/mock-api-server /app/mock-api-server
COPY --from=builder /app/docs/savings-api/savings-api.yaml /app/docs/savings-api/savings-api.yaml

EXPOSE 1337

CMD ["/app/mock-api-server", "-input", "/app/docs/savings-api/savings-api.yaml"]
