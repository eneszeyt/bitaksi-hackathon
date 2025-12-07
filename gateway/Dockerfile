# Start from the official Golang image
FROM golang:1.24-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o gateway-service cmd/main.go

# Expose port 8000
EXPOSE 8000

# Command to run the executable
CMD ["./gateway-service"]