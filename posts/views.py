# Create your views here.

from rest_framework import viewsets
from rest_framework.permissions import AllowAny, IsAuthenticated

from users.permissions import IsOwner
from .models import Post
from .serializers import PostSerializer


class PostViewSet(viewsets.ModelViewSet):
    queryset = Post.objects.all()
    serializer_class = PostSerializer

    def get_permissions(self):
        """
        Apply different permissions for different actions.
        """
        if self.action in ['update', 'partial_update', 'destroy']:
            return [IsAuthenticated(), IsOwner()]  # Require authentication & ownership
        return [AllowAny()]  # Other actions are public

    def perform_create(self, serializer):
        serializer.save(user=self.request.user)

    def get_serializer_context(self):
        """Pass the request to the serializer context to access the authenticated user."""
        context = super().get_serializer_context()
        context["request"] = self.request
        return context
