# Start from an official image with Docker installed
FROM docker:20.10.24-dind

# Install Docker Compose
RUN apk add --no-cache docker-compose

# Set the working directory
WORKDIR /app

# Copy project files
COPY . .

# Run Docker Compose
CMD ["docker-compose", "up", "--build"]