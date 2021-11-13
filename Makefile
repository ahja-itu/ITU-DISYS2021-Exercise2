.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./src/service/service.proto;

.PHONY: run
run:
	go run src/*.go

.PHONY: docker-build
docker-build:
	docker build -t disysexercise2 .


.PHONY: up
up:
	docker-compose up

.PHONY: down
down:
	docker-compose down


.PHONY: refresh
refresh:
	docker build -t disysexercise2 .
	docker-compose up
