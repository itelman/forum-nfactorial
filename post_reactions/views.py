from django.db import transaction
from django.shortcuts import get_object_or_404
from rest_framework import status, permissions, viewsets
from rest_framework.response import Response

from posts.models import Post
from posts.utils import update_post_reaction_counts
from .models import PostReaction


class PostReactionViewSet(viewsets.ViewSet):
    permission_classes = [permissions.IsAuthenticated]

    def create(self, request, post_id):
        post = get_object_or_404(Post, id=post_id)
        is_like = request.data.get("is_like")
        make_insertion = True

        if is_like not in [0, 1]:
            return Response({"error": "is_like must be 0 or 1"}, status=status.HTTP_400_BAD_REQUEST)

        # Retrieve existing reaction BEFORE starting the transaction
        existing_reaction = PostReaction.objects.filter(post=post, user=request.user).first()

        with transaction.atomic():
            if existing_reaction:
                # Step 2: Delete the existing reaction
                existing_reaction.delete()

                # Step 3: Check if the old reaction was the same as the new one
                if existing_reaction.is_like == is_like:
                    make_insertion = False

            # Step 4: Insert the new reaction (only if it's different)
            if make_insertion:
                PostReaction.objects.create(post=post, user=request.user, is_like=is_like)

            # Step 5: Update post's like/dislike counts
            update_post_reaction_counts(post.id)

        if make_insertion:
            return Response({"message": "Reaction added"}, status=status.HTTP_200_OK)
        else:
            return Response({"message": "Reaction removed"}, status=status.HTTP_200_OK)
