# Build stage
FROM golang:1.24 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./main.go

################################################################
FROM alpine:3.21

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
ENTRYPOINT ["/app/main"]
