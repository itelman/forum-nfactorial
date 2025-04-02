from django.db import transaction
from django.shortcuts import get_object_or_404
from rest_framework import permissions, status, viewsets
from rest_framework.response import Response

from comments.models import Comment
from comments.utils import update_comment_reaction_counts
from .models import CommentReaction


class CommentReactionViewSet(viewsets.ViewSet):
    permission_classes = [permissions.IsAuthenticated]

    def create(self, request, post_id, comment_id):
        comment = get_object_or_404(Comment, id=comment_id, post_id=post_id)
        is_like = request.data.get("is_like")
        make_insertion = True

        if is_like not in [0, 1]:
            return Response({"error": "is_like must be 0 or 1"}, status=status.HTTP_400_BAD_REQUEST)

        existing_reaction = CommentReaction.objects.filter(comment=comment, user=request.user).first()

        with transaction.atomic():
            if existing_reaction:
                existing_reaction.delete()

                if existing_reaction.is_like == is_like:
                    make_insertion = False

            if make_insertion:
                CommentReaction.objects.create(comment=comment, user=request.user, is_like=is_like)

            update_comment_reaction_counts(comment.id)

        if make_insertion:
            return Response({"message": "Reaction added"}, status=status.HTTP_200_OK)
        else:
            return Response({"message": "Reaction removed"}, status=status.HTTP_200_OK)
