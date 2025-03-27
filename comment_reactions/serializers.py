from rest_framework import serializers

from .models import CommentReaction


class CommentReactionSerializer(serializers.ModelSerializer):
    class Meta:
        model = CommentReaction
        fields = ['id', 'comment', 'user', 'is_like', 'created']
        read_only_fields = ['id', 'user', 'created']
