// server/db/db.go
package db

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "github.com/google/uuid"
    "your-project/models"
)

type Database struct {
    db *sql.DB
}

func NewDatabase(connStr string) (*Database, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    return &Database{db: db}, nil
}

// User methods
func (d *Database) CreateUser(user *models.User) error {
    query := `
        INSERT INTO users (username, password_hash, created_at)
        VALUES ($1, $2, $3)
    `
    _, err := d.db.Exec(query, user.Username, user.PasswordHash, user.CreatedAt)
    return err
}

func (d *Database) GetUser(username string) (*models.User, error) {
    query := `
        SELECT username, password_hash, karma, created_at
        FROM users
        WHERE username = $1
    `
    user := &models.User{}
    err := d.db.QueryRow(query, username).Scan(
        &user.Username,
        &user.PasswordHash,
        &user.Karma,
        &user.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    return user, nil
}

// Subreddit methods
func (d *Database) CreateSubreddit(subreddit *models.Subreddit) error {
    query := `
        INSERT INTO subreddits (name, description, created_at)
        VALUES ($1, $2, $3)
    `
    _, err := d.db.Exec(query, subreddit.Name, subreddit.Description, subreddit.CreatedAt)
    return err
}

// Post methods
func (d *Database) CreatePost(post *models.Post) error {
    query := `
        INSERT INTO posts (id, title, content, author_name, subreddit_name, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
    _, err := d.db.Exec(query,
        post.ID,
        post.Title,
        post.Content,
        post.AuthorName,
        post.SubredditName,
        post.CreatedAt,
    )
    return err
}

// Comment methods
func (d *Database) CreateComment(comment *models.Comment) error {
    query := `
        INSERT INTO comments (id, content, author_name, post_id, parent_id, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
    _, err := d.db.Exec(query,
        comment.ID,
        comment.Content,
        comment.AuthorName,
        comment.PostID,
        comment.ParentID,
        comment.CreatedAt,
    )
    return err
}

// Vote methods
func (d *Database) AddVote(vote *models.Vote) error {
    query := `
        INSERT INTO votes (username, target_id, is_upvote, created_at)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (username, target_id) 
        DO UPDATE SET is_upvote = $3
    `
    _, err := d.db.Exec(query,
        vote.UserName,
        vote.TargetID,
        vote.IsUpvote,
        vote.CreatedAt,
    )
    return err
}

// Message methods
func (d *Database) CreateMessage(message *models.DirectMessage) error {
    query := `
        INSERT INTO messages (id, from_user, to_user, content, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `
    _, err := d.db.Exec(query,
        message.ID,
        message.FromUser,
        message.ToUser,
        message.Content,
        message.CreatedAt,
    )
    return err
}
