from django.contrib import admin

from .models import PostReaction


@admin.register(PostReaction)
class PostReactionAdmin(admin.ModelAdmin):
    list_display = ("id", "post", "user", "is_like", "created")
    list_filter = ("is_like", "created")
    search_fields = ("user__username", "post__title")
    ordering = ("-created",)
