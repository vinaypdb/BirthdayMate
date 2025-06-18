# ---------- Build stage ----------
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o app main.go

# ---------- Final stage ----------
FROM alpine:latest

WORKDIR /root
COPY --from=builder /app/app .

EXPOSE 9090
CMD ["./app"]
