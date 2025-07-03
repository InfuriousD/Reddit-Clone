// server/websocket/hub.go
package websocket

import (
    "encoding/json"
    "github.com/gorilla/websocket"
    "log"
    "net/http"
    "sync"
)

type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte
    username string
}

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    userMap    map[string]*Client
    mu         sync.RWMutex
}

func newHub() *Hub {
    return &Hub{
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        clients:    make(map[*Client]bool),
        userMap:    make(map[string]*Client),
    }
}

func (h *Hub) run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.userMap[client.username] = client
            h.mu.Unlock()
            
        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                delete(h.userMap, client.username)
                close(client.send)
            }
            h.mu.Unlock()
            
        case message := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                    delete(h.userMap, client.username)
                }
            }
            h.mu.RUnlock()
        }
    }
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // In production, implement proper origin checking
    },
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    
    username := r.URL.Query().Get("username")
    client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), username: username}
    client.hub.register <- client
    
    go client.writePump()
    go client.readPump()
}
