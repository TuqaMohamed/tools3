# Use an official Go runtime as a base image
FROM golang:latest

# Set the working directory
WORKDIR /phase1

# Copy the local package files to the container's workspace
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the entire project to the container's workspace
COPY . .

# Build the Go application
RUN go build -o main .

# Expose port 8080 to the outside world
ARG PORT_ARG=8080
ENV PORT=${PORT_ARG}
EXPOSE $PORT

# Command to run the executable
CMD ["./main"]