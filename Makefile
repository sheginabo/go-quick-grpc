IMAGE=quick-grpc-api
TAG=latest

all: build_app

build_app:
	go mod download
	go build -o bin/app ./cmd

test:
	go mod download
	go test ./...

local_run:
	go mod download
	go run cmd/main.go

run:
	./bin/app

clean:
	rm -rf ./bin

docker_build:
	docker build -t ${IMAGE}:${TAG} .

docker_run:
	docker run --rm -d -p 8080:8080 --name ${IMAGE} ${IMAGE}:${TAG}

docker_stop:
	docker stop ${IMAGE}

docker_log:
	docker logs ${IMAGE}

docker_test:
	docker build --rm -f Dockerfile.test .
server:
	go run ./cmd/main.go
proto_grpc:
	rm -f ./internal/pb/*.go
	protoc --proto_path=./internal/proto --go_out=./internal/pb --go_opt=paths=source_relative \
        --go-grpc_out=./internal/pb --go-grpc_opt=paths=source_relative \
        ./internal/proto/*.proto