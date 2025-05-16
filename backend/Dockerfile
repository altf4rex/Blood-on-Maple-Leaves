# -------- stage 1: build --------
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd

# -------- stage 2: runtime --------
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server        .
# копируем SQL и YAML – нужны рантайму
COPY --from=builder /app/migrations    ./migrations
COPY --from=builder /app/scenes        ./scenes
ENTRYPOINT ["./server"]
