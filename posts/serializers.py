from django.db import transaction
from rest_framework import serializers

from categories.models import Category
from comments.serializers import CommentSerializer
from post_reactions.models import PostReaction
from .models import Post
from .models import PostCategory


class PostSerializer(serializers.ModelSerializer):
    comments = CommentSerializer(many=True, read_only=True)
    auth_user_reaction = serializers.SerializerMethodField()
    categories_id = serializers.ListField(child=serializers.IntegerField(), write_only=True, required=False)

    class Meta:
        model = Post
        fields = '__all__'
        read_only_fields = ['id', 'user', 'likes', 'dislikes', 'created', 'edited']

    def create(self, validated_data):
        """Override to handle categories separately."""
        categories_id = validated_data.pop("categories_id", [])  # Extract category IDs

        with transaction.atomic():  # Ensure rollback if any step fails
            post = Post.objects.create(**validated_data)  # Create post
            categories = Category.objects.filter(id__in=categories_id)  # Retrieve categories

            # Assign categories to the post
            PostCategory.objects.bulk_create([
                PostCategory(post=post, category=category) for category in categories
            ])

        return post

    def get_auth_user_reaction(self, obj):
        """Retrieve the authenticated user's reaction to this post."""
        user = self.context.get('request').user
        if user.is_authenticated:
            reaction = PostReaction.objects.filter(post=obj, user=user).first()
            return reaction.is_like if reaction else -1
        return -1  # Return -1 if user is not authenticated

    def to_representation(self, instance):
        """Customize post serialization output"""
        data = super().to_representation(instance)

        # Apply custom formatting
        return {
            "id": data["id"],
            "user": {"id": instance.user.id, "username": instance.user.username},
            "title": data["title"],
            "content": data["content"],
            "categories": [category.name for category in instance.categories.all()],
            "likes": data["likes"],
            "dislikes": data["dislikes"],
            "created": data["created"],
            "auth_user_reaction": self.get_auth_user_reaction(instance),
            "comments": CommentSerializer(instance.comments.all().order_by('-created'), many=True, context=self.context).data,
        }
