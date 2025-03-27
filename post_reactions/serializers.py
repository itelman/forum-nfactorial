from rest_framework import serializers

from .models import PostReaction


class PostReactionSerializer(serializers.ModelSerializer):
    class Meta:
        model = PostReaction
        fields = ['is_like']
