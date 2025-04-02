from django.contrib.auth.models import User
from django.db import models

from comments.models import Comment


class CommentReaction(models.Model):
    comment = models.ForeignKey(Comment, on_delete=models.CASCADE, related_name="reactions")
    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name="comment_reactions")
    is_like = models.IntegerField(choices=[(1, 'Like'), (0, 'Dislike')])
    created = models.DateTimeField(auto_now_add=True)

    class Meta:
        unique_together = ('comment', 'user')

    def __str__(self):
        return f"{self.user} {'liked' if self.is_like else 'disliked'} {self.comment}"
