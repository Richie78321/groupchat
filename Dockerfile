FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN apk add --no-cache protobuf make \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

RUN make protoc && make clientcli && make servercli

FROM alpine

COPY --from=builder /app/bin /bin

ENTRYPOINT ["/bin/sh"]
