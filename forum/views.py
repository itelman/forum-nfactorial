import platform

from rest_framework import status
from rest_framework.response import Response
from rest_framework.views import APIView


class HealthCheckView(APIView):
    """Health check endpoint to verify API availability and system info."""

    def get(self, request):
        data = {
            "status": "available",
            "system_info": {
                "environment": "development",  # Change this based on actual environment
                "version": "1.0",  # Update with actual version
                "os": platform.system(),
                "python_version": platform.python_version(),
            }
        }
        return Response(data, status=status.HTTP_200_OK)
