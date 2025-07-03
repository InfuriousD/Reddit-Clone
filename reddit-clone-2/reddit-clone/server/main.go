// server/main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reddit-clone/models"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Server struct {
	router   *mux.Router
	posts    map[uuid.UUID]*models.Post
	comments map[uuid.UUID]*models.Comment
	messages map[uuid.UUID]*models.DirectMessage
	users    map[string]*models.User
	hub      *Hub
}

func NewServer() *Server {
	hub := newHub()
	go hub.run()

	s := &Server{
		router:   mux.NewRouter(),
		posts:    make(map[uuid.UUID]*models.Post),
		comments: make(map[uuid.UUID]*models.Comment),
		messages: make(map[uuid.UUID]*models.DirectMessage),
		users:    make(map[string]*models.User),
		hub:      hub,
	}
	s.routes()
	return s
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // For development
	},
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func initDB() (*sql.DB, error) {
	// For development, using default PostgreSQL settings
	connStr := "postgres://postgres:postgres@localhost:5432/reddit_clone?sslmode=disable"
	return sql.Open("postgres", connStr)
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (s *Server) routes() {
	// Auth routes
	s.router.HandleFunc("/api/register", s.handleRegister()).Methods("POST")
	s.router.HandleFunc("/api/login", s.handleLogin()).Methods("POST")

	// Subreddit routes
	s.router.HandleFunc("/api/subreddits", s.handleCreateSubreddit()).Methods("POST")
	s.router.HandleFunc("/api/subreddits/{name}/join", s.handleJoinSubreddit()).Methods("POST")

	// Post routes
	s.router.HandleFunc("/api/posts", s.handleCreatePost()).Methods("POST")
	s.router.HandleFunc("/api/posts/{id}", s.handleGetPost()).Methods("GET")
	s.router.HandleFunc("/api/posts/{id}/vote", s.handleVotePost()).Methods("POST")

	// Comment routes
	s.router.HandleFunc("/api/posts/{id}/comments", s.handleCreateComment()).Methods("POST")
	s.router.HandleFunc("/api/comments/{id}", s.handleGetComment()).Methods("GET")
	s.router.HandleFunc("/api/comments/{id}/vote", s.handleVoteComment()).Methods("POST")

	// WebSocket
	s.router.HandleFunc("/ws", s.handleWebSocket())

	s.router.HandleFunc("/api/messages", s.handleSendMessage()).Methods("POST")
	s.router.HandleFunc("/api/messages", s.handleGetMessages()).Methods("GET")

	s.router.HandleFunc("/api/messages", s.handleGetMessages()).Methods("GET")
	s.router.HandleFunc("/api/messages", s.handleSendMessage()).Methods("POST")

	s.router.HandleFunc("/api/search", s.handleSearch()).Methods("GET")

}

func (s *Server) handleSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get search query from URL parameters
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "Search query is required", http.StatusBadRequest)
			return
		}

		query = strings.ToLower(query)
		var results []*models.Post

		// Search through all posts
		for _, post := range s.posts {
			// Search in title and content
			if strings.Contains(strings.ToLower(post.Title), query) ||
				strings.Contains(strings.ToLower(post.Content), query) ||
				strings.Contains(strings.ToLower(post.SubredditName), query) {
				results = append(results, post)
			}
		}

		// Sort results by creation time (newest first)
		sort.Slice(results, func(i, j int) bool {
			return results[i].CreatedAt.After(results[j].CreatedAt)
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

func (s *Server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Check if user exists and password matches
		user, exists := s.users[req.Username]
		if !exists {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Verify password
		if user.PasswordHash != req.Password {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Generate token for successful login
		token := fmt.Sprintf("token-%s-%d", req.Username, time.Now().Unix())

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"token":    token,
			"username": req.Username,
		})
	}
}

func (s *Server) handleCreateSubreddit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// For development, return success
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Subreddit created successfully",
		})
	}
}

func (s *Server) handleJoinSubreddit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		// For development, return success
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Joined subreddit: " + name,
		})
	}
}

func (s *Server) handleCreateComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		postIDStr := vars["id"]
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		var req struct {
			Content  string     `json:"content"`
			ParentID *uuid.UUID `json:"parent_id,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create new comment
		comment := &models.Comment{
			ID:         uuid.New(),
			Content:    req.Content,
			PostID:     postID,
			ParentID:   req.ParentID,
			AuthorName: r.Header.Get("X-User"), // In production, get from auth token
		}

		// Store comment
		s.comments[comment.ID] = comment

		// Update post's comments if it exists
		if post, exists := s.posts[postID]; exists {
			post.CommentsCount++
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(comment)
	}
}

func (s *Server) handleCreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Title     string `json:"title"`
			Content   string `json:"content"`
			Subreddit string `json:"subreddit"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		post := &models.Post{
			ID:            uuid.New(),
			Title:         req.Title,
			Content:       req.Content,
			SubredditName: req.Subreddit,
			AuthorName:    r.Header.Get("X-User"), // In production, get from auth token
		}

		s.posts[post.ID] = post

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(post)
	}
}

func (s *Server) handleGetPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		postID, err := uuid.Parse(vars["id"])
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		post, exists := s.posts[postID]
		if !exists {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(post)
	}
}

func (s *Server) handleVotePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		postID, err := uuid.Parse(vars["id"])
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		var req struct {
			IsUpvote bool `json:"is_upvote"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get user from auth header (for demo purposes)
		username := r.Header.Get("X-User")
		if username == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Find the post
		post, exists := s.posts[postID]
		if !exists {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		// Record the vote
		if req.IsUpvote {
			post.Upvotes++
			post.Score++
		} else {
			post.Downvotes++
			post.Score--
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(post)
	}
}

func (s *Server) handleVoteComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		commentID, err := uuid.Parse(vars["id"])
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		var req struct {
			IsUpvote bool `json:"is_upvote"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get user from auth header
		username := r.Header.Get("X-User")
		if username == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Find the comment
		comment, exists := s.comments[commentID]
		if !exists {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}

		// Record the vote
		if req.IsUpvote {
			comment.Upvotes++
			comment.Score++
		} else {
			comment.Downvotes++
			comment.Score--
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(comment)
	}
}
func (s *Server) handleGetComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		commentID, err := uuid.Parse(vars["id"])
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		comment, exists := s.comments[commentID]
		if !exists {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(comment)
	}
}

func (s *Server) handleWebSocket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			return
		}

		client := &Client{
			hub:  s.hub,
			conn: conn,
			send: make(chan []byte, 256),
		}
		client.hub.register <- client

		go client.writePump()
		go client.readPump()
	}
}

func (s *Server) handleSendMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ToUser  string `json:"to_user"`
			Content string `json:"content"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get sender from auth token (for demo, we'll get it from header)
		fromUser := r.Header.Get("X-User")
		if fromUser == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		message := &models.DirectMessage{
			ID:        uuid.New(),
			FromUser:  fromUser,
			ToUser:    req.ToUser,
			Content:   req.Content,
			CreatedAt: time.Now(),
		}

		// Store the message
		s.messages[message.ID] = message

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(message)
	}
}

func (s *Server) handleGetMessages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user from auth token (for demo, we'll get it from header)
		username := r.Header.Get("X-User")
		if username == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Collect all messages for this user
		var userMessages []*models.DirectMessage
		for _, msg := range s.messages {
			if msg.ToUser == username || msg.FromUser == username {
				userMessages = append(userMessages, msg)
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(userMessages)
	}
}

func (s *Server) handleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Check if user already exists
		if _, exists := s.users[req.Username]; exists {
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}

		// Create new user
		s.users[req.Username] = &models.User{
			Username:     req.Username,
			PasswordHash: req.Password, // In production, hash the password
			CreatedAt:    time.Now(),
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User registered successfully",
		})
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			w.Close()
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		c.hub.broadcast <- message
	}
}
func main() {
	server := NewServer()

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", server.router); err != nil {
		log.Fatal(err)
	}
}
