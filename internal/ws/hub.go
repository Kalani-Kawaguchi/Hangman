package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/Kalani-Kawaguchi/Hangman/internal/session"
	"github.com/gorilla/websocket"
)

type Hub struct {
	Lobbies map[string]*session.Lobby
	lock    sync.Mutex
}
