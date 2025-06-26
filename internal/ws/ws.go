package ws

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
		log.Println("no connection")
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

	// Get player's "id" cookie and save it in the client connection
	playerID, err := r.Cookie("id")
	if err != nil {
		return nil, "", fmt.Errorf("player not identified")
	}
	lobby.ConnLock.Lock()
	lobby.Clients[conn] = playerID.Value
	lobby.ConnLock.Unlock()

	log.Println("Added connection to lobby")
	return conn, lobbyID, nil
}

func handleMessage(conn *websocket.Conn, lobbyID string, msg WSMessage) {
	log.Printf("Received from %s: %v\n", lobbyID, msg)
	// data := map[string]string{"type": "update", "word": "hangman"}
	// conn.WriteJSON(data)

	lobby, err := session.GetLobby(lobbyID)
	if err != nil {
		return
	}
	playerID := lobby.Clients[conn]

	switch msg.Type {
	case "update":
		handleUpdate(conn, lobbyID, playerID, msg.Payload)
	case "guess":
		handleGuess(conn, lobbyID, playerID, msg.Payload)
	case "submit":
		handleSubmit(conn, lobbyID, playerID, msg.Payload)
	case "restart":
		handleRestart(lobbyID)
	default:
		log.Println("Unknown message type:", msg.Type)
	}
}

func handleRestart(lobbyID string) {
	lobby := wsHub.Lobbies[lobbyID]
	lobby.State = session.StateWaiting
}

func handleUpdate(conn *websocket.Conn, lobbyID string, playerID string, payload interface{}) {
	return
}

func handleGuess(conn *websocket.Conn, lobbyID string, playerID string, payload interface{}) {
	lobby := wsHub.Lobbies[lobbyID]
	letter, ok := payload.(string)
	if !ok {
		log.Print("Letter could not be asserted to string")
		return
	}

	if lobby.State == session.StateWaiting {
		log.Print("Lobby not ready")
		return
	}

	if playerID == lobby.Player1ID {
		lobby.Game2.Guess(rune(letter[0]))
		sendWinLost(conn, lobby.Game2)
	} else if playerID == lobby.Player2ID {
		lobby.Game1.Guess(rune(letter[0]))
		sendWinLost(conn, lobby.Game1)
	}

	// check if player2 is in lobby and both games are finished OR Only player1 is in the lobby and their game is finished
	if (lobby.Player2 != "" && (lobby.Game1.Status != game.InProgress && lobby.Game2.Status != game.InProgress)) ||
		(lobby.Game2.Status != game.InProgress && lobby.Player2 == "") {
		lobby.State = session.StateEnded
		BroadcastToLobby(lobbyID, "end")
		return
	}

	log.Print("handleGuess")
	BroadcastToLobby(lobbyID, "update")
}

func sendWinLost(conn *websocket.Conn, g game.Game) {
	if g.WinOrLost() {
		if g.Status == game.Won {
			data := map[string]string{"type": "win", "payload": "win"}
			conn.WriteJSON(data)
		} else if g.Status == game.Lost {
			data := map[string]string{"type": "lost", "payload": g.Word}
			conn.WriteJSON(data)
		}
	}
}

func handleSubmit(conn *websocket.Conn, lobbyID string, playerID string, payload interface{}) {
	lobby := wsHub.Lobbies[lobbyID]
	word, ok := payload.(string)
	if !ok {
		log.Print("Letter could not be asserted to string")
		return
	}

	if playerID == lobby.Player1ID {
		log.Println("checking player1ID")
		if game.ValidateWord(word) {
			log.Println("validating word for player1")
			lobby.Game1 = game.NewGame(word)
			lobby.Game1Ready = true
			log.Println("New Game1 Created")
		}
	} else if playerID == lobby.Player2ID {
		log.Println("checking player2ID")
		if game.ValidateWord(word) {
			log.Println("validating word for player2")
			lobby.Game2 = game.NewGame(word)
			lobby.Game2Ready = true
			log.Println("New Game2 created")
		}
	}

	if lobby.Game1Ready && lobby.Game2Ready {
		lobby.State = session.StateReady
		BroadcastToLobby(lobbyID, "update")
		BroadcastToLobby(lobbyID, "start_game")
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

	// If no more clients are connected, delete the lobby
	if len(lobby.Clients) == 0 {
		log.Printf("Lobby %s is empty. Deleting it.", lobbyID)
		session.DeleteLobby(lobbyID)
		delete(wsHub.Lobbies, lobbyID)
	}
}

func BroadcastToLobby(lobbyID string, t string) {
	wsHub.Lock.Lock()
	defer wsHub.Lock.Unlock()

	lobby, exists := wsHub.Lobbies[lobbyID]
	if !exists {
		log.Printf("Lobby %s not found for broadcast", lobbyID)
		return
	}
	lobby.ConnLock.Lock()
	defer lobby.ConnLock.Unlock()
	for conn, id := range lobby.Clients {
		switch t {
		case "update":
			data := map[string]string{"type": "update", "revealed": "", "attempts": "6"}
			if id == lobby.Player1ID {
				data["revealed"] = string(lobby.Game2.Revealed)
				data["attempts"] = strconv.Itoa(lobby.Game2.AttemptsLeft)
				conn.WriteJSON(data)
			} else if id == lobby.Player2ID {
				data["revealed"] = string(lobby.Game1.Revealed)
				data["attempts"] = strconv.Itoa(lobby.Game1.AttemptsLeft)
				conn.WriteJSON(data)
			}
		case "start_game":
			start_message := map[string]string{"type": "start_game", "start": "x"}
			conn.WriteJSON(start_message)
		case "closeAll":
			data := map[string]string{"type": "close", "message": "close"}
			conn.WriteJSON(data)
		case "closeOne":
			if id == lobby.Player2ID {
				data := map[string]string{"type": "close", "message": "close"}
				conn.WriteJSON(data)
			}
		case "end":
			data := map[string]string{"type": "end", "message": "end"}
			conn.WriteJSON(data)
			resetLobby(lobbyID) // may have an issue with double reset
		}
		log.Printf("Broadcast msg: %s to lobby: %s", t, lobbyID)
	}
}

func resetLobby(lobbyID string) {
	lobby, err := session.GetLobby(lobbyID)
	if err != nil {
		log.Fatalln("Lobby not found while resetting")
	}

	lobby.Game1Ready = false
	lobby.Game2Ready = false
	lobby.Game1 = game.Game{}
	lobby.Game2 = game.Game{}
}
