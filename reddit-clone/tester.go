package main

import (
	"fmt"
	"reddit-clone/engine"
	"time"
)

func testRedditFunctionality() {
	fmt.Println("=== Starting Reddit Functionality Tests ===")
	e := engine.NewEngine()

	// Test 1: Account Registration
	fmt.Println("\n1. Testing Account Registration:")
	user1 := e.RegisterAccount("testUser1")
	user2 := e.RegisterAccount("testUser2")
	fmt.Printf("- Created users: %s, %s\n", user1.Username, user2.Username)

	// Test duplicate registration
	dupUser := e.RegisterAccount("testUser1")
	if dupUser == nil {
		fmt.Println("✓ Duplicate user registration prevented")
	}

	// Test 2: Subreddit Creation and Membership
	fmt.Println("\n2. Testing Subreddit Operations:")
	e.CreateSubreddit("testSubreddit")
	e.JoinSubreddit("testUser1", "testSubreddit")
	fmt.Println("- Created and joined subreddit")

	// Test 3: Posting
	fmt.Println("\n3. Testing Posting:")
	post := e.PostInSubreddit("testSubreddit", "testUser1", "Test post content")
	if post != nil {
		fmt.Printf("✓ Post created with ID: %s\n", post.ID)
	}

	// Test 4: Commenting
	fmt.Println("\n4. Testing Comments:")
	comment1 := e.CommentOnPost("testSubreddit", post.ID, "testUser2", "Test comment")
	if comment1 != nil {
		fmt.Println("✓ Comment created")
		// Test reply to comment
		reply := e.ReplyToComment(post.ID, comment1.ID, "testUser1", "Reply to comment")
		if reply != nil {
			fmt.Println("✓ Reply to comment created")
		}
	}

	// Test 5: Voting System
	fmt.Println("\n5. Testing Voting System:")
	initialKarma := user1.Karma
	e.Upvote(post.ID, true, "testUser2")
	if user1.Karma > initialKarma {
		fmt.Println("✓ Karma system working")
	}

	// Test 6: User Feed
	fmt.Println("\n6. Testing User Feed:")
	feed := e.GetUserFeed("testUser1", 10)
	fmt.Printf("- Feed contains %d posts\n", len(feed))

	// Test 7: Direct Messages
	fmt.Println("\n7. Testing Direct Messages:")
	e.SendDirectMessage("testUser1", "testUser2", "Test message")
	messages := e.GetDirectMessages("testUser2")
	if len(messages) > 0 {
		fmt.Printf("✓ Message delivered from %s to %s\n", messages[0].From.Username, messages[0].To.Username)
	}

	// Test 8: Repost Functionality
	fmt.Println("\n8. Testing Repost:")
	repost := e.Repost(post.ID, "testUser2", "testSubreddit")
	if repost != nil {
		fmt.Println("✓ Repost created")
	}

	// Test 9: Connection Status
	fmt.Println("\n9. Testing Connection Status:")
	e.SetUserConnection("testUser1", false)
	time.Sleep(time.Millisecond * 100)
	e.SetUserConnection("testUser1", true)
	fmt.Println("✓ Connection status changes simulated")

	fmt.Println("\n=== Functionality Test Summary ===")
	fmt.Printf("Users: %d\n", len(e.GetSubreddits()))
	fmt.Printf("Posts: %d\n", len(post.Subreddit.Posts))
	fmt.Printf("User1 Karma: %d\n", user1.Karma)
}

func main() {
	testRedditFunctionality()
}
