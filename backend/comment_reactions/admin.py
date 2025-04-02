from django.contrib import admin

from .models import CommentReaction


@admin.register(CommentReaction)
class CommentReactionAdmin(admin.ModelAdmin):
    list_display = ('id', 'comment', 'user', 'is_like', 'created')
    list_filter = ('is_like', 'created')
    search_fields = ('user__username', 'comment__content')
