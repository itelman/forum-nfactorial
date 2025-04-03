# Start from an official image with Docker installed
FROM docker:latest

# Install Docker Compose
RUN apk add --no-cache docker-compose

# Set the working directory
WORKDIR /app

# Copy project files
COPY . .

EXPOSE 8080

# Run Docker Compose
CMD ["docker-compose", "up", "--build"]