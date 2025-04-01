from datetime import datetime

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
    categories = serializers.ListField(child=serializers.CharField(), write_only=True, required=False)

    class Meta:
        model = Post
        fields = '__all__'
        read_only_fields = ['id', 'user', 'likes', 'dislikes', 'created', 'edited']

    def validate_title(self, value):
        if value != value.strip():
            raise serializers.ValidationError("Please enter a valid title.")

        if not (5 <= len(value.strip()) <= 50):
            raise serializers.ValidationError("Title must be between 5 and 50 characters long.")

        return value.strip()

    def validate_content(self, value):
        if value != value.strip():
            raise serializers.ValidationError("Please enter a valid content.")

        return value.strip()

    def validate_categories(self, value):
        try:
            category_ids = [int(cat_id) for cat_id in value]
        except ValueError:
            raise serializers.ValidationError("Please enter valid categories.")

        existing_categories = set(Category.objects.filter(id__in=category_ids).values_list("id", flat=True))
        if len(existing_categories) != len(category_ids):
            raise serializers.ValidationError("Please enter valid categories.")

        return category_ids

    def create(self, validated_data):
        categories_id = validated_data.pop("categories", [])

        with transaction.atomic():
            post = Post.objects.create(**validated_data)
            categories = Category.objects.filter(id__in=categories_id)

            PostCategory.objects.bulk_create([
                PostCategory(post=post, category=category) for category in categories
            ])

        return post

    def update(self, instance, validated_data):
        if 'title' in validated_data:
            validated_data['title'] = self.validate_title(validated_data['title'])
        if 'content' in validated_data:
            validated_data['content'] = self.validate_content(validated_data['content'])

        if validated_data['title'] == instance.title and validated_data['content'] == instance.content:
            raise serializers.ValidationError({"generic": "You haven't made any changes"})

        instance.edited = datetime.now()

        for attr, value in validated_data.items():
            setattr(instance, attr, value)

        instance.save()
        return instance

    def get_auth_user_reaction(self, obj):
        user = self.context.get('request').user
        if user.is_authenticated:
            reaction = PostReaction.objects.filter(post=obj, user=user).first()
            return reaction.is_like if reaction else -1
        return -1

    def to_representation(self, instance):
        data = super().to_representation(instance)

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
            "comments": CommentSerializer(instance.comments.all().order_by('-created'), many=True,
                                          context=self.context).data,
        }
