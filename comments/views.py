from datetime import datetime

from django.shortcuts import get_object_or_404
from rest_framework import serializers
from rest_framework import viewsets, permissions, status
from rest_framework.response import Response

from posts.models import Post
from .models import Comment
from .serializers import CommentSerializer


class CommentViewSet(viewsets.ViewSet):
    permission_classes = [permissions.IsAuthenticatedOrReadOnly]

    def list(self, request, post_id):
        post = get_object_or_404(Post, id=post_id)
        comments = post.comments.all().order_by('-created')
        serializer = CommentSerializer(comments, many=True, context={'request': request, 'post': post})
        return Response(serializer.data)

    def create(self, request, post_id):
        post = get_object_or_404(Post, id=post_id)
        serializer = CommentSerializer(data=request.data, context={'request': request, 'post': post})
        serializer.is_valid(raise_exception=True)
        comment = serializer.save(user=request.user)

        return Response(CommentSerializer(comment, context={'request': request}).data,
                        status=status.HTTP_200_OK)

    def retrieve(self, request, post_id, comment_id):
        comment = get_object_or_404(Comment, id=comment_id, post_id=post_id)
        serializer = CommentSerializer(comment, context={'request': request, 'post': comment.post})
        return Response(serializer.data, status=status.HTTP_200_OK)

    def update(self, request, post_id, comment_id):
        comment = get_object_or_404(Comment, id=comment_id, post_id=post_id)

        if comment.user != request.user:
            return Response({"error": "You can only edit your own comments."}, status=status.HTTP_403_FORBIDDEN)

        new_content = request.data.get("content", None)
        if new_content:
            serializer = CommentSerializer(comment, data={"content": new_content}, partial=True,
                                           context={'request': request, 'post': comment.post})
            serializer.is_valid(raise_exception=True)

            if new_content == comment.content:
                raise serializers.ValidationError({"content": "You haven't made any changes"})

        serializer = CommentSerializer(comment, data=request.data, partial=True,
                                       context={'request': request, 'post': comment.post})
        serializer.is_valid(raise_exception=True)
        comment.edited = datetime.now()
        serializer.save()
        
        return Response(CommentSerializer(comment, context={'request': request}).data, status=status.HTTP_200_OK)

    def destroy(self, request, post_id, comment_id):
        comment = get_object_or_404(Comment, id=comment_id, post_id=post_id)

        if comment.user != request.user:
            return Response({"error": "You can only delete your own comments."}, status=status.HTTP_403_FORBIDDEN)

        comment.delete()
        return Response(status=status.HTTP_200_OK)
