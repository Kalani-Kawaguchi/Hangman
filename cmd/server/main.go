package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Kalani-Kawaguchi/Hangman/internal/game"
	"github.com/Kalani-Kawaguchi/Hangman/internal/session"
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

	if !game.ValidateWord(word) {
		http.Error(w, "invalid word try again", http.StatusBadRequest)
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
	}

	if playerName == lobby_pointer.Player1 {
		lobby_pointer.Game1 = game.NewGame(word)
		fmt.Fprintf(w, "Word: '%s' chosen for %s. \n", word, lobby_pointer.Player2)
	} else if playerName == lobby_pointer.Player2 {
		lobby_pointer.Game2 = game.NewGame(word)
		fmt.Fprintf(w, "Word: '%s' chosen for %s. \n", word, lobby_pointer.Player1)
	}

	fmt.Fprintf(w, "Game created successfully.")
}

func handleGuessLetter(w http.ResponseWriter, r *http.Request) {
	var req LetterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	letter := req.Letter
	if len(letter) != 1 {
		http.Error(w, "Enter a single letter.", http.StatusBadRequest)
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
	}

	if playerName == lobby_pointer.Player1 {
		if !lobby_pointer.Game2.Guess(rune(letter[0])) {
			http.Error(w, "Letter already guessed", http.StatusBadRequest)
		} else {
			lobby_pointer.Game2.Guess(rune(letter[0]))
			fmt.Fprintf(w, "Letter, %s, guessed successful \n", letter)
		}
	} else if playerName == lobby_pointer.Player2 {
		if !lobby_pointer.Game1.Guess(rune(letter[0])) {
			http.Error(w, "Letter already guessed", http.StatusBadRequest)
		} else {
			lobby_pointer.Game1.Guess(rune(letter[0]))
			fmt.Fprintf(w, "Letter, %s, guessed successful \n", letter)
		}
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
