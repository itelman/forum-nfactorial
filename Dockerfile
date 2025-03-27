# Use official Python image
FROM python:3.11

# Set working directory
WORKDIR /app

# Copy project files
COPY . .

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Collect static files
RUN python manage.py collectstatic --noinput

# Expose port (Render sets a dynamic port)
EXPOSE 8000

# Start the application (Render provides $PORT)
CMD ["gunicorn", "--bind", "0.0.0.0:$PORT", "forum.wsgi:application"]
