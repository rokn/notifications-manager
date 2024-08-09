FROM golang:1.22.0

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy sources
COPY pkg/ ./pkg
COPY cmd/notifier_service/ ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /notifier_service

EXPOSE 8080

# Run
CMD ["/notifier_service"]
