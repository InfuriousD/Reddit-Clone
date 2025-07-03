// client/reddit/client.go
package reddit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reddit-clone/models"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	ws         *websocket.Conn
	authToken  string
	username   string
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// Auth methods
func (c *Client) Register(username, password string) error {
	payload := map[string]string{
		"username": username,
		"password": password,
	}
	return c.post("/api/register", payload, nil)
}

func (c *Client) Login(username, password string) error {
	payload := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal login payload: %v", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create login request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send login request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("Invalid username or password")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned error: %d", resp.StatusCode)
	}

	var response struct {
		Token    string `json:"token"`
		Username string `json:"username"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode login response: %v", err)
	}

	c.authToken = response.Token
	c.username = username // Store username for future requests

	return nil
}

// Subreddit methods
func (c *Client) CreateSubreddit(ctx context.Context, name, description string) error {
	payload := map[string]string{
		"name":        name,
		"description": description,
	}
	return c.post("/api/subreddits", payload, nil)
}

func (c *Client) JoinSubreddit(ctx context.Context, name string) error {
	return c.post(fmt.Sprintf("/api/subreddits/%s/join", name), nil, nil)
}

func (c *Client) LeaveSubreddit(ctx context.Context, name string) error {
	return c.post(fmt.Sprintf("/api/subreddits/%s/leave", name), nil, nil)
}

// Post methods
func (c *Client) CreatePost(ctx context.Context, title, content, subreddit string) (*models.Post, error) {
	payload := map[string]string{
		"title":     title,
		"content":   content,
		"subreddit": subreddit,
	}
	var post models.Post
	err := c.post("/api/posts", payload, &post)
	return &post, err
}

func (c *Client) GetPost(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	var post models.Post
	err := c.get(fmt.Sprintf("/api/posts/%s", postID), &post)
	return &post, err
}

// Comment methods
func (c *Client) CreateComment(ctx context.Context, postID uuid.UUID, content string, parentID *uuid.UUID) (*models.Comment, error) {
	payload := map[string]interface{}{
		"content":   content,
		"post_id":   postID,
		"parent_id": parentID,
	}
	var comment models.Comment
	err := c.post(fmt.Sprintf("/api/posts/%s/comments", postID), payload, &comment)
	return &comment, err
}

// Vote methods
func (c *Client) Vote(ctx context.Context, targetID uuid.UUID, isUpvote bool, targetType string) error {
	payload := map[string]interface{}{
		"is_upvote": isUpvote,
	}

	var endpoint string
	if targetType == "post" {
		endpoint = fmt.Sprintf("/api/posts/%s/vote", targetID)
	} else if targetType == "comment" {
		endpoint = fmt.Sprintf("/api/comments/%s/vote", targetID)
	} else {
		return fmt.Errorf("invalid target type: %s", targetType)
	}

	return c.post(endpoint, payload, nil)
}

// Feed methods
func (c *Client) GetFeed(ctx context.Context, offset, limit int) (*models.FeedResponse, error) {
	var feed models.FeedResponse
	err := c.get(fmt.Sprintf("/api/feed?offset=%d&limit=%d", offset, limit), &feed)
	return &feed, err
}

// Message methods
func (c *Client) SendMessage(ctx context.Context, toUser, content string) error {
	payload := map[string]string{
		"to_user": toUser,
		"content": content,
	}
	return c.post("/api/messages", payload, nil)
}

func (c *Client) GetMessages(ctx context.Context) ([]*models.DirectMessage, error) {
	var messages []*models.DirectMessage
	err := c.get("/api/messages", &messages)
	return messages, err
}

// Search method
func (c *Client) Search(ctx context.Context, query string) ([]*models.Post, error) {
	var posts []*models.Post
	err := c.get(fmt.Sprintf("/api/search?q=%s", url.QueryEscape(query)), &posts)
	return posts, err
}

// WebSocket methods
func (c *Client) ConnectWebSocket(username string) error {
	wsURL := fmt.Sprintf("ws://%s/ws?username=%s", c.baseURL, username)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}
	c.ws = conn
	return nil
}

// Helper methods
func (c *Client) post(endpoint string, payload interface{}, response interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
		req.Header.Set("X-User", c.username) // Add this line
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("server returned error: %d", resp.StatusCode)
	}

	if response != nil {
		return json.NewDecoder(resp.Body).Decode(response)
	}
	return nil
}

func (c *Client) get(endpoint string, response interface{}) error {
	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return err
	}

	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
		req.Header.Set("X-User", c.username) // Add this line
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("server returned error: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(response)
}
