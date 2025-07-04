package session

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/Kalani-Kawaguchi/Hangman/internal/game"
	"github.com/gorilla/websocket"
)

type LobbyState string

const (
	StateWaiting LobbyState = "waiting"
	StateReady   LobbyState = "ready"
	StatePlaying LobbyState = "playing"
	StateEnded   LobbyState = "ended"
)

type Lobby struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Player1          string `json:"player1"`
	Player2          string `json:"player2,omitempty"`
	Player1ID        string
	Player2ID        string
	State            LobbyState `json:"state"`
	Created          time.Time  `json:"created"`
	Game1            game.Game
	Game2            game.Game
	Game1Ready       bool
	Game2Ready       bool
	Clients          map[*websocket.Conn]string // active WebSocket clients. Client: PlayerID
	ConnLock         sync.Mutex                 // protects Clients map
	Player1Restarted bool
	Player2Restarted bool
	PlayerCount      string
}

// Might move this somewhere else
type WordRequest struct {
	Word string `json:"word"`
}

type LobbySummary struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	State       LobbyState `json:"state"`
	Player1     string     `json:"player1"`
	Player2     string     `json:"player2"`
	PlayerCount string     `json:"playerCount"`
}

// Thread-safe map to store active lobbies
var (
	lobbies   = make(map[string]*Lobby)
	lobbiesMu sync.Mutex
)

// CreateLobby initializes a new lobby and returns it
func CreateLobby(name string) *Lobby {
	lobbiesMu.Lock()
	defer lobbiesMu.Unlock()

	id := GenerateID()
	lobby := &Lobby{
		ID:          id,
		Name:        name,
		State:       StateWaiting,
		Created:     time.Now(),
		Game1Ready:  false,
		Game2Ready:  false,
		Clients:     make(map[*websocket.Conn]string),
		PlayerCount: "1",
	}

	lobbies[id] = lobby

	return lobby
}

// JoinLobby assigns a player to an existing lobby
func JoinLobby(lobbyID, playerName string, playerID string) (*Lobby, error) {
	lobbiesMu.Lock()
	defer lobbiesMu.Unlock()
	lobby, exists := lobbies[lobbyID]
	if !exists {
		return nil, errors.New("lobby not found")
	}

	// Check which Lobby player to assign to
	if lobby.Player1 == "" {
		lobby.Player1 = playerName
		lobby.Player1ID = playerID
	} else if lobby.State != StateWaiting {
		return nil, errors.New("lobby is busy")
	} else if lobby.Player2 == "" {
		lobby.Player2 = playerName
		lobby.Player2ID = playerID
		lobby.PlayerCount = "2"
	} else {
		return nil, errors.New("lobby already full")
	}

	return lobby, nil
}

// GetLobby returns a pointer to the lobby if it exists
func GetLobby(lobbyID string) (*Lobby, error) {
	lobbiesMu.Lock()
	defer lobbiesMu.Unlock()

	lobby, ok := lobbies[lobbyID]
	if !ok {
		fmt.Println(lobbyID)
		return nil, errors.New(lobbyID)
	}
	return lobby, nil
}

func GetLobbyList() []LobbySummary {
	lobbiesMu.Lock()
	defer lobbiesMu.Unlock()

	var availableLobbies []LobbySummary
	for id, lobby := range lobbies {
		availableLobbies = append(availableLobbies, LobbySummary{ID: id, Name: lobby.Name, State: lobby.State, Player1: lobby.Player1, Player2: lobby.Player2, PlayerCount: lobby.PlayerCount})
	}
	return availableLobbies
}

// Helper to generate a random lobby or Player ID
func GenerateID() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6
	rand.New(rand.NewSource(time.Now().UnixNano()))
	id := make([]byte, length)
	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}
	return string(id)
}

func DeleteLobby(lobby_id string) {
	lobbiesMu.Lock()
	defer lobbiesMu.Unlock()
	delete(lobbies, lobby_id)
}
