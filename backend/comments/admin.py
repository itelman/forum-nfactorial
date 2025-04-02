from django.contrib import admin

from .models import Comment


@admin.register(Comment)
class CommentAdmin(admin.ModelAdmin):
    list_display = ('id', 'post', 'user', 'content', 'likes', 'dislikes', 'created', 'edited')
    list_filter = ('created', 'edited')
    search_fields = ('content', 'user__username', 'post__title')
    ordering = ('-created',)
