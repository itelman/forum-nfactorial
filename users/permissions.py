from rest_framework.permissions import BasePermission


class IsNotAuthenticated(BasePermission):
    """
    Custom permission to allow access only to unauthenticated users.
    """

    def has_permission(self, request, view):
        return not request.user or not request.user.is_authenticated


class IsOwner(BasePermission):
    """
    Custom permission to allow only the owner of a post to edit or delete it.
    """

    def has_object_permission(self, request, view, obj):
        return obj.user == request.user  # Allow only if user owns the post
