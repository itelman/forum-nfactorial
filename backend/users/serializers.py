import re

from django.contrib.auth.models import User
from django.core.exceptions import ValidationError
from django.core.validators import validate_email
from rest_framework import serializers


class UserSerializer(serializers.ModelSerializer):
    class Meta:
        model = User
        fields = ["id", "username", "email", "date_joined", "password"]
        extra_kwargs = {
            "password": {"write_only": True},
            "email": {"required": True},
            "username": {"required": True},
        }

    def validate_username(self, value):
        username_regex = r"^[a-zA-Z]{5,}([._]{0,1}[a-zA-Z0-9]{2,})*$"

        if not re.fullmatch(username_regex, value):
            raise serializers.ValidationError("Please enter a valid username.")

        if not (5 <= len(value) <= 30):
            raise serializers.ValidationError("Username must be between 5 and 30 characters long.")

        if User.objects.filter(username__iexact=value).exists():
            raise serializers.ValidationError("This username is already taken.")

        return value.lower()

    def validate_email(self, value):
        value = value.strip().lower()

        try:
            validate_email(value)
        except ValidationError:
            raise serializers.ValidationError("Please enter a valid email address.")

        if User.objects.filter(email__iexact=value).exists():
            raise serializers.ValidationError("This email is already registered.")

        return value

    def validate_password(self, value):
        password_regex = r"^[a-zA-Z0-9_.-]+$"

        if not re.fullmatch(password_regex, value):
            raise serializers.ValidationError(
                "Please enter a valid password. Only letters, numbers, '_', '.', and '-' are allowed."
            )

        if not (6 <= len(value) <= 20):
            raise serializers.ValidationError("Password must be between 6 and 20 characters long.")

        return value

    def create(self, validated_data):
        password = validated_data.pop("password")
        validated_data["username"] = validated_data["username"].lower()
        validated_data["email"] = validated_data["email"].lower()

        user = User(**validated_data)
        user.set_password(password)
        user.save()
        return user
