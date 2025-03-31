from django.conf import settings
from django.contrib.auth import authenticate
from rest_framework import status
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response
from rest_framework.views import APIView
from rest_framework_simplejwt.tokens import RefreshToken

from posts.models import Post
from posts.serializers import PostSerializer
from .permissions import IsNotAuthenticated
from .serializers import UserSerializer

SECRET_KEY = settings.SECRET_KEY  # Use Django's secret key for signing JWTs


class SignupView(APIView):
    permission_classes = [IsNotAuthenticated]  # ✅ Only allow unauthenticated users

    def post(self, request):
        serializer = UserSerializer(data=request.data)
        if serializer.is_valid():
            serializer.save()
            return Response({"message": "User registered successfully"}, status=status.HTTP_200_OK)
        return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)


class LoginView(APIView):
    permission_classes = [IsNotAuthenticated]

    def post(self, request):
        username = request.data.get("username")
        password = request.data.get("password")

        if not username or not password:
            return Response({"error": "Username and password are required"}, status=status.HTTP_400_BAD_REQUEST)

        # Convert username to lowercase before authentication
        user = authenticate(username=username.lower(), password=password)
        if user is None:
            return Response({"error": "Invalid credentials"}, status=status.HTTP_401_UNAUTHORIZED)

        refresh = RefreshToken.for_user(user)

        return Response({
            'access_token': str(refresh.access_token),
            "type": "bearer",
            'refresh_token': str(refresh),
        }, status=status.HTTP_200_OK)


class UserProfileView(APIView):
    """Returns the authenticated user's profile details"""
    permission_classes = [IsAuthenticated]

    def get(self, request):
        serializer = UserSerializer(request.user, context={"request": request})
        return Response(serializer.data, status=status.HTTP_200_OK)


class UserCreatedPostsView(APIView):
    """Lists posts created by the authenticated user"""
    permission_classes = [IsAuthenticated]

    def get(self, request):
        created_posts = Post.objects.filter(user=request.user)
        serializer = PostSerializer(created_posts, many=True, context={"request": request})
        return Response(serializer.data, status=status.HTTP_200_OK)


class UserReactedPostsView(APIView):
    """Lists posts the authenticated user has liked or disliked"""
    permission_classes = [IsAuthenticated]

    def get(self, request):
        reacted_posts = Post.objects.filter(reactions__user=request.user)
        serializer = PostSerializer(reacted_posts, many=True, context={"request": request})
        return Response(serializer.data, status=status.HTTP_200_OK)
