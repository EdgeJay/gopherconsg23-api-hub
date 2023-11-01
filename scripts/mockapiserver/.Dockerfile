FROM golang:1.21-alpine

RUN apk add --no-cache bash
RUN which bash

WORKDIR /app

COPY . .

RUN ls -al
RUN go mod download
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
RUN chmod +x ./scripts/mockapiserver/entrypoint.sh

ENTRYPOINT ["/app/scripts/mockapiserver/entrypoint.sh"]
