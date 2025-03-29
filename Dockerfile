FROM tinygo/tinygo:0.37.0 AS builder

WORKDIR /workspace
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ARG BUILDVCS=false
ENV GOFLAGS="-buildvcs=${BUILDVCS}"
RUN tinygo build -target wasi -o plugin.wasm .

FROM scratch
WORKDIR /
COPY --from=builder /workspace/plugin.wasm /plugin.wasm
