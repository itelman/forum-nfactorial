from django.contrib import admin

from .models import Category


@admin.register(Category)
class CategoryAdmin(admin.ModelAdmin):
    list_display = ('id', 'name', 'created')  # Display these fields in the category list
    search_fields = ('name',)  # Enable search by category name
    list_filter = ('created',)  # Add a filter by creation date
    ordering = ('-created',)  # Order by newest first
    readonly_fields = ('created',)  # Make the created field read-only
    fieldsets = (
        ('Category Information', {
            'fields': ('name', 'created')
        }),
    )
