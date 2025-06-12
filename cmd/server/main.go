package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Kalani-Kawaguchi/Hangman/internal/session"
	"github.com/gorilla/mux"
)

// Request/Response Structs
type CreateLobbyRequest struct {
	Name string `json:"name"`
	Host string `json:"host"`
}

type JoinLobbyRequest struct {
	LobbyID string `json:"lobby_id"`
	Name    string `json:"name"`
}

func main() {
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/", handleRoot)
	r.HandleFunc("/create-lobby", handleCreateLobby).Methods("POST")
	r.HandleFunc("/join-lobby", handleJoinLobby).Methods("POST")
	r.HandleFunc("/lobby/{id}", handleGetLobby).Methods("GET")

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

	lobby := session.CreateLobby(req.Name, req.Host)
	json.NewEncoder(w).Encode(lobby)
	fmt.Fprintf(w, "Created A Lobby")
}

func handleJoinLobby(w http.ResponseWriter, r *http.Request) {
	var req JoinLobbyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	lobby, err := session.JoinLobby(req.LobbyID, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(lobby)
	fmt.Fprintf(w, "Joined lobby")
}

func handleChooseword(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// Once Host and Guesser fields are filled we need to have them both choose a word
	// Once the word is given it's passed into game struct
	// Once both games are created lobby set to ready
	// once lobby is ready both games will start
	// Look into web sockets
	return
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
