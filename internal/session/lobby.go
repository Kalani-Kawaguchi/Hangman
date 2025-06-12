package session

import (
	"errors"
	"fmt"
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
	ID       string     `json:"id"`
	Host     string     `json:"host"`
	Guesser  string     `json:"guesser,omitempty"`
	State    LobbyState `json:"state"`
	Created  time.Time  `json:"created"`
	Word     string     `json:"-"` // hidden from API response
	Attempts int        `json:"attempts"`
	Guesses  []string   `json:"guesses"`
}

// Thread-safe map to store active lobbies
var (
	lobbies   = make(map[string]*Lobby)
	lobbiesMu sync.Mutex
)

// CreateLobby initializes a new lobby and returns it
func CreateLobby(hostName string) *Lobby {
	lobbiesMu.Lock()
	defer lobbiesMu.Unlock()

	id := generateLobbyID()
	lobby := &Lobby{
		ID:      id,
		Host:    hostName,
		State:   StateWaiting,
		Created: time.Now(),
	}

	lobbies[id] = lobby
	return lobby
}

// JoinLobby assigns a guesser to an existing lobby
func JoinLobby(lobbyID, guesserName string) (*Lobby, error) {
	lobbiesMu.Lock()
	defer lobbiesMu.Unlock()
	fmt.Printf("lobbyID = [%s]\n", lobbyID)
	fmt.Println("WHYYY")
	lobby, exists := lobbies[lobbyID]
	if !exists {
		return nil, errors.New("lobby not found")
	}
	if lobby.Guesser != "" {
		return nil, errors.New("lobby already has a guesser")
	}

	lobby.Guesser = guesserName
	lobby.State = StateReady
	return lobby, nil
}

// GetLobby returns a pointer to the lobby if it exists
func GetLobby(lobbyID string) (*Lobby, error) {
	lobbiesMu.Lock()
	defer lobbiesMu.Unlock()

	lobby, ok := lobbies[lobbyID]
	if !ok {
		return nil, errors.New("lobby not found")
	}
	return lobby, nil
}

// Helper to generate a random lobby ID
func generateLobbyID() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6
	rand.Seed(time.Now().UnixNano())

	id := make([]byte, length)
	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}
	return string(id)
}
