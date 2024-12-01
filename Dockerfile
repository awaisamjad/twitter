# Use the latest Golang base image
FROM golang:latest AS builder

# Set the work directory in the container
WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./src/main.go

# Start a new build stage
FROM alpine:latest  

# Install certificates
RUN apk --no-cache add ca-certificates

# Set work directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy templates directory
COPY --from=builder /app/templates ./templates

# Expose port 3000
EXPOSE 3000

# Command to run the executable
CMD ["./main"]