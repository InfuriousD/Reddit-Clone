// cmd/test/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"reddit-clone/client/reddit"
	"sync"
	"time"
)

const numClients = 3 // Define this as a constant at package level

func main() {
	log.Println("Starting Reddit Clone Test")

	// Test with multiple concurrent clients
	var wg sync.WaitGroup

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			testClient(clientID)
		}(i)
		// Slight delay between client starts for better log readability
		time.Sleep(time.Second)
	}

	wg.Wait()
	log.Println("All tests completed")
}

func testClient(clientID int) {
	client := reddit.NewClient("http://localhost:8080")
	ctx := context.Background()

	username := fmt.Sprintf("testuser%d", clientID)
	password := fmt.Sprintf("testpass%d", clientID)

	// 1. Register
	log.Printf("Client %d: Registering user %s", clientID, username)
	if err := client.Register(username, password); err != nil {
		log.Printf("Client %d: Registration failed: %v", clientID, err)
		return
	}
	log.Printf("Client %d: Registration successful", clientID)

	// 2. Login
	log.Printf("Client %d: Logging in as %s", clientID, username)
	if err := client.Login(username, password); err != nil {
		log.Printf("Client %d: Login failed: %v", clientID, err)
		return
	}
	log.Printf("Client %d: Login successful", clientID)

	// 3. Create Subreddit
	subredditName := fmt.Sprintf("subreddit%d", clientID)
	log.Printf("Client %d: Creating subreddit %s", clientID, subredditName)
	if err := client.CreateSubreddit(ctx, subredditName, "Test subreddit"); err != nil {
		log.Printf("Client %d: Create subreddit failed: %v", clientID, err)
		return
	}
	log.Printf("Client %d: Subreddit created successfully", clientID)

	// 4. Create Post
	log.Printf("Client %d: Creating post in %s", clientID, subredditName)
	post, err := client.CreatePost(ctx,
		fmt.Sprintf("Test post from client %d", clientID),
		fmt.Sprintf("This is test content from client %d", clientID),
		subredditName)
	if err != nil {
		log.Printf("Client %d: Create post failed: %v", clientID, err)
		return
	}
	log.Printf("Client %d: Post created successfully with ID: %s", clientID, post.ID)

	// 5. Add Comment
	log.Printf("Client %d: Adding comment to post", clientID)
	comment, err := client.CreateComment(ctx, post.ID,
		fmt.Sprintf("Test comment from client %d", clientID), nil)
	if err != nil {
		log.Printf("Client %d: Create comment failed: %v", clientID, err)
		return
	}
	log.Printf("Client %d: Comment added successfully with ID: %s", clientID, comment.ID)

	// 6. Vote on post
	log.Printf("Client %d: Voting on post", clientID)
	if err := client.Vote(ctx, post.ID, true, "post"); err != nil {
		log.Printf("Client %d: Vote failed: %v", clientID, err)
		return
	}
	log.Printf("Client %d: Vote recorded successfully", clientID)

	// 7. Get Feed
	log.Printf("Client %d: Fetching feed", clientID)
	feed, err := client.GetFeed(ctx, 0, 10)
	if err != nil {
		log.Printf("Client %d: Get feed failed: %v", clientID, err)
		return
	}
	log.Printf("Client %d: Retrieved %d posts from feed", clientID, len(feed.Posts))

	// 8. Send Message to another user
	otherUser := fmt.Sprintf("testuser%d", (clientID+1)%numClients)
	log.Printf("Client %d: Sending message to %s", clientID, otherUser)
	if err := client.SendMessage(ctx, otherUser,
		fmt.Sprintf("Test message from client %d", clientID)); err != nil {
		log.Printf("Client %d: Send message failed: %v", clientID, err)
		return
	}
	log.Printf("Client %d: Message sent successfully", clientID)

	// 9. View Messages
	log.Printf("Client %d: Checking messages", clientID)
	messages, err := client.GetMessages(ctx)
	if err != nil {
		log.Printf("Client %d: Get messages failed: %v", clientID, err)
		return
	}
	log.Printf("Client %d: Retrieved %d messages", clientID, len(messages))

	log.Printf("Client %d: All tests completed successfully", clientID)
}
