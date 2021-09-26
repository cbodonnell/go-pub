# Build stage
FROM golang:1-alpine as builder

RUN mkdir /app
WORKDIR /app

COPY . .

RUN go get -d -v ./...
RUN go build

# Production stage
FROM alpine:latest as prod

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/go-pub ./

CMD [ "./go-pub" ]