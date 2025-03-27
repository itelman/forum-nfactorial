# Use official Python image as a base
FROM python:3.11

# Set working directory in the container
WORKDIR /app

# Copy project files
COPY . .

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Expose port (Render sets a dynamic port)
CMD ["sh", "-c", "python manage.py migrate && python manage.py collectstatic --noinput && gunicorn forum.wsgi:application --bind 0.0.0.0:$PORT"]
