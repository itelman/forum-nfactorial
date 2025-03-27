# Register your models here.
from django.contrib import admin

from .models import Post


@admin.register(Post)
class PostAdmin(admin.ModelAdmin):
    list_display = ("id", "title", "user", "likes", "dislikes", "created", "edited")
    list_filter = ("created", "edited")
    search_fields = ("title", "content", "user__username")
    ordering = ("-created",)
