# Multi-stage build for Autobattle server + frontend

# Build frontend assets
FROM node:20-alpine AS frontend-builder
WORKDIR /frontend
COPY frontend/package.json ./
RUN npm install
COPY frontend/. ./
RUN npm run build

# Build Go backend
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
COPY --from=frontend-builder /frontend/dist /app/frontend/dist
COPY openapi.yml /app/openapi.yml
COPY Readme.md /app/Readme.md

EXPOSE 8080

ENTRYPOINT ["/app/autobattle"]
