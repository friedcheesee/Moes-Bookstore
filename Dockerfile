# Use an official Go runtime as a parent image
FROM golang:1.21.0

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose a port for your Go application
EXPOSE 8080

# Run your application
CMD ["./main"]