# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Create a directory for your "uploads" files and set it as a volume
VOLUME /app/uploads

# Copy the Go source code and any necessary files into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose port 8080
EXPOSE 8080

# Define the command to run your Go server when the container starts
CMD ["./main"]
