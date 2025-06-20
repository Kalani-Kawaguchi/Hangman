package ws

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Kalani-Kawaguchi/Hangman/internal/game"
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

	playerCookie, err := r.Cookie("player")
	if err != nil {
		return nil, "", fmt.Errorf("player not identified")
	}
	lobby.ConnLock.Lock()
	lobby.Clients[conn] = playerCookie.Value
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
	case "submit":
		handleSubmit(conn, lobbyID, msg.Payload)
	default:
		log.Println("Unknown message type:", msg.Type)
	}
}

func handleUpdate(conn *websocket.Conn, lobbyID string, payload interface{}) {
	return
}

func handleGuess(conn *websocket.Conn, lobbyID string, payload interface{}) {
	lobby := wsHub.Lobbies[lobbyID]
	playerName := lobby.Clients[conn]
	letter, ok := payload.(string)
	if !ok {
		log.Print("Letter could not be asserted to string")
		return
	}

	if lobby.State == session.StateWaiting {
		log.Print("Lobby not ready")
		return
	}

	data := map[string]string{"type": "update", "revealed": ""}
	// USE BROADCAST
	if playerName == lobby.Player1 {
		lobby.Game2.Guess(rune(letter[0]))
		data["revealed"] = string(lobby.Game2.Revealed)

	} else if playerName == lobby.Player2 {
		lobby.Game1.Guess(rune(letter[0]))
		data["revealed"] = string(lobby.Game1.Revealed)
	}

}

func handleSubmit(conn *websocket.Conn, lobbyID string, payload interface{}) {
	lobby := wsHub.Lobbies[lobbyID]
	playerName := lobby.Clients[conn]
	word, ok := payload.(string)
	if !ok {
		log.Print("Letter could not be asserted to string")
		return
	}

	if playerName == lobby.Player1 {
		if game.ValidateWord(word) {
			lobby.Game1 = game.NewGame(word)
			lobby.Game1Ready = true
		}
	} else if playerName == lobby.Player2 {
		if game.ValidateWord(word) {
			lobby.Game2 = game.NewGame(word)
			lobby.Game2Ready = true
		}
	}

	if lobby.Game1Ready && lobby.Game2Ready {
		lobby.State = session.StateReady
		// USE BROADCAST
		// for c, name := range lobby.Clients {
		// 	data := map[string]string{"type": "update", "revealed": ""}
		// 	if name == lobby.Player1 {
		// 		data["revealed"] = string(lobby.Game2.Revealed)
		// 	} else if name == lobby.Player2 {
		// 		data["revealed"] = string(lobby.Game1.Revealed)
		// 	}
		// 	c.WriteJSON(data)
		// 	start_message := map[string]string{"type": "start_game", "start": "x"}
		// 	c.WriteJSON(start_message)
		// }
	}
}

func cleanupConnection(lobbyID string, conn *websocket.Conn) {
	wsHub.Lock.Lock()
	defer wsHub.Lock.Unlock()

	lobby, ok := wsHub.Lobbies[lobbyID]
	if !ok {
		return
	}

	// Remove the WebSocket connection
	lobby.ConnLock.Lock()
	delete(lobby.Clients, conn)
	lobby.ConnLock.Unlock()
	conn.Close()

	if lobby.Player1 == lobby.Clients[conn] {
		log.Printf("Host left lobby %s", lobbyID)
		// You could auto-delete the lobby or notify remaining player
		delete(wsHub.Lobbies, lobbyID)
		session.DeleteLobby(lobbyID) // delete lobby from session layer too
	} else if lobby.Player2 == lobby.Clients[conn] {
		log.Printf("Guest left lobby %s", lobbyID)
	}
	// If no more clients are connected, delete the lobby
	if len(lobby.Clients) == 0 {
		log.Printf("Lobby %s is empty. Deleting it.", lobbyID)
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
