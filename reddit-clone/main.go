package main

import (
	"fmt"
	"math"
	"math/rand"
	"reddit-clone/engine"
	"sync"
	"time"
)

// Generate a Zipf distribution for subreddit memberships
func simulateZipfDistribution(engine *engine.Engine, subreddits []string, users []*engine.User) {
	ratio := 1.07
	for i, subreddit := range subreddits {
		memberCount := int(float64(len(users)) / math.Pow(float64(i+1), ratio))
		for j := 0; j < memberCount; j++ {
			user := users[rand.Intn(len(users))]
			engine.JoinSubreddit(user.Username, subreddit)
		}
	}
}

func simulateUserActivity(engine *engine.Engine, user *engine.User, subreddits []string, wg *sync.WaitGroup) {
	defer wg.Done()

	connected := true
	engine.SetUserConnection(user.Username, connected)

	// Store post IDs instead of objects
	var postIDs []string
	var commentIDs []string

	// Activity simulation
	for i := 0; i < 15; i++ {
		// Randomly toggle connection status
		if rand.Float32() < 0.1 {
			connected = !connected
			engine.SetUserConnection(user.Username, connected)
		}

		if connected {
			// Create posts
			if rand.Float32() < 0.3 {
				subreddit := subreddits[rand.Intn(len(subreddits))]
				content := fmt.Sprintf("Post %d by %s", i, user.Username)
				if post := engine.PostInSubreddit(subreddit, user.Username, content); post != nil {
					postIDs = append(postIDs, post.ID)
				}
			}

			// Create comments
			if len(postIDs) > 0 && rand.Float32() < 0.4 {
				postID := postIDs[rand.Intn(len(postIDs))]
				subreddit := subreddits[rand.Intn(len(subreddits))]
				content := fmt.Sprintf("Comment on post %s", postID)
				if comment := engine.CommentOnPost(subreddit, postID, user.Username, content); comment != nil {
					commentIDs = append(commentIDs, comment.ID)
				}
			}

			// Create reposts
			if len(postIDs) > 0 && rand.Float32() < 0.1 {
				postID := postIDs[rand.Intn(len(postIDs))]
				targetSubreddit := subreddits[rand.Intn(len(subreddits))]
				engine.Repost(postID, user.Username, targetSubreddit)
			}

			// Send messages
			if rand.Float32() < 0.2 {
				targetUser := fmt.Sprintf("user%d", rand.Intn(1000)+1)
				content := fmt.Sprintf("Message from %s", user.Username)
				engine.SendDirectMessage(user.Username, targetUser, content)
			}

			// Voting activity
			if len(postIDs) > 0 {
				for j := 0; j < 5; j++ {
					postID := postIDs[rand.Intn(len(postIDs))]
					engine.Upvote(postID, rand.Float32() > 0.3, user.Username)
				}
			}

			// Get feed and interact
			feed := engine.GetUserFeed(user.Username, 10)
			for _, post := range feed {
				if rand.Float32() < 0.3 {
					engine.Upvote(post.ID, rand.Float32() > 0.3, user.Username)
				}
			}
		}

		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	}
}

func displayDetailedMetrics(engine *engine.Engine, users []*engine.User, numUsers int, elapsed time.Duration) {
	totalPosts := 0
	totalComments := 0
	totalVotes := 0
	activeUsers := 0

	for _, user := range users {
		if user.Connected {
			activeUsers++
		}
		totalPosts += len(user.Posts)
		totalComments += len(user.Comments)
		totalVotes += user.Karma
	}

	fmt.Printf("\nPerformance Metrics:\n")
	fmt.Printf("Runtime: %s\n", elapsed)
	fmt.Printf("Total Users: %d (Active: %d)\n", numUsers, activeUsers)
	fmt.Printf("Total Posts: %d\n", totalPosts)
	fmt.Printf("Total Comments: %d\n", totalComments)
	fmt.Printf("Total Votes: %d\n", totalVotes)
	fmt.Printf("Operations/sec: %.2f\n", float64(totalPosts+totalComments+totalVotes)/elapsed.Seconds())
}

func main() {
	redditEngine := engine.NewEngine()

	numUsers := 1000
	var users []*engine.User
	for i := 0; i < numUsers; i++ {
		user := redditEngine.RegisterAccount(fmt.Sprintf("user%d", i+1))
		users = append(users, user)
	}

	subreddits := []string{"golang", "python", "java", "csharp", "rust"}
	for _, subreddit := range subreddits {
		redditEngine.CreateSubreddit(subreddit)
	}

	simulateZipfDistribution(redditEngine, subreddits, users)

	var wg sync.WaitGroup
	wg.Add(len(users))
	start := time.Now()

	batchSize := 500
	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}

		for j := i; j < end; j++ {
			go simulateUserActivity(redditEngine, users[j], subreddits, &wg)
		}
	}

	wg.Wait()
	elapsed := time.Since(start)

	displayDetailedMetrics(redditEngine, users, numUsers, elapsed)
}
