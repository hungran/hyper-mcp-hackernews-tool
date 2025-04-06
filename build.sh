#!/bin/bash
set -e

# Download dependencies
go mod tidy

# Build the WebAssembly module with TinyGo - suppress debug info
tinygo build -o plugin.wasm -target=wasi -no-debug .
docker build -t ghcr.io/hungran/hyper-mcp-hackernews-tool .
echo "Build completed successfully. Plugin is at plugin.wasm"
