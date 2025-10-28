# The build stage
FROM golang:1.25 AS builder
WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the binary statically for Linux (Cloud Run requires Linux binaries)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api cmd/api/*.go

# The run stage
FROM gcr.io/distroless/base-debian12
WORKDIR /app
# Copr CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/api .
ENV PORT=8080
EXPOSE 8080
CMD ["./api"]
