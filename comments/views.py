from django.shortcuts import get_object_or_404
from rest_framework import viewsets, permissions, status
from rest_framework.response import Response

from posts.models import Post
from .models import Comment
from .serializers import CommentSerializer


class CommentViewSet(viewsets.ViewSet):
    permission_classes = [permissions.IsAuthenticatedOrReadOnly]

    def list(self, request, post_id):
        """Retrieve all comments under a specific post."""
        post = get_object_or_404(Post, id=post_id)
        comments = post.comments.all()
        serializer = CommentSerializer(comments, many=True, context={'request': request})
        return Response(serializer.data)

    def create(self, request, post_id):
        """Create a comment under a post."""
        post = get_object_or_404(Post, id=post_id)
        serializer = CommentSerializer(data=request.data, context={'request': request})

        if serializer.is_valid():
            comment = serializer.save(user=request.user, post=post)
            return Response(CommentSerializer(comment, context={'request': request}).data,
                            status=status.HTTP_200_OK)

        return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)

    def retrieve(self, request, post_id, comment_id):
        """Retrieve a specific comment under a post."""
        comment = get_object_or_404(Comment, id=comment_id, post_id=post_id)
        serializer = CommentSerializer(comment, context={'request': request})
        return Response(serializer.data, status=status.HTTP_200_OK)

    def update(self, request, post_id, comment_id):
        """Update a comment. Only the owner can edit."""
        comment = get_object_or_404(Comment, id=comment_id, post_id=post_id)

        if comment.user != request.user:
            return Response({"error": "You can only edit your own comments."}, status=status.HTTP_403_FORBIDDEN)

        serializer = CommentSerializer(comment, data=request.data, partial=True, context={'request': request})
        if serializer.is_valid():
            serializer.save()
            return Response(CommentSerializer(comment, context={'request': request}).data, status=status.HTTP_200_OK)

        return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)

    def destroy(self, request, post_id, comment_id):
        """Delete a comment. Only the owner can delete."""
        comment = get_object_or_404(Comment, id=comment_id, post_id=post_id)

        if comment.user != request.user:
            return Response({"error": "You can only delete your own comments."}, status=status.HTTP_403_FORBIDDEN)

        comment.delete()
        return Response(status=status.HTTP_200_OK)
