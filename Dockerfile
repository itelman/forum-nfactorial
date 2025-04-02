# Use docker-compose inside Fly.io
FROM docker:latest

WORKDIR /app

# Copy the project files
COPY . .

# Install docker-compose
RUN apk add --no-cache docker-compose

# Expose the frontend port that Koyeb should serve
EXPOSE 8080

# Start services using docker-compose
CMD ["docker-compose", "up"]