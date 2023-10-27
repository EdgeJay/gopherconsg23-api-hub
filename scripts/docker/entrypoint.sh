#!/bin/bash

set -e

INPUT_FILE="$1"
INPUT_HASH_FILE="/app/input_file.sha256"

echo "Loading OpenAPI specs file from $INPUT_FILE"

generate_input_hash () {
    # Generate hash of input file
    sha256sum "$INPUT_FILE" > "$INPUT_HASH_FILE"
    # cat /app/input_file.sha256
}

init_server () {
    oapi-codegen \
        -templates /app/cmd/mock-api-server/templates/ \
        -generate types,server \
        -package main \
        $INPUT_FILE > /app/cmd/mock-api-server/server.go

    echo "Mock server codebase generated, re-building server..."

    CGO_ENABLED=0 GOOS=linux go build -o /app/build/mock-api-server /app/cmd/mock-api-server/*.go

    echo "Server re-built"
}

if [[ -s $INPUT_HASH_FILE ]]; then
    if ! sha256sum -c $INPUT_HASH_FILE; then
        echo "OpenAPI specs file changed, re-building server..."
        generate_input_hash
        init_server
    fi
else
    generate_input_hash
    init_server
fi

echo "Starting up server..."
/app/build/mock-api-server -input $INPUT_FILE