from django.urls import path, include

from .views import PostViewSet

post_list = PostViewSet.as_view({'get': 'list', 'post': 'create'})
post_detail = PostViewSet.as_view({'get': 'retrieve', 'put': 'update', 'delete': 'destroy'})  # Added "get": "retrieve"

urlpatterns = [
    path('', post_list, name='post-list/post-create'),  # GET & POST at "/posts"
    path('/<int:pk>', post_detail, name='post-detail/update/delete'),  # GET, PUT, DELETE at "/posts/{id}"
    path('/<int:post_id>/comments', include('comments.urls')),  # Include comments under "/posts/{post_id}/comments"
    path('/<int:post_id>/react', include('post_reactions.urls')),
]
