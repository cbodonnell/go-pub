# Build stage
FROM golang:1-alpine as builder

RUN apk add git make

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN make build

# Production stage
FROM alpine:latest as prod

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/bin/pub ./

CMD [ "./pub", "serve" ]
