from django.contrib.auth.models import AnonymousUser
from django.contrib.auth.models import User
from django.urls import resolve
from django.utils.deprecation import MiddlewareMixin
from rest_framework.exceptions import AuthenticationFailed
from rest_framework_simplejwt.exceptions import TokenError
from rest_framework_simplejwt.tokens import AccessToken


class AuthenticationMiddleware(MiddlewareMixin):
    def process_request(self, request):
        if resolve(request.path_info).app_name == "admin":
            return

        request.user = AnonymousUser()

        auth_header = request.headers.get("Authorization")
        if not auth_header:
            return

        try:
            token_str = auth_header.split(" ")[1]  # Extract token from "Bearer <token>"
            token = AccessToken(token_str)  # Decode token

            user_id = token.get("user_id")  # Extract user_id
            if not user_id:
                raise AuthenticationFailed("Invalid token: No user ID found")

            request.user = User.objects.get(id=user_id)  # Attach user to request

        except TokenError as e:
            if "expired" in str(e).lower():
                request.user = AnonymousUser()
                return
            else:
                raise AuthenticationFailed(f"Authentication error: {str(e)}")
