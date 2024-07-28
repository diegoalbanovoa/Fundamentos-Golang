# Use the official Golang image as the base image
FROM golang:1.18

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port on which the server will run
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

