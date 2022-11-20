FROM tinygo/tinygo:0.26.0 as builder

FROM golang:alpine3.16
COPY --from=builder /usr/local/tinygo /usr/local/tinygo
RUN chmod +x /usr/local/tinygo

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct

COPY . /go/src/envoy-wasm

RUN /usr/local/tinygo build -o optimized.wasm -scheduler=none -target=wasi main.go
