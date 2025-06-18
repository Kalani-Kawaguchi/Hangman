package ws

import (
	"sync"

	"github.com/Kalani-Kawaguchi/Hangman/internal/session"
)

type Hub struct {
	Lobbies map[string]*session.Lobby
	Lock    sync.Mutex
}
