package websocket

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/leehai1107/cmm_server/pkg/logger"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client struct for websocket connection and message sending
type Client struct {
	ID   string
	Conn *websocket.Conn
	send chan Message
	hub  *Hub
}

// NewClient creates a new client
func NewClient(id string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{ID: id, Conn: conn, send: make(chan Message, 256), hub: hub}
}

// Client goroutine to read messages from client
func (c *Client) Read() {

	defer func() {
		c.hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			logger.Errorf("Error: %v", err)
			break
		}
		c.hub.broadcast <- msg
	}
}

// Client goroutine to write messages to client
func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			} else {
				err := c.Conn.WriteJSON(message)
				if err != nil {
					logger.Errorf("Error: %v", err)
					break
				}
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}

// Client closing channel to unregister client
func (c *Client) Close() {
	close(c.send)
}

// Function to handle websocket connection and register client to hub and start goroutines
func serveWS(ctx *gin.Context, roomId string, hub *Hub) {
	// Validate hub
	if hub == nil {
		logger.Errorf("Hub is nil for RoomId: %s", roomId)
		return
	}

	// Validate Gin context
	if ctx == nil || ctx.Writer == nil || ctx.Request == nil {
		logger.Errorf("Invalid Gin context for RoomId: %s", roomId)
		return
	}

	logger.Infof("New connection: %s, RoomId: %s", ctx.Request.RemoteAddr, roomId)

	// WebSocket upgrade
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Errorf("WebSocket upgrade failed: %v", err)
		return
	}
	if ws == nil {
		logger.Errorf("WebSocket connection is nil after upgrade, RoomId: %s", roomId)
		return
	}

	// Create and register client
	client := NewClient(roomId, ws, hub)
	if client == nil {
		logger.Errorf("Failed to create a new client for RoomId: %s", roomId)
		ws.Close()
		return
	}

	// Register the client to the hub
	hub.register <- client

	// Start client goroutines
	go client.Write()
	go client.Read()
}
