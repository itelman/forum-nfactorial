from django.urls import path

from .views import CommentReactionViewSet

urlpatterns = [
    path('', CommentReactionViewSet.as_view({'post': 'create'}), name='comment-react'),
]
