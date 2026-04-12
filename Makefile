BINARY := bin/pub
CMD    := ./cmd/pub/main.go

.PHONY: build build-embed serve dev-client compose-up compose-down compose-build docker-build

# Build binary (no embedded client; use CLIENT env var to point at Vite dev server)
build:
	go build -o $(BINARY) $(CMD)

# Build binary with client embedded (requires running `make client-build` first)
build-embed:
	go build -tags embed -o $(BINARY) $(CMD)

# Build the React client
client-build:
	cd client && npm ci && npm run build && cp -r dist ../internal/static/dist

# Run the server (expects CLIENT env var or embedded dist)
serve:
	./$(BINARY) serve

# Run the Vite dev server (proxies API calls to localhost:8080)
dev-client:
	cd client && npm run dev

# Docker Compose (from repo root, references deploy/compose/)
compose-up:
	docker compose -f deploy/compose/docker-compose.yml --env-file deploy/compose/.env up -d

compose-down:
	docker compose -f deploy/compose/docker-compose.yml --env-file deploy/compose/.env down

compose-build:
	docker compose -f deploy/compose/docker-compose.yml --env-file deploy/compose/.env build

# Build the production Docker image directly
docker-build:
	docker build -f deploy/docker/Dockerfile -t cheebz/go-pub .
