from django.urls import path

from .views import CategoryListView, CategoryPostsView

urlpatterns = [
    path('', CategoryListView.as_view(), name='category-list'),
    path('/<int:pk>/posts', CategoryPostsView.as_view(), name='category-posts'),
]
