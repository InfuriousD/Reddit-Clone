// models/models.go
package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a Reddit user account
type User struct {
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never sent to client
	Karma        int       `json:"karma"`
	CreatedAt    time.Time `json:"created_at"`
	Subreddits   []string  `json:"subreddits"` // List of subscribed subreddit names
}

// Subreddit represents a community
type Subreddit struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Moderators  []string  `json:"moderators"` // List of moderator usernames
	Subscribers int       `json:"subscribers"`
}

// Post represents content submitted to a subreddit
type Post struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	AuthorName    string    `json:"author_name"`
	SubredditName string    `json:"subreddit_name"`
	CreatedAt     time.Time `json:"created_at"`
	Score         int       `json:"score"`
	CommentsCount int       `json:"comments_count"`
	Upvotes       int       `json:"upvotes"`
	Downvotes     int       `json:"downvotes"`
}

// Comment represents a response to a post or another comment
type Comment struct {
	ID           uuid.UUID  `json:"id"`
	Content      string     `json:"content"`
	AuthorName   string     `json:"author_name"`
	PostID       uuid.UUID  `json:"post_id"`
	ParentID     *uuid.UUID `json:"parent_id,omitempty"` // Null for top-level comments
	CreatedAt    time.Time  `json:"created_at"`
	Score        int        `json:"score"`
	Upvotes      int        `json:"upvotes"`
	Downvotes    int        `json:"downvotes"`
	RepliesCount int        `json:"replies_count"`
}

// Vote represents a user's vote on a post or comment
type Vote struct {
	UserName  string    `json:"username"`
	TargetID  uuid.UUID `json:"target_id"` // ID of post or comment
	IsUpvote  bool      `json:"is_upvote"`
	CreatedAt time.Time `json:"created_at"`
}

// DirectMessage represents a private message between users
type DirectMessage struct {
	ID        uuid.UUID  `json:"id"`
	FromUser  string     `json:"from_user"`
	ToUser    string     `json:"to_user"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
}

// FeedResponse represents a paginated feed of posts
type FeedResponse struct {
	Posts      []Post `json:"posts"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
