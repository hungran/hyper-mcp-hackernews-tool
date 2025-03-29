FROM tinygo/tinygo:0.37.0 AS builder

WORKDIR /workspace
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN tinygo build -o plugin.wasm -target=wasi -no-debug .

FROM scratch
WORKDIR /
COPY --from=builder /workspace/plugin.wasm /plugin.wasm
