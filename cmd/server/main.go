package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Kalani-Kawaguchi/Hangman/internal/game"
	"github.com/Kalani-Kawaguchi/Hangman/internal/session"
	"github.com/Kalani-Kawaguchi/Hangman/internal/ws"
	"github.com/gorilla/mux"
)

// Request/Response Structs
type CreateLobbyRequest struct {
	LobbyName string `json:"lobby_name"`
	HostName  string `json:"host_name"`
}

type JoinLobbyRequest struct {
	LobbyID    string `json:"lobby_id"`
	PlayerName string `json:"player_name"`
}

type LetterRequest struct {
	Letter string `json:"guess"`
}

func main() {
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/", handleRoot)
	r.HandleFunc("/create-lobby", handleCreateLobby).Methods("POST")
	r.HandleFunc("/join-lobby", handleJoinLobby).Methods("POST")
	r.HandleFunc("/choose-word", handleChooseWord).Methods("POST")
	r.HandleFunc("/guess-letter", handleGuessLetter).Methods("POST")
	r.HandleFunc("/lobby/{id}", handleGetLobby).Methods("GET")
	r.HandleFunc("/list-lobbies", handleListLobbies).Methods("GET")
	r.HandleFunc("/leave-lobby", handleLeaveLobby).Methods("Post")
	r.HandleFunc("/ws", ws.HandleWebSocket)

	// Start server
	log.Println("Hangman running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hangman Multiplayer Game")
}

// Handlers
func handleCreateLobby(w http.ResponseWriter, r *http.Request) {
	var req CreateLobbyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "player",
		Value: req.HostName,
	})

	lobby := session.CreateLobby(req.LobbyName, req.HostName)

	http.SetCookie(w, &http.Cookie{
		Name:  "lobby",
		Value: lobby.ID,
	})

	json.NewEncoder(w).Encode(lobby)
	fmt.Fprintf(w, "Created A Lobby")
}

func handleJoinLobby(w http.ResponseWriter, r *http.Request) {
	var req JoinLobbyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	lobby, err := session.JoinLobby(req.LobbyID, req.PlayerName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "player",
		Value: req.PlayerName,
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "lobby",
		Value: req.LobbyID,
	})

	json.NewEncoder(w).Encode(lobby)
	fmt.Fprintf(w, "Joined lobby")
}

func handleChooseWord(w http.ResponseWriter, r *http.Request) {
	var req session.WordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	word := req.Word

	if !game.ValidateWord(word, w) {
		return
	}

	player_cookie, err := r.Cookie("player")
	if err != nil {
		http.Error(w, "Player not identified", http.StatusUnauthorized)
		return
	}
	playerName := player_cookie.Value

	lobby_cookie, err := r.Cookie("lobby")
	if err != nil {
		http.Error(w, "Lobby not identified", http.StatusUnauthorized)
		return
	}
	lobby := lobby_cookie.Value
	fmt.Fprintf(w, "Lobby: %s \n", lobby)

	lobby_pointer, exists := session.GetLobby(lobby)
	if exists != nil {
		http.Error(w, "Lobby not identified", http.StatusUnauthorized)
		return
	}

	if playerName == lobby_pointer.Player1 {
		lobby_pointer.Game1 = game.NewGame(word)
		fmt.Fprintf(w, "Word: '%s' chosen for %s. \n", word, lobby_pointer.Player2)
		lobby_pointer.Game1Ready = true

	} else if playerName == lobby_pointer.Player2 {
		lobby_pointer.Game2 = game.NewGame(word)
		fmt.Fprintf(w, "Word: '%s' chosen for %s. \n", word, lobby_pointer.Player1)
		lobby_pointer.Game2Ready = true
	}

	fmt.Fprintf(w, "Game created successfully.")

	if lobby_pointer.Game1Ready && lobby_pointer.Game2Ready {
		lobby_pointer.State = session.StateReady
	}
}

func handleGuessLetter(w http.ResponseWriter, r *http.Request) {
	var req LetterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	player_cookie, err := r.Cookie("player")
	if err != nil {
		http.Error(w, "Player not identified", http.StatusUnauthorized)
		return
	}
	playerName := player_cookie.Value

	lobby_cookie, err := r.Cookie("lobby")
	if err != nil {
		http.Error(w, "Lobby not identified", http.StatusUnauthorized)
		return
	}
	lobby := lobby_cookie.Value

	lobby_pointer, exists := session.GetLobby(lobby)
	if exists != nil {
		http.Error(w, "Lobby not identified", http.StatusUnauthorized)
		return
	}

	if lobby_pointer.State == session.StateWaiting {
		http.Error(w, "Lobby not ready", http.StatusUnauthorized)
		return
	}

	letter := req.Letter

	if playerName == lobby_pointer.Player1 {
		if lobby_pointer.Game2.WinOrLost(w) {
			return
		}
	} else if playerName == lobby_pointer.Player2 {
		if lobby_pointer.Game1.WinOrLost(w) {
			return
		}
	}

	if len(letter) != 1 {
		http.Error(w, "Enter a single letter.", http.StatusBadRequest)
		return
	}

	if playerName == lobby_pointer.Player1 {
		lobby_pointer.Game2.Guess(rune(letter[0]), w)
	} else if playerName == lobby_pointer.Player2 {
		lobby_pointer.Game1.Guess(rune(letter[0]), w)
	}
}

func handleGetLobby(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	lobby, err := session.GetLobby(id)
	if err != nil {
		http.Error(w, "lobby not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(lobby)
}

func handleListLobbies(w http.ResponseWriter, r *http.Request) {
	lobbies := session.GetLobbyList()
	json.NewEncoder(w).Encode(lobbies)
}

func handleLeaveLobby(w http.ResponseWriter, r *http.Request) {

	// Get player name and lobby id from cookies
	player_cookie, err := r.Cookie("player")
	if err != nil {
		http.Error(w, "Player not identified", http.StatusUnauthorized)
		return
	}
	player_name := player_cookie.Value

	lobby_cookie, err := r.Cookie("lobby")
	if err != nil {
		http.Error(w, "Lobby not identified", http.StatusUnauthorized)
		return
	}
	lobby_id := lobby_cookie.Value

	// Get lobby using lobbyid from cookie
	lobby, err := session.GetLobby(lobby_id)
	if err != nil {
		http.Error(w, "Error getting Lobby with specified ID.", http.StatusInternalServerError)
		return
	}

	// Check which player is trying to leave
	if player_name == lobby.Player1 {
		// kick out player2 if they exist
		if lobby.Player2 != "" {
			// We'll need to use websocket to notify player2 that they've been kicked out and lobby deleted
			fmt.Fprintln(w, "Player 2 kicked")
		}

		// Clear player1s lobby cookies
		http.SetCookie(w, &http.Cookie{
			Name:   "lobby",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		// delete lobby from lobbies
		session.DeleteLobby(lobby_id)
		fmt.Fprintln(w, "Lobby deleted. Host has left.")
		return

	} else if player_name == lobby.Player2 {
		// just have player2 leave the lobby
		lobby.Player2 = ""

		// update player2 lobby cookie and lobby player2 info
		http.SetCookie(w, &http.Cookie{
			Name:   "lobby",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		lobby.State = session.StateWaiting
		fmt.Fprintln(w, "You have left lobby.")
		return

	} else {
		// return error, since this player shouldn't be in this lobby
		http.Error(w, "You are not part of this lobby.", http.StatusUnauthorized)
		return
	}
}
