from django.contrib.auth.models import User
from django.db import models

from posts.models import Post


class Comment(models.Model):
    post = models.ForeignKey(Post, on_delete=models.CASCADE, related_name='comments')  # Add related_name
    user = models.ForeignKey(User, on_delete=models.CASCADE)
    content = models.TextField()
    likes = models.IntegerField(default=0)
    dislikes = models.IntegerField(default=0)
    created = models.DateTimeField(auto_now_add=True)
    edited = models.DateTimeField(null=True, blank=True)

    def __str__(self):
        return f"Comment by {self.user.username} on {self.post.id}"
