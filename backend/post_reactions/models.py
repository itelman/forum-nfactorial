from django.contrib.auth.models import User
from django.db import models

from posts.models import Post


class PostReaction(models.Model):
    post = models.ForeignKey(Post, on_delete=models.CASCADE, related_name="reactions")
    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name="post_reactions")
    is_like = models.IntegerField(choices=[(1, "Like"), (0, "Dislike")])
    created = models.DateTimeField(auto_now_add=True)

    class Meta:
        unique_together = ("post", "user")  # Ensures one reaction per user per post

    def __str__(self):
        return f"{self.user} {'liked' if self.is_like else 'disliked'} {self.post}"
