# Use an official Docker image that includes Docker and Docker Compose
FROM docker:latest 

# Install Docker Compose
RUN apk add --no-cache docker-compose

# Set the working directory
WORKDIR /app

# Copy your entire project
COPY . .

# Ensure Docker Compose file has correct permissions
RUN chmod +x /app/docker-compose.yml

# Expose the frontend port that Koyeb should serve
EXPOSE 8080

# Run Docker Compose when the container starts
CMD ["docker-compose", "up"]