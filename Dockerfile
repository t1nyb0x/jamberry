# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build arguments for version info
ARG VERSION=unknown
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown

# Build with version info
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-X github.com/t1nyb0x/jamberry/internal/version.Version=${VERSION} \
              -X github.com/t1nyb0x/jamberry/internal/version.GitCommit=${GIT_COMMIT} \
              -X github.com/t1nyb0x/jamberry/internal/version.BuildDate=${BUILD_DATE}" \
    -o jamberry ./cmd/server

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/jamberry .

# Run as non-root user
RUN adduser -D -g '' appuser
USER appuser

ENTRYPOINT ["./jamberry"]
