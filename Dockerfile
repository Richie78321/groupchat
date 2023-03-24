FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN apk add --no-cache \
    # Required for building
    protobuf make \
    # Required for go-sqlite3
    gcc musl-dev \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Required for go-sqlite3
ENV CGO_ENABLED=1
RUN make protoc && make clientcli && make servercli

FROM alpine

# Required for networking emulation / testing
RUN apk add iproute2

COPY --from=builder /app/bin /bin
COPY --from=builder /app/scripts .

ENTRYPOINT ["/bin/sh"]
