from django.urls import path

from .views import PostReactionViewSet

urlpatterns = [
    path('', PostReactionViewSet.as_view({'post': 'create'}), name='post-react'),
]
