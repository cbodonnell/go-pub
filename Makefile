build:
	go build -o bin/pub cmd/pub/main.go

serve:
	./bin/pub serve

container:
	docker build -f deploy/docker/Dockerfile -t cheebz/go-pub .