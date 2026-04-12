# Stage 1: Build the React client
FROM node:22-alpine AS client-builder
WORKDIR /app/client
COPY client/package.json client/package-lock.json ./
RUN npm install
COPY client/ ./
RUN npm run build

# Stage 2: Build the Go binary (with client embedded)
FROM golang:1.24-alpine AS server-builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
# Copy built client dist into the embed location
COPY --from=client-builder /app/client/dist ./internal/static/dist
RUN go build -tags embed -o bin/pub ./cmd/pub/main.go

# Stage 3: Minimal runtime image
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=server-builder /app/bin/pub ./
CMD ["./pub", "serve"]
