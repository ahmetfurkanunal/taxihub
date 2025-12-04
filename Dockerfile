# Build aşaması
FROM golang:1.25.4 AS builder

WORKDIR /app

# Go mod ve go sum kopyalanır, modüller indirilir
COPY go.mod go.sum ./
RUN go mod download

# Tüm kodu kopyala
COPY . .

# Binary build
RUN go build -o driver-service main.go

# Runtime imajı
FROM debian:bookworm-slim

WORKDIR /root/

COPY --from=builder /app/driver-service .

EXPOSE 8080

CMD ["./driver-service"]
