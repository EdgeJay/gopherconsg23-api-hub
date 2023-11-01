#!/bin/bash

set -e

INPUT_FILE="$1"
SERVER_PORT="${2:-1337}"
INPUT_HASH_FILE="/app/input_file.sha256"

echo "Server port set to $SERVER_PORT"
echo "Loading OpenAPI specs file from $INPUT_FILE"

generate_input_hash () {
    # Generate hash of input file
    sha256sum "$INPUT_FILE" > "$INPUT_HASH_FILE"
    # cat /app/input_file.sha256
}

init_server () {
    oapi-codegen \
        -templates /app/templates/mockapiserver/ \
        -generate types,server \
        -package main \
        $INPUT_FILE > /app/cmd/mockapiserver/server.gen.go

    echo "Mock server codebase generated, re-building server..."

    CGO_ENABLED=0 GOOS=linux go build -o /app/build/mockapiserver /app/cmd/mockapiserver/*.go

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
/app/build/mockapiserver -input $INPUT_FILE -port $SERVER_PORT