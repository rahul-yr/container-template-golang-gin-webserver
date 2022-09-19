FROM golang:1.19.1-alpine3.16 as builder
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy go mod and sum files
COPY go.* ./
# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY . ./
# Build the Go app
RUN go build -v -o server
# Build a small image
FROM alpine:3.16
# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/server /app/server
# Expose port 8080 to the outside world
EXPOSE 8080
# Command to run the executable
CMD ["/app/server"] 