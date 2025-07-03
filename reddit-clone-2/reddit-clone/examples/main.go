// examples/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"reddit-clone/client/reddit"
)

func main() {
	// Create a new client
	client := reddit.NewClient("http://localhost:8080")

	// Test all functionality
	if err := testAllFunctionality(client); err != nil {
		log.Fatalf("Test failed: %v", err)
	}
}

func testAllFunctionality(client *reddit.Client) error {
	ctx := context.Background()

	// 1. Register and login users
	users := []struct {
		username string
		password string
	}{
		{"testuser1", "password1"},
		{"testuser2", "password2"},
	}

	for _, user := range users {
		if err := client.Register(user.username, user.password); err != nil {
			return fmt.Errorf("failed to register user %s: %v", user.username, err)
		}
		fmt.Printf("✓ Registered user: %s\n", user.username)
	}

	// Login as first user
	if err := client.Login(users[0].username, users[0].password); err != nil {
		return fmt.Errorf("failed to login: %v", err)
	}
	fmt.Println("✓ Logged in successfully")

	// 2. Create and join subreddit
	subredditName := "testsubreddit"
	if err := client.CreateSubreddit(ctx, subredditName, "A test subreddit"); err != nil {
		return fmt.Errorf("failed to create subreddit: %v", err)
	}
	fmt.Printf("✓ Created subreddit: %s\n", subredditName)

	if err := client.JoinSubreddit(ctx, subredditName); err != nil {
		return fmt.Errorf("failed to join subreddit: %v", err)
	}
	fmt.Printf("✓ Joined subreddit: %s\n", subredditName)

	// 3. Create a post
	post, err := client.CreatePost(ctx, "Test Post", "This is a test post content", subredditName)
	if err != nil {
		return fmt.Errorf("failed to create post: %v", err)
	}
	fmt.Printf("✓ Created post: %s\n", post.Title)

	// 4. Create comments
	comment1, err := client.CreateComment(ctx, post.ID, "This is a top-level comment", nil)
	if err != nil {
		return fmt.Errorf("failed to create comment: %v", err)
	}
	fmt.Println("✓ Created top-level comment")

	// Create a reply to the first comment
	_, err = client.CreateComment(ctx, post.ID, "This is a reply to the first comment", &comment1.ID)
	if err != nil {
		return fmt.Errorf("failed to create reply: %v", err)
	}
	fmt.Println("✓ Created reply comment")

	// 5. Vote on post and comments
	if err := client.Vote(ctx, post.ID, true, "post"); err != nil {
		return fmt.Errorf("failed to upvote post: %v", err)
	}
	fmt.Println("✓ Upvoted post")

	if err := client.Vote(ctx, comment1.ID, true, "comment"); err != nil {
		return fmt.Errorf("failed to upvote comment: %v", err)
	}
	fmt.Println("✓ Upvoted comment")

	// 6. Send direct message
	if err := client.SendMessage(ctx, users[1].username, "Hello! This is a test message"); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	fmt.Println("✓ Sent direct message")

	// 7. Get feed
	feed, err := client.GetFeed(ctx, 0, 10)
	if err != nil {
		return fmt.Errorf("failed to get feed: %v", err)
	}
	fmt.Printf("✓ Got feed with %d posts\n", len(feed.Posts))

	fmt.Println("\nAll functionality tested successfully!")
	return nil
}
