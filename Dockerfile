# Use an official Go runtime as a parent image
FROM golang:1.22.5-alpine

# Install necessary build tools and SQLite
RUN apk add --no-cache gcc musl-dev sqlite

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the working directory inside the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o main ./cmd/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
