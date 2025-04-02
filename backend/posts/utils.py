from django.db import connection


def update_post_reaction_counts(post_id):
    """Update likes and dislikes count for a specific post."""
    with connection.cursor() as cursor:
        cursor.execute("""
            UPDATE posts_post
            SET likes = (SELECT COUNT(*) FROM post_reactions_postreaction WHERE post_id = %s AND is_like = 1),
                dislikes = (SELECT COUNT(*) FROM post_reactions_postreaction WHERE post_id = %s AND is_like = 0)
            WHERE id = %s
        """, [post_id, post_id, post_id])
