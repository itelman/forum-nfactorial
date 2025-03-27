from django.urls import path, include

from .views import CommentViewSet

comment_list = CommentViewSet.as_view({'get': 'list', 'post': 'create'})
comment_detail = CommentViewSet.as_view({'get': 'retrieve', 'put': 'update', 'delete': 'destroy'})

urlpatterns = [
    path('', comment_list, name='comment-list/comment-create'),  # GET & POST comments for a post
    path('/<int:comment_id>', comment_detail, name='comment-detail/update/delete'),
    path('/<int:comment_id>/react', include('comment_reactions.urls')),
]
