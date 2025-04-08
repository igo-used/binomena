
FROM golang:1.22-alpine

WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o binomena

# Expose ports
EXPOSE 8080 9000

# Command to run
CMD ["./binomena", "--api-port", "8080", "--p2p-port", "9000", "--id", "genesis-node"]
