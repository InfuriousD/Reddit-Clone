package engine

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// Engine struct holds all subreddits, users, and messages
type Engine struct {
	subreddits map[string]*Subreddit
	users      map[string]*User
	messages   map[string][]*DirectMessage
	mu         sync.Mutex
}

type Subreddit struct {
	Name    string
	Posts   []*Post
	Members map[string]*User
}

type User struct {
	Username  string
	Karma     int
	Posts     []*Post
	Comments  []*Comment
	Connected bool
}

type Post struct {
	ID           string
	Author       *User
	Subreddit    *Subreddit
	Content      string
	Timestamp    time.Time
	Upvotes      int
	Downvotes    int
	Comments     []*Comment
	IsRepost     bool
	OriginalPost *Post
}

type Comment struct {
	ID        string
	Author    *User
	Content   string
	Timestamp time.Time
	Parent    *Post
	ReplyTo   *Comment
	Replies   []*Comment
	Upvotes   int
	Downvotes int
}

type DirectMessage struct {
	ID        string
	From      *User
	To        *User
	Content   string
	Timestamp time.Time
	ReplyTo   string
}

func (e *Engine) LeaveSubreddit(username, subredditName string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	subreddit, exists := e.subreddits[subredditName]
	if !exists {
		return
	}
	delete(subreddit.Members, username)
}

func (e *Engine) GetFeed(username string) []*Post {
	e.mu.Lock()
	defer e.mu.Unlock()

	var feed []*Post
	for _, subreddit := range e.subreddits {
		if _, isMember := subreddit.Members[username]; isMember {
			feed = append(feed, subreddit.Posts...)
		}
	}
	return feed
}

func (e *Engine) ReplyToComment(postID, parentCommentID, username, content string) *Comment {
	e.mu.Lock()
	defer e.mu.Unlock()

	user, exists := e.users[username]
	if !exists {
		return nil
	}

	var parentComment *Comment
	var parentPost *Post

	// Find the parent post and comment
	for _, subreddit := range e.subreddits {
		for _, post := range subreddit.Posts {
			if post.ID == postID {
				parentPost = post
				for _, comment := range post.Comments {
					if comment.ID == parentCommentID {
						parentComment = comment
						break
					}
				}
				break
			}
		}
	}

	if parentComment == nil || parentPost == nil {
		return nil
	}

	reply := &Comment{
		ID:        fmt.Sprintf("comment-%d", rand.Int()),
		Author:    user,
		Content:   content,
		Timestamp: time.Now(),
		Parent:    parentPost,
		ReplyTo:   parentComment,
	}

	parentComment.Replies = append(parentComment.Replies, reply)
	user.Comments = append(user.Comments, reply)
	return reply
}

func (e *Engine) Repost(originalPostID, username, subredditName string) *Post {
	e.mu.Lock()
	defer e.mu.Unlock()

	var originalPost *Post
	for _, subreddit := range e.subreddits {
		for _, post := range subreddit.Posts {
			if post.ID == originalPostID {
				originalPost = post
				break
			}
		}
	}

	if originalPost == nil {
		return nil
	}

	user := e.users[username]
	subreddit := e.subreddits[subredditName]

	repost := &Post{
		ID:           fmt.Sprintf("post-%d", rand.Int()),
		Author:       user,
		Content:      originalPost.Content,
		Subreddit:    subreddit,
		Timestamp:    time.Now(),
		IsRepost:     true,
		OriginalPost: originalPost,
	}

	subreddit.Posts = append(subreddit.Posts, repost)
	user.Posts = append(user.Posts, repost)
	return repost
}

func (e *Engine) SetUserConnection(username string, connected bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if user, exists := e.users[username]; exists {
		user.Connected = connected
	}
}

// NewEngine creates and initializes a new Reddit-like engine
func NewEngine() *Engine {
	return &Engine{
		subreddits: make(map[string]*Subreddit),
		users:      make(map[string]*User),
		messages:   make(map[string][]*DirectMessage),
	}
}

// CreateSubreddit creates a new subreddit if it doesn't already exist
func (e *Engine) CreateSubreddit(name string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.subreddits[name]; exists {
		fmt.Println("Subreddit already exists!")
		return
	}
	e.subreddits[name] = &Subreddit{Name: name, Members: make(map[string]*User)}
	fmt.Printf("Subreddit %s created.\n", name)
}

// RegisterAccount registers a new user in the engine
func (e *Engine) RegisterAccount(username string) *User {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.users[username]; exists {
		fmt.Println("Account already exists!")
		return nil
	}
	user := &User{Username: username, Karma: 0}
	e.users[username] = user
	fmt.Printf("User %s registered.\n", username)
	return user
}

// JoinSubreddit allows a user to join a subreddit
func (e *Engine) JoinSubreddit(username, subredditName string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	user, exists := e.users[username]
	if !exists {
		fmt.Println("User not found!")
		return
	}

	subreddit, exists := e.subreddits[subredditName]
	if !exists {
		fmt.Println("Subreddit not found!")
		return
	}

	subreddit.Members[username] = user
	fmt.Printf("%s joined the subreddit %s.\n", username, subredditName)
}

// PostInSubreddit allows a user to post in a subreddit
func (e *Engine) PostInSubreddit(subredditName, username, content string) *Post {
	e.mu.Lock()
	defer e.mu.Unlock()

	user, exists := e.users[username]
	if !exists {
		fmt.Println("User not found!")
		return nil
	}

	subreddit, exists := e.subreddits[subredditName]
	if !exists {
		fmt.Println("Subreddit not found!")
		return nil
	}

	post := &Post{
		ID:        fmt.Sprintf("%d", rand.Int()), // Generate a random post ID
		Author:    user,
		Subreddit: subreddit,
		Content:   content,
		Upvotes:   0,
		Downvotes: 0,
	}
	subreddit.Posts = append(subreddit.Posts, post)
	user.Posts = append(user.Posts, post)

	fmt.Printf("%s posted in %s: %s\n", username, subredditName, content)
	return post
}

// CommentOnPost allows a user to comment on a post
func (e *Engine) CommentOnPost(subredditName string, postID string, username string, content string) *Comment {
	e.mu.Lock()
	defer e.mu.Unlock()

	user, exists := e.users[username]
	if !exists {
		fmt.Println("User not found!")
		return nil
	}

	subreddit, exists := e.subreddits[subredditName]
	if !exists {
		fmt.Println("Subreddit not found!")
		return nil
	}

	var post *Post
	for _, p := range subreddit.Posts {
		if p.ID == postID {
			post = p
			break
		}
	}

	if post == nil {
		fmt.Println("Post not found!")
		return nil
	}

	comment := &Comment{
		Author:  user,
		Content: content,
		Parent:  post,
	}
	post.Comments = append(post.Comments, comment)
	user.Comments = append(user.Comments, comment)

	fmt.Printf("%s commented on post %s: %s\n", username, postID, content)
	return comment
}

// Upvote - Upvote or downvote a post or comment
// Upvote - Upvote or downvote a post or comment
func (e *Engine) Upvote(postID string, upvote bool, username string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	user, exists := e.users[username]
	if !exists {
		fmt.Println("User not found!")
		return
	}

	var post *Post
	for _, subreddit := range e.subreddits {
		for _, p := range subreddit.Posts {
			if p.ID == postID {
				post = p
				break
			}
		}
	}
	if post != nil {
		if upvote {
			post.Upvotes++
			post.Author.Karma++ // Update post author's karma
			user.Karma++        // Update voting user's karma
		} else {
			post.Downvotes++
			post.Author.Karma-- // Update post author's karma
			user.Karma--        // Update voting user's karma
		}
		fmt.Printf("%s voted on post %s. Upvotes: %d, Downvotes: %d\n", user.Username, postID, post.Upvotes, post.Downvotes)
		return
	}

	// Check for comment
	var comment *Comment
	for _, subreddit := range e.subreddits {
		for _, p := range subreddit.Posts {
			for _, c := range p.Comments {
				if c.Parent.ID == postID {
					comment = c
					break
				}
			}
		}
	}
	if comment != nil {
		if upvote {
			comment.Upvotes++
			comment.Author.Karma++ // Update comment author's karma
			user.Karma++           // Update voting user's karma
		} else {
			comment.Downvotes++
			comment.Author.Karma-- // Update comment author's karma
			user.Karma--           // Update voting user's karma
		}
		fmt.Printf("%s voted on comment. Upvotes: %d, Downvotes: %d\n", user.Username, comment.Upvotes, comment.Downvotes)
	}
}

// SendDirectMessage allows one user to send a message to another
func (e *Engine) SendDirectMessage(fromUsername, toUsername, content string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	fromUser, exists := e.users[fromUsername]
	if !exists {
		fmt.Println("Sender user not found!")
		return
	}

	toUser, exists := e.users[toUsername]
	if !exists {
		fmt.Println("Recipient user not found!")
		return
	}

	message := &DirectMessage{
		From:    fromUser,
		To:      toUser,
		Content: content,
	}
	e.messages[toUsername] = append(e.messages[toUsername], message)

	fmt.Printf("Message sent from %s to %s: %s\n", fromUsername, toUsername, content)
}

// GetDirectMessages retrieves all direct messages for a user
func (e *Engine) GetDirectMessages(username string) []*DirectMessage {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.messages[username]
}

// GetSubreddits returns all subreddits
func (e *Engine) GetSubreddits() map[string]*Subreddit {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.subreddits
}

// ReplyToDirectMessage allows replying to a direct message
func (e *Engine) ReplyToDirectMessage(messageID string, fromUsername string, content string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	fromUser, exists := e.users[fromUsername]
	if !exists {
		fmt.Println("Sender user not found!")
		return
	}

	// Find original message and recipient
	var originalMessage *DirectMessage
	var toUser *User
	for _, messages := range e.messages {
		for _, msg := range messages {
			if msg.ID == messageID {
				originalMessage = msg
				toUser = msg.From // Reply goes to original sender
				break
			}
		}
	}

	if originalMessage == nil || toUser == nil {
		fmt.Println("Original message not found!")
		return
	}

	reply := &DirectMessage{
		ID:        fmt.Sprintf("msg-%d", rand.Int()),
		From:      fromUser,
		To:        toUser,
		Content:   content,
		Timestamp: time.Now(),
		ReplyTo:   messageID,
	}

	e.messages[toUser.Username] = append(e.messages[toUser.Username], reply)
	fmt.Printf("Reply sent from %s to %s: %s\n", fromUsername, toUser.Username, content)
}

// GetUserFeed returns recent posts from subscribed subreddits
func (e *Engine) GetUserFeed(username string, limit int) []*Post {
	e.mu.Lock()
	defer e.mu.Unlock()

	var feed []*Post
	// Remove the unused variable declaration
	_, exists := e.users[username]
	if !exists {
		return feed
	}

	// Collect posts from subscribed subreddits
	for _, subreddit := range e.subreddits {
		if _, isMember := subreddit.Members[username]; isMember {
			feed = append(feed, subreddit.Posts...)
		}
	}

	// Sort by timestamp (newest first)
	sort.Slice(feed, func(i, j int) bool {
		return feed[i].Timestamp.After(feed[j].Timestamp)
	})

	// Apply limit if specified
	if limit > 0 && len(feed) > limit {
		return feed[:limit]
	}
	return feed
}

// Helper function to find a comment in a post's comment tree
func findComment(comments []*Comment, commentID string) *Comment {
	for _, comment := range comments {
		if comment.ID == commentID {
			return comment
		}
		// Search in replies
		if len(comment.Replies) > 0 {
			if found := findComment(comment.Replies, commentID); found != nil {
				return found
			}
		}
	}
	return nil
}

// GetPopularSubreddits returns subreddits sorted by member count
func (e *Engine) GetPopularSubreddits() []*Subreddit {
	e.mu.Lock()
	defer e.mu.Unlock()

	subreddits := make([]*Subreddit, 0, len(e.subreddits))
	for _, s := range e.subreddits {
		subreddits = append(subreddits, s)
	}

	sort.Slice(subreddits, func(i, j int) bool {
		return len(subreddits[i].Members) > len(subreddits[j].Members)
	})

	return subreddits
}
