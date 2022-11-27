build:
	go build -o bin/pub cmd/pub/main.go

run:
	go run cmd/pub/main.go

container:
	docker build -f deploy/docker/Dockerfile -t cheebz/go-pub .