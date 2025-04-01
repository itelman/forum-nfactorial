from django.conf import settings
from django.contrib.auth import authenticate
from rest_framework import serializers
from rest_framework import status
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response
from rest_framework.views import APIView
from rest_framework_simplejwt.tokens import RefreshToken

from posts.models import Post
from posts.serializers import PostSerializer
from .permissions import IsNotAuthenticated
from .serializers import UserSerializer

SECRET_KEY = settings.SECRET_KEY


class SignupView(APIView):
    permission_classes = [IsNotAuthenticated]

    def post(self, request):
        serializer = UserSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        serializer.save()
        return Response({"message": "User registered successfully"}, status=status.HTTP_200_OK)


class LoginView(APIView):
    permission_classes = [IsNotAuthenticated]

    def post(self, request):
        username = request.data.get("username")
        password = request.data.get("password")

        if not username or not password:
            raise serializers.ValidationError({"generic": "Username and password are required"})

        user = authenticate(username=username.lower(), password=password)
        if user is None:
            raise serializers.ValidationError(
                {"generic": "Authentication failed. Please check your credentials and try again"})

        refresh = RefreshToken.for_user(user)

        return Response({
            "access_token": str(refresh.access_token),
            "type": "Bearer",
        }, status=status.HTTP_200_OK)


class UserProfileView(APIView):
    permission_classes = [IsAuthenticated]

    def get(self, request):
        serializer = UserSerializer(request.user, context={"request": request})
        return Response(serializer.data, status=status.HTTP_200_OK)


class UserCreatedPostsView(APIView):
    permission_classes = [IsAuthenticated]

    def get(self, request):
        created_posts = Post.objects.filter(user=request.user)
        serializer = PostSerializer(created_posts, many=True, context={"request": request})
        return Response(serializer.data, status=status.HTTP_200_OK)


class UserReactedPostsView(APIView):
    permission_classes = [IsAuthenticated]

    def get(self, request):
        reacted_posts = Post.objects.filter(reactions__user=request.user)
        serializer = PostSerializer(reacted_posts, many=True, context={"request": request})
        return Response(serializer.data, status=status.HTTP_200_OK)
