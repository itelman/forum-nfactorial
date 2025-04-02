from django.urls import path

from .views import SignupView, LoginView, UserCreatedPostsView, UserReactedPostsView, UserProfileView

urlpatterns = [
    path("/signup", SignupView.as_view(), name="user-signup"),
    path("/login", LoginView.as_view(), name="user-login"),
    path("/posts/created", UserCreatedPostsView.as_view(), name="user-created-posts"),
    path("/posts/reacted", UserReactedPostsView.as_view(), name="user-reacted-posts"),
    path("/me", UserProfileView.as_view(), name="user-profile"),
]
