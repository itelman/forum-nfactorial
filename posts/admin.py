# Register your models here.
from django.contrib import admin

from .models import Post


@admin.register(Post)
class PostAdmin(admin.ModelAdmin):
    list_display = ("id", "user", "title", "content", "likes", "dislikes", "created", "edited", "display_categories")
    list_filter = ("created", "edited", "categories")
    search_fields = ("title", "content", "user__username")
    ordering = ("-created",)

    def display_categories(self, obj):
        return ", ".join(category.name for category in obj.categories.all())

    display_categories.short_description = "Categories"  # Custom column header
