# Use an official Golang image to build our application
FROM golang:1.21-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
# This caches the dependencies layer in Docker
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o /goat-monitor

# Expose the port that our app will listen on
EXPOSE 8080

# The command to run when the container starts
CMD [ "/goat-monitor" ]