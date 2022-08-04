package models

import (
	"database/sql"

	"github.com/appak21/forum/pkg/config"
)

var db *sql.DB

func init() {
	config.Connect()
	db = config.GetDB()
}

type User struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	CreatedAt  string `json:"createdAt"`
	Reputation int
}

type Post struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	Username  string `json:"username"`
	UserID    int64  `json:"userId"`
	CreatedAt string `json:"createdAt"`
	When      string
	Tags      []string
	Comments  []Comment
	Votes     Vote
}

type Vote struct {
	Likes    int //like numbers
	Dislikes int //dislike numbers
}

type Comment struct {
	ID        int64  `json:"id"`
	Text      string `json:"text"`
	Username  string `json:"username"`
	UserID    int64  `json:"userId"`
	PostID    int64  `json:"postId"`
	CreatedAt string `json:"createdAt"`
	When      string
	Votes     Vote
}

type Like struct {
	ID     int64 `json:"id,omitempty"`
	UserID int64 `json:"user_id,omitempty"`
	PostID int64 `json:"post_id,omitempty"`
	IsLike bool  `json:"is_like,omitempty"`
}
