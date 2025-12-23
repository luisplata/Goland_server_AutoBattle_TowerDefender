# Multi-stage build for Autobattle server
FROM golang:1.22.5-alpine AS builder
WORKDIR /app

# Preload dependencies
COPY go.mod .
RUN go mod download

# Copy source
COPY . .

# Build statically-linked binary for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/autobattle ./

# Runtime image
FROM gcr.io/distroless/base-debian12
WORKDIR /app

COPY --from=builder /bin/autobattle /app/autobattle

EXPOSE 8080

ENTRYPOINT ["/app/autobattle"]
