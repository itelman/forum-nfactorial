from django.contrib.auth.models import User
from rest_framework import serializers


class UserSerializer(serializers.ModelSerializer):
    class Meta:
        model = User
        fields = ["id", "username", "email", "date_joined", "password"]
        extra_kwargs = {
            "password": {"write_only": True},  # Ensure password is not exposed
            "email": {"required": True},
            "username": {"required": True},
        }

    def validate_username(self, value):
        """Ensure username is unique (case-insensitive)"""
        if User.objects.filter(username__iexact=value).exists():
            raise serializers.ValidationError("This username is already taken.")
        return value.lower()

    def validate_email(self, value):
        """Ensure email is unique (case-insensitive)"""
        if User.objects.filter(email__iexact=value).exists():
            raise serializers.ValidationError("This email is already registered.")
        return value.lower()

    def create(self, validated_data):
        """Save the user with a properly hashed password"""
        password = validated_data.pop("password")  # Extract password before creating the user
        validated_data["username"] = validated_data["username"].lower()
        validated_data["email"] = validated_data["email"].lower()

        user = User(**validated_data)
        user.set_password(password)  # ✅ Hashes the password before saving
        user.save()
        return user
