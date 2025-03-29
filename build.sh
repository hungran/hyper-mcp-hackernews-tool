#!/bin/bash
set -e

# Download dependencies
go mod tidy

# Build the WebAssembly module with TinyGo - suppress debug info
tinygo build -o plugin.wasm -target=wasi -no-debug .

echo "Build completed successfully. Plugin is at plugin.wasm" 
