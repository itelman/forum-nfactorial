package adapters

import (
	"database/sql"
	"errors"

	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/post_reactions/domain"
)

type PostsRepositorySqlite struct {
	db *sql.DB
}

func NewPostsRepositorySqlite(db *sql.DB) *PostsRepositorySqlite {
	return &PostsRepositorySqlite{db}
}

func (r *PostsRepositorySqlite) Get(input domain.GetPostInput) (*dto.Post, error) {
	query := "SELECT posts.id, users.id, users.username, posts.title, posts.content, posts.likes, posts.dislikes, posts.created FROM posts INNER JOIN users ON posts.user_id = users.id WHERE posts.id = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	post := &dto.Post{User: &dto.User{}}
	if err := stmt.QueryRow(input.ID).Scan(
		&post.ID,
		&post.User.ID,
		&post.User.Username,
		&post.Title,
		&post.Content,
		&post.Likes,
		&post.Dislikes,
		&post.Created,
	); errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrPostNotFound
	} else if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *PostsRepositorySqlite) UpdateReactionsCount(tx *sql.Tx, input domain.UpdatePostReactionsCountInput) error {
	query := "UPDATE posts SET likes = (SELECT COUNT(*) FROM post_reactions WHERE post_id = posts.id AND is_like = 1), dislikes = (SELECT COUNT(*) FROM post_reactions WHERE post_id = posts.id AND is_like = 0) WHERE id = ?"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(input.PostID); err != nil {
		return err
	}

	return nil
}
