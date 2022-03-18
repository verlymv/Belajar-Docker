#pake linux yang udah diinstall golang
FROM golang:1.16-alpine AS builder 

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Add a work directory and move to working directory /app
WORKDIR /app

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o app

# Build a small image
FROM alpine:3.14

# Copy built binary from builder
COPY --from=builder app .

# Expose port
EXPOSE 8081

# Command to run
CMD ./app
