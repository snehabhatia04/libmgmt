# Step 1: Use an official Golang runtime as the base image
FROM golang:1.23 AS builder

# Step 2: Set the current working directory inside the container
WORKDIR /app

# Step 3: Copy the current directory contents into the container
COPY . .

# Step 4: Install the Go dependencies
RUN ls -la /app
RUN go mod tidy

# Step 5: Expose the port your application runs on
EXPOSE 8080

# Step 6: Define the command to run your application
CMD ["go", "run", "main.go"]

