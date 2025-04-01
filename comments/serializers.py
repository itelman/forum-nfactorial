from rest_framework import serializers

from comment_reactions.models import CommentReaction
from .models import Comment


class CommentSerializer(serializers.ModelSerializer):
    auth_user_reaction = serializers.SerializerMethodField()

    class Meta:
        model = Comment
        fields = '__all__'
        read_only_fields = ['id', 'user', 'likes', 'dislikes', 'created', 'edited', 'post']

    def validate_content(self, value):
        if value != value.strip():
            raise serializers.ValidationError("Please provide a valid content.")
        return value.strip()

    def create(self, validated_data):
        post = self.context.get('post')
        return Comment.objects.create(post=post, **validated_data)

    def get_auth_user_reaction(self, obj):
        user = self.context.get('request').user
        if user.is_authenticated:
            reaction = CommentReaction.objects.filter(comment=obj, user=user).first()
            return reaction.is_like if reaction else -1
        return -1

    def to_representation(self, instance):
        data = super().to_representation(instance)

        return {
            "id": data["id"],
            "post_id": instance.post.id,
            "user": {"id": instance.user.id, "username": instance.user.username},
            "content": data["content"],
            "likes": data["likes"],
            "dislikes": data["dislikes"],
            "created": data["created"],
            "auth_user_reaction": self.get_auth_user_reaction(instance),
        }
