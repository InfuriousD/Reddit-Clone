// cmd/demo/main.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"reddit-clone/client/reddit"
	"strings"

	"github.com/google/uuid"
)

type Demo struct {
	client   *reddit.Client
	reader   *bufio.Reader
	loggedIn bool
	username string
}

func NewDemo() *Demo {
	return &Demo{
		client:   reddit.NewClient("http://localhost:8080"),
		reader:   bufio.NewReader(os.Stdin),
		loggedIn: false,
	}
}

func (d *Demo) Start() {
	fmt.Println("=== Reddit Clone Demo ===")
	for {
		if !d.loggedIn {
			d.showAuthMenu()
		} else {
			d.showMainMenu()
		}
	}
}

func (d *Demo) showAuthMenu() {
	fmt.Println("\nAuthentication Menu:")
	fmt.Println("1. Create Account")
	fmt.Println("2. Login")
	fmt.Println("3. Exit")
	fmt.Print("Choose an option: ")

	choice := d.readLine()
	switch choice {
	case "1":
		d.handleRegister()
	case "2":
		d.handleLogin()
	case "3":
		os.Exit(0)
	default:
		fmt.Println("Invalid option")
	}
}

func (d *Demo) showMainMenu() {
	fmt.Printf("\nLogged in as: %s\n", d.username)
	fmt.Println("\nMain Menu:")
	fmt.Println("1. Create Subreddit")
	fmt.Println("2. Join Subreddit")
	fmt.Println("3. Create Post")
	fmt.Println("4. Search Posts")
	fmt.Println("5. View Feed")
	fmt.Println("6. Comment on Post")
	fmt.Println("7. Vote on Post")
	fmt.Println("8. Send Message")
	fmt.Println("9. View Messages")
	fmt.Println("10. Logout")
	fmt.Print("Choose an option: ")

	choice := d.readLine()
	switch choice {
	case "1":
		d.handleCreateSubreddit()
	case "2":
		d.handleJoinSubreddit()
	case "3":
		d.handleCreatePost()
	case "4":
		d.handleSearch()
	case "5":
		d.handleViewFeed()
	case "6":
		d.handleComment()
	case "7":
		d.handleVote()
	case "8":
		d.handleSendMessage()
	case "9":
		d.handleViewMessages()
	case "10":
		d.handleLogout()
	default:
		fmt.Println("Invalid option")
	}
}

func (d *Demo) handleRegister() {
	fmt.Print("Enter desired username: ")
	username := d.readLine()
	fmt.Print("Enter password: ")
	password := d.readLine()

	log.Printf("Sending POST request to /api/register")
	err := d.client.Register(username, password)
	if err != nil {
		fmt.Printf("Registration failed: %v\n", err)
		return
	}
	fmt.Println("Registration successful!")
}

func (d *Demo) handleLogin() {
	fmt.Print("Enter username: ")
	username := d.readLine()
	fmt.Print("Enter password: ")
	password := d.readLine()

	log.Printf("Sending POST request to /api/login")
	err := d.client.Login(username, password)
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		return
	}
	d.loggedIn = true
	d.username = username
	fmt.Println("Login successful!")
}

func (d *Demo) handleCreateSubreddit() {
	fmt.Print("Enter subreddit name: ")
	name := d.readLine()
	fmt.Print("Enter description: ")
	description := d.readLine()

	log.Printf("Sending POST request to /api/subreddits")
	ctx := context.Background()
	err := d.client.CreateSubreddit(ctx, name, description)
	if err != nil {
		fmt.Printf("Failed to create subreddit: %v\n", err)
		return
	}
	fmt.Println("Subreddit created successfully!")
}

func (d *Demo) handleJoinSubreddit() {
	fmt.Print("Enter subreddit name to join: ")
	name := d.readLine()

	ctx := context.Background()
	err := d.client.JoinSubreddit(ctx, name)
	if err != nil {
		fmt.Printf("Failed to join subreddit: %v\n", err)
		return
	}
	fmt.Printf("Successfully joined r/%s!\n", name)
}

func (d *Demo) handleCreatePost() {
	fmt.Print("Enter subreddit name: ")
	subreddit := d.readLine()
	fmt.Print("Enter post title: ")
	title := d.readLine()
	fmt.Print("Enter post content: ")
	content := d.readLine()

	ctx := context.Background()
	post, err := d.client.CreatePost(ctx, title, content, subreddit)
	if err != nil {
		fmt.Printf("Failed to create post: %v\n", err)
		return
	}
	fmt.Printf("Post created successfully! Post ID: %s\n", post.ID)
}

func (d *Demo) handleSearch() {
	fmt.Print("Enter search query: ")
	query := d.readLine()

	ctx := context.Background()
	posts, err := d.client.Search(ctx, query)
	if err != nil {
		fmt.Printf("Search failed: %v\n", err)
		return
	}

	fmt.Printf("\nFound %d posts:\n", len(posts))
	for i, post := range posts {
		fmt.Printf("\n%d. %s\n", i+1, post.Title)
		fmt.Printf("   Posted by u/%s in r/%s\n", post.AuthorName, post.SubredditName)
		fmt.Printf("   Score: %d\n", post.Score)
	}
}

func (d *Demo) handleViewFeed() {
	ctx := context.Background()
	feed, err := d.client.GetFeed(ctx, 0, 10)
	if err != nil {
		fmt.Printf("Failed to get feed: %v\n", err)
		return
	}

	fmt.Println("\n=== Your Feed ===")
	for i, post := range feed.Posts {
		fmt.Printf("\n%d. %s\n", i+1, post.Title)
		fmt.Printf("   Posted by u/%s in r/%s\n", post.AuthorName, post.SubredditName)
		fmt.Printf("   Score: %d Comments: %d\n", post.Score, post.CommentsCount)
	}
}

func (d *Demo) handleComment() {
	fmt.Print("Enter post ID: ")
	postIDStr := d.readLine()
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		fmt.Printf("Invalid post ID: %v\n", err)
		return
	}

	fmt.Print("Enter your comment: ")
	content := d.readLine()

	ctx := context.Background()
	comment, err := d.client.CreateComment(ctx, postID, content, nil)
	if err != nil {
		fmt.Printf("Failed to create comment: %v\n", err)
		return
	}
	fmt.Printf("Comment posted successfully! Comment ID: %s\n", comment.ID)
}

func (d *Demo) handleVote() {
	fmt.Print("Enter post ID: ")
	postIDStr := d.readLine()
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		fmt.Printf("Invalid post ID: %v\n", err)
		return
	}

	fmt.Print("Upvote? (y/n): ")
	voteType := d.readLine()
	isUpvote := strings.ToLower(voteType) == "y"

	ctx := context.Background()
	err = d.client.Vote(ctx, postID, isUpvote, "post")
	if err != nil {
		fmt.Printf("Failed to vote: %v\n", err)
		return
	}
	fmt.Println("Vote recorded successfully!")
}

func (d *Demo) handleSendMessage() {
	fmt.Print("Enter recipient username: ")
	recipient := d.readLine()
	fmt.Print("Enter message content: ")
	content := d.readLine()

	ctx := context.Background()
	err := d.client.SendMessage(ctx, recipient, content)
	if err != nil {
		fmt.Printf("Failed to send message: %v\n", err)
		return
	}
	fmt.Println("Message sent successfully!")
}

func (d *Demo) handleViewMessages() {
	ctx := context.Background()
	messages, err := d.client.GetMessages(ctx)
	if err != nil {
		fmt.Printf("Failed to get messages: %v\n", err)
		return
	}

	fmt.Println("\n=== Your Messages ===")
	if len(messages) == 0 {
		fmt.Println("No messages found.")
		return
	}

	for i, msg := range messages {
		fmt.Printf("\n%d. From: u/%s\n", i+1, msg.FromUser)
		fmt.Printf("   Content: %s\n", msg.Content)
		fmt.Printf("   Sent at: %s\n", msg.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("   ---")
	}
}

func (d *Demo) handleLogout() {
	d.loggedIn = false
	d.username = ""
	fmt.Println("Logged out successfully!")
}

func (d *Demo) readLine() string {
	text, _ := d.reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func main() {
	demo := NewDemo()
	demo.Start()
}
