package chat

import (
	"job_swipe/internal/database"
	"job_swipe/internal/models"
	"job_swipe/internal/response"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Hub    *Hub
	UserID uint
	Conn   *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	Clients    map[uint]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uint]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("User %d connected", client.UserID)
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
				log.Printf("User %d disconnected", client.UserID)
			}
			h.mu.Unlock()
		}
	}
}

var GlobalHub = NewHub()

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Printf("Received message from %d: %s", c.UserID, message)
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{Hub: hub, UserID: userID.(uint), Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

type MessageInput struct {
	ReceiverID uint   `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

func SendMessage(c *gin.Context) {
	senderID, _ := c.Get("user_id")

	var input MessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	var receiver models.User
	if err := database.DB.First(&receiver, input.ReceiverID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Receiver not found", nil)
		return
	}

	msg := models.Message{
		SenderID:   senderID.(uint),
		ReceiverID: input.ReceiverID,
		Content:    input.Content,
		IsRead:     false,
	}

	if err := database.DB.Create(&msg).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to send message", err.Error())
		return
	}

	GlobalHub.mu.Lock()
	if client, ok := GlobalHub.Clients[input.ReceiverID]; ok {
		client.Send <- []byte(input.Content)
	}
	GlobalHub.mu.Unlock()

	response.Success(c, "Message sent", msg)
}

func GetChatHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	otherUserIDStr := c.Param("user_id")

	otherUserID, err := strconv.Atoi(otherUserIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	var messages []models.Message
	if err := database.DB.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userID, otherUserID, otherUserID, userID,
	).Order("created_at asc").Find(&messages).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch messages", err.Error())
		return
	}

	response.Success(c, "Chat history", messages)
}

func GetConversations(c *gin.Context) {
	// userID, _ := c.Get("user_id")
	type Conversation struct {
		UserID      uint      `json:"user_id"`
		LastMessage string    `json:"last_message"`
		Timestamp   time.Time `json:"timestamp"`
	}

	response.Success(c, "Conversations list (Not implemented fully yet)", nil)
}
