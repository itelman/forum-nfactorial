from django.contrib.auth.models import User
from django.db import models

from categories.models import Category


# Create your models here.
class Post(models.Model):
    user = models.ForeignKey(User, on_delete=models.CASCADE)
    title = models.TextField()
    content = models.TextField()
    likes = models.IntegerField(default=0)
    dislikes = models.IntegerField(default=0)
    created = models.DateTimeField(auto_now_add=True)
    edited = models.DateTimeField(null=True, blank=True)

    # Many-to-Many with Category through PostCategory
    categories = models.ManyToManyField(Category, related_name='posts', through='PostCategory')

    def __str__(self):
        return self.title


class PostCategory(models.Model):
    post = models.ForeignKey(Post, on_delete=models.CASCADE)
    category = models.ForeignKey(Category, on_delete=models.CASCADE)
    created = models.DateTimeField(auto_now_add=True)

    class Meta:
        unique_together = ('post', 'category')  # Ensure a post-category pair is unique
