package ws

import (
	"log"
	"net/http"

	"github.com/Kalani-Kawaguchi/Hangman/internal/session"
	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

var wsHub = &Hub{
	Lobbies: make(map[string]*session.Lobby),
}

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	lobbyID := r.URL.Query().Get("lobby")
	if lobbyID == "" {
		http.Error(w, "Missing lobby ID", http.StatusBadRequest)
		return
	}

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	wsHub.Lock.Lock()
	lobby, err := session.GetLobby(lobbyID)
	if err != nil {
		log.Println("Lobby not found")
	}

	if _, exists := wsHub.Lobbies[lobbyID]; !exists {

		wsHub.Lobbies[lobbyID] = lobby
	}

	lobby.ConnLock.Lock()
	if !lobby.Clients[conn] {
		lobby.Clients[conn] = true
	}
	lobby.ConnLock.Unlock()
	log.Println("Added connection to lobby")

	wsHub.Lock.Unlock()

	defer func() {
		wsHub.Lock.Lock()
		delete(wsHub.Lobbies[lobbyID].Clients, conn)
		wsHub.Lock.Unlock()
		conn.Close()
	}()

	for {
		var msg WSMessage
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		log.Printf("Received from %s: %v\n", lobbyID, msg)

		data := map[string]string{"type": "update", "word": "hangman"}
		conn.WriteJSON(data)
	}
}

func BroadcastToLobby(lobbyID string, msg WSMessage) {
	wsHub.Lock.Lock()
	defer wsHub.Lock.Unlock()

	lobby, exists := wsHub.Lobbies[lobbyID]
	if !exists {
		log.Printf("Lobby %s not found for broadcast", lobbyID)
		return
	}
	lobby.ConnLock.Lock()
	defer lobby.ConnLock.Unlock()
	for conn := range lobby.Clients {
		if err := conn.WriteJSON(msg); err != nil {
			log.Println("Broadcast error:", err)
			conn.Close()
			delete(lobby.Clients, conn)
		}
		log.Printf("Broadcast msg: %v to lobby: %s", msg, lobbyID)
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
