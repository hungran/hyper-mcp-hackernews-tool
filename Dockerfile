# FROM tinygo/tinygo:0.36.0 AS builder

# WORKDIR /app
# COPY go.mod .
# COPY go.sum .
# RUN go mod download
# COPY . .
# ARG BUILDVCS=false
# ENV GOFLAGS="-buildvcs=${BUILDVCS}"
# RUN GOOS=wasip1 GOARCH=wasm tinygo build -o /app/plugin.wasm -target=wasi .


FROM scratch
WORKDIR /
# COPY --from=builder /app/plugin.wasm /plugin.wasm
COPY plugin.wasm /plugin.wasm
