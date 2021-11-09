
.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./src/service/service.proto;

.PHONY: run
run:
	go run src/agent/*.go
