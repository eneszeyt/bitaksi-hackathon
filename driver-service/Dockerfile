# Start from the official Golang image
FROM golang:1.24-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files first (for better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
# We build the file located at cmd/server/main.go and name the binary "driver-service"
RUN go build -o driver-service cmd/server/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./driver-service"]