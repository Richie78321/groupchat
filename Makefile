protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	chatservice/chat_service.proto

servercli:
	go build -o bin/ cmd/servercli/servercli.go

clientcli:
	go build -o bin/ cmd/clientcli/clientcli.go
