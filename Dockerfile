# Stage 1: Build the Go app
FROM golang:1.21 AS builder

WORKDIR /app
COPY . .
RUN go mod init vnypdb-app
RUN go mod tidy
RUN go build -o app main.go

# Stage 2: Final lightweight image
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/app .

EXPOSE 8080
CMD ["./app"]
