from django.db import connection


def update_comment_reaction_counts(comment_id):
    """Update the likes and dislikes count for a comment."""
    with connection.cursor() as cursor:
        cursor.execute("""
            UPDATE comments_comment
            SET likes = (SELECT COUNT(*) FROM comment_reactions_commentreaction WHERE comment_id = %s AND is_like = 1),
                dislikes = (SELECT COUNT(*) FROM comment_reactions_commentreaction WHERE comment_id = %s AND is_like = 0)
            WHERE id = %s
        """, [comment_id, comment_id, comment_id])
