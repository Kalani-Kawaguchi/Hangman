package session

import (
	"errors"
	"fmt"
	"github.com/Kalani-Kawaguchi/Hangman/internal/game"
	"math/rand"
	"sync"
	"time"
)

type LobbyState string

const (
	StateWaiting LobbyState = "waiting"
	StateReady   LobbyState = "ready"
	StatePlaying LobbyState = "playing"
	StateEnded   LobbyState = "ended"
)

type Lobby struct {
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Player1 string     `json:"player1"`
	Player2 string     `json:"player2,omitempty"`
	State   LobbyState `json:"state"`
	Created time.Time  `json:"created"`
	Game1   game.Game
	Game2   game.Game
}

// Might move this somewhere else
type WordRequest struct {
	Word string `json:"word"`
}

// Thread-safe map to store active lobbies
var (
	lobbies   = make(map[string]*Lobby)
	lobbiesMu sync.Mutex
)

// CreateLobby initializes a new lobby and returns it
func CreateLobby(name string, player1 string) *Lobby {
	lobbiesMu.Lock()
	defer lobbiesMu.Unlock()

	id := generateLobbyID()
	lobby := &Lobby{
		ID:      id,
		Name:    name,
		Player1: player1,
		State:   StateWaiting,
		Created: time.Now(),
	}

	lobbies[id] = lobby

	return lobby
}

// JoinLobby assigns a guesser to an existing lobby
func JoinLobby(lobbyID, player2 string) (*Lobby, error) {
	lobbiesMu.Lock()
	defer lobbiesMu.Unlock()
	lobby, exists := lobbies[lobbyID]
	if !exists {
		return nil, errors.New("lobby not found")
	}
	if lobby.Player2 != "" {
		return nil, errors.New("lobby already full")
	}

	lobby.Player2 = player2
	lobby.State = StateReady
	return lobby, nil
}

// GetLobby returns a pointer to the lobby if it exists
func GetLobby(lobbyID string) (*Lobby, error) {
	lobbiesMu.Lock()
	defer lobbiesMu.Unlock()
	fmt.Println("Available lobbies:")
	for id := range lobbies {
		fmt.Println("-", id)
	}

	lobby, ok := lobbies[lobbyID]
	if !ok {
		fmt.Println(lobbyID)
		return nil, errors.New(lobbyID)
	}
	return lobby, nil
}

// Helper to generate a random lobby ID
func generateLobbyID() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6
	rand.New(rand.NewSource(time.Now().UnixNano()))
	id := make([]byte, length)
	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}
	return string(id)
}
