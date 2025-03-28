package dto

import (
	"time"
)

type User struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Created  time.Time `json:"date_joined"`
}

type Post struct {
	ID               int        `json:"id"`
	User             *User      `json:"user"`
	Title            string     `json:"title"`
	Content          string     `json:"content"`
	Categories       []string   `json:"categories"`
	Likes            int        `json:"likes"`
	Dislikes         int        `json:"dislikes"`
	Created          time.Time  `json:"created"`
	AuthUserReaction int        `json:"auth_user_reaction"`
	Comments         []*Comment `json:"comments"`
}

type Comment struct {
	ID               int       `json:"id"`
	PostID           int       `json:"post_id"`
	User             *User     `json:"user"`
	Content          string    `json:"content"`
	Likes            int       `json:"likes"`
	Dislikes         int       `json:"dislikes"`
	Created          time.Time `json:"created"`
	AuthUserReaction int       `json:"auth_user_reaction"`
}

type Category struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

type PostReaction struct {
	ID      int
	PostID  int
	User    *User
	IsLike  int
	Created time.Time
}

type CommentReaction struct {
	ID        int
	CommentID int
	User      *User
	IsLike    int
	Created   time.Time
}
