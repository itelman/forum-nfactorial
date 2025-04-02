# Use docker-compose inside Fly.io
FROM docker:latest

WORKDIR /app

# Copy the project files
COPY . .

# Install docker-compose
RUN apk add --no-cache docker-compose

# Start services using docker-compose
CMD ["docker-compose", "up"]
