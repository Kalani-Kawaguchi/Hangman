package ws

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/Kalani-Kawaguchi/Hangman/internal/session"
	"github.com/gorilla/websocket"
)

type Hub struct {
	Lobbies map[string]*session.Lobby
	Lock    sync.Mutex
}

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
	conn, lobbyID, err := setupWebSocket(w, r)
	if err != nil {
		log.Println("Setup error:", err)
		return
	}
	defer cleanupConnection(lobbyID, conn)

	for {
		var msg WSMessage
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("WebSocket read error:", err)
			break
		}
		handleMessage(conn, lobbyID, msg)

	}
}

func setupWebSocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, string, error) {
	lobbyID := r.URL.Query().Get("lobby")
	if lobbyID == "" {
		http.Error(w, "Missing lobby ID", http.StatusBadRequest)
		return nil, "", errors.New("missing lobby ID")
	}

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, "", err
	}

	wsHub.Lock.Lock()
	defer wsHub.Lock.Unlock()

	lobby, err := session.GetLobby(lobbyID)
	if err != nil {
		return nil, "", err
	}
	if _, exists := wsHub.Lobbies[lobbyID]; !exists {
		wsHub.Lobbies[lobbyID] = lobby
	}

	lobby.ConnLock.Lock()
	lobby.Clients[conn] = true
	lobby.ConnLock.Unlock()

	log.Println("Added connection to lobby")
	return conn, lobbyID, nil
}

func handleMessage(conn *websocket.Conn, lobbyID string, msg WSMessage) {
	log.Printf("Received from %s: %v\n", lobbyID, msg)
	data := map[string]string{"type": "update", "word": "hangman"}
	conn.WriteJSON(data)
	switch msg.Type {
	case "update":
		handleUpdate(conn, lobbyID, msg.Payload)
	case "guess":
		handleGuess(conn, lobbyID, msg.Payload)
	// add the other message type cases...
	default:
		log.Println("Unknown message type:", msg.Type)
	}
}

func handleUpdate(conn *websocket.Conn, lobbyID string, payload interface{}) {
	return
}

func handleGuess(conn *websocket.Conn, lobbyID string, payload interface{}) {
	return
}

func cleanupConnection(lobbyID string, conn *websocket.Conn) {
	wsHub.Lock.Lock()
	defer wsHub.Lock.Unlock()

	delete(wsHub.Lobbies[lobbyID].Clients, conn)
	conn.Close()
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
