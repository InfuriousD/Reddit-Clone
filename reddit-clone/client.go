package main

import (
	"fmt"
	"math/rand"
	"reddit-clone/engine"
	"time"
)

func simulateUserActivity(engine *engine.Engine, user *engine.User) {
	// Simulate random activity
	for i := 0; i < 10; i++ {
		// Randomly post in a random subreddit
		if rand.Intn(2) == 0 {
			subredditName := "golang" // Example subreddit
			content := fmt.Sprintf("Random post %d", i)
			engine.PostInSubreddit(subredditName, user.Username, content)
		}

		// Randomly comment on a post
		if rand.Intn(2) == 0 {
			postID := fmt.Sprintf("%d", rand.Intn(100)) // Example post ID
			content := fmt.Sprintf("Random comment %d", i)
			engine.CommentOnPost("golang", postID, user.Username, content)
		}

		// Randomly upvote or downvote
		if rand.Intn(2) == 0 {
			postID := fmt.Sprintf("%d", rand.Intn(100)) // Example post ID
			upvote := rand.Intn(2) == 0
			engine.Upvote(postID, upvote, user.Username)
		}

		// Simulate a delay for user activity
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func main() {
	// Initialize engine
	engine := engine.NewEngine()

	// Register some users
	user1 := engine.RegisterAccount("user1")
	user2 := engine.RegisterAccount("user2")

	// Create a subreddit
	engine.CreateSubreddit("golang")

	// Simulate user activity
	go simulateUserActivity(engine, user1)
	go simulateUserActivity(engine, user2)

	// Simulate users joining subreddit
	engine.JoinSubreddit(user1.Username, "golang")
	engine.JoinSubreddit(user2.Username, "golang")

	// Give the simulation some time to complete
	time.Sleep(5 * time.Second)
}
