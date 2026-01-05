# Multi-stage build for minimal image
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o facilitator ./cmd/server

# Final minimal image
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/facilitator /facilitator

EXPOSE 3000

ENTRYPOINT ["/facilitator"]
