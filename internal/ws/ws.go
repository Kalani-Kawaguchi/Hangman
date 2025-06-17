package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Hub struct {
	clients map[string]map[*websocket.Conn]bool
	lock    sync.Mutex
}

var wsHub = &Hub{
	clients: make(map[string]map[*websocket.Conn]bool),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	lobbyID := r.URL.Query().Get("lobby")
	if lobbyID == "" {
		http.Error(w, "Missing lobby ID", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	wsHub.lock.Lock()
	if wsHub.clients[lobbyID] == nil {
		wsHub.clients[lobbyID] = make(map[*websocket.Conn]bool)
	}
	wsHub.clients[lobbyID][conn] = true
	wsHub.lock.Unlock()

	defer func() {
		wsHub.lock.Lock()
		delete(wsHub.clients[lobbyID], conn)
		wsHub.lock.Unlock()
		conn.Close()
	}()

	for {
		var msg WSMessage
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		log.Printf("Received from %s: %v\n", lobbyID, msg)

		// Handle specific message types later
	}
}

func BroadcastToLobby(lobbyID string, msg WSMessage) {
	wsHub.lock.Lock()
	defer wsHub.lock.Unlock()

	for conn := range wsHub.clients[lobbyID] {
		if err := conn.WriteJSON(msg); err != nil {
			log.Println("Broadcast error:", err)
			conn.Close()
			delete(wsHub.clients[lobbyID], conn)
		}
		log.Printf("Broadcast msg: %s to lobby: %s", msg, lobbyID)
	}
}

func HandleBroadcastTest(w http.ResponseWriter, r *http.Request) {
	lobbyID := r.URL.Query().Get("lobby")
	if lobbyID == "" {
		http.Error(w, "Missing Lobby ID", http.StatusBadRequest)
		return
	}

	msg := WSMessage{
		Type:    "broadcast",
		Payload: "Hello from server!",
	}

	BroadcastToLobby(lobbyID, msg)
	w.Write([]byte("Broadcast sent"))
}
