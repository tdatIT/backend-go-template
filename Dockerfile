# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a smaller static binary
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /out/app ./cmd/main.go

FROM alpine:3.21.6

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /out/app ./app

EXPOSE 5000

ENTRYPOINT ["/app/app"]
