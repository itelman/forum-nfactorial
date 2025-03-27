from rest_framework.generics import ListAPIView, get_object_or_404
from rest_framework.response import Response
from rest_framework.views import APIView

from posts.models import Post
from posts.serializers import PostSerializer
from .models import Category
from .serializers import CategorySerializer


# ✅ List all categories
class CategoryListView(ListAPIView):
    queryset = Category.objects.all()
    serializer_class = CategorySerializer


# ✅ List all posts in a specific category
class CategoryPostsView(APIView):
    def get(self, request, pk):
        category = get_object_or_404(Category, id=pk)
        posts = Post.objects.filter(postcategory__category=category)
        return Response(PostSerializer(posts, many=True, context={'request': request}).data)
