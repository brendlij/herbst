# Build frontend
FROM node:22-alpine AS frontend-builder

WORKDIR /app/web

# Install dependencies
COPY web/package.json web/package-lock.json* web/bun.lock* ./
RUN npm install

# Build frontend
COPY web/ ./
RUN npm run build


# Build backend
FROM golang:1.25.2-alpine AS backend-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Download Go modules
COPY go.mod go.sum* ./
RUN go mod download

# Build binary
COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN CGO_ENABLED=0 GOOS=linux go build -o herbst ./cmd/herbst


# Final image
FROM alpine:3.21

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk add --no-cache ca-certificates tzdata

# Copy binary
COPY --from=backend-builder /app/herbst .

# Copy frontend build
COPY --from=frontend-builder /app/web/dist ./web/dist

# Create directories for config and static files
RUN mkdir -p /app/config /app/static

# Environment variables
ENV HERBST_CONFIG_DIR=/app/config
ENV HERBST_STATIC_DIR=/app/static

EXPOSE 8080

CMD ["./herbst"]
