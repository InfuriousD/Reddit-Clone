// server/engine/engine.go
package engine

import (
	"context"
	"reddit-clone/models"
	"reddit-clone/server/db"
	"time"

	"github.com/google/uuid"
)

type RedditEngine struct {
	db *db.Database
}

func NewRedditEngine(db *db.Database) *RedditEngine {
	return &RedditEngine{db: db}
}

// User operations
func (e *RedditEngine) CreateUser(ctx context.Context, username, passwordHash string) error {
	user := &models.User{
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}
	return e.db.CreateUser(user)
}

// Subreddit operations
func (e *RedditEngine) CreateSubreddit(ctx context.Context, name, description, creator string) error {
	subreddit := &models.Subreddit{
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		Moderators:  []string{creator},
	}
	return e.db.CreateSubreddit(subreddit)
}

// Post operations
func (e *RedditEngine) CreatePost(ctx context.Context, title, content, authorName, subredditName string) error {
	post := &models.Post{
		ID:            uuid.New(),
		Title:         title,
		Content:       content,
		AuthorName:    authorName,
		SubredditName: subredditName,
		CreatedAt:     time.Now(),
	}
	return e.db.CreatePost(post)
}

// Comment operations
func (e *RedditEngine) CreateComment(ctx context.Context, content, authorName string, postID uuid.UUID, parentID *uuid.UUID) error {
	comment := &models.Comment{
		ID:         uuid.New(),
		Content:    content,
		AuthorName: authorName,
		PostID:     postID,
		ParentID:   parentID,
		CreatedAt:  time.Now(),
	}
	return e.db.CreateComment(comment)
}

// Vote operations
func (e *RedditEngine) Vote(ctx context.Context, username string, targetID uuid.UUID, isUpvote bool) error {
	vote := &models.Vote{
		UserName:  username,
		TargetID:  targetID,
		IsUpvote:  isUpvote,
		CreatedAt: time.Now(),
	}
	return e.db.AddVote(vote)
}

// Message operations
func (e *RedditEngine) SendMessage(ctx context.Context, fromUser, toUser, content string) error {
	message := &models.DirectMessage{
		ID:        uuid.New(),
		FromUser:  fromUser,
		ToUser:    toUser,
		Content:   content,
		CreatedAt: time.Now(),
	}
	return e.db.CreateMessage(message)
}

// Feed operations
func (e *RedditEngine) GetFeed(ctx context.Context, username string, offset, limit int) ([]*models.Post, error) {
	// Get user's subscribed subreddits
	user, err := e.db.GetUser(username)
	if err != nil {
		return nil, err
	}

	// Get posts from subscribed subreddits
	return e.db.GetFeedPosts(user.Subreddits, offset, limit)
}
