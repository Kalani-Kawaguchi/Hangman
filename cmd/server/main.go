package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Kalani-Kawaguchi/Hangman/internal/game"
	"github.com/Kalani-Kawaguchi/Hangman/internal/session"
	"github.com/Kalani-Kawaguchi/Hangman/internal/ws"
	"github.com/gorilla/handlers"
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

type LobbyRequest struct {
	LobbyID string `json:"lobby_id"`
}

type LetterRequest struct {
	Letter string `json:"guess"`
}

func newRest() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handleRoot)
	r.HandleFunc("/create-lobby", handleCreateLobby).Methods("POST")
	r.HandleFunc("/join-lobby", handleJoinLobby).Methods("POST")
	r.HandleFunc("/choose-word", handleChooseWord).Methods("POST")
	r.HandleFunc("/guess-letter", handleGuessLetter).Methods("POST")
	r.HandleFunc("/lobby/{id}", handleGetLobby).Methods("GET")
	r.HandleFunc("/list-lobbies", handleListLobbies).Methods("GET")
	r.HandleFunc("/list-games", handleListGames).Methods("POST")
	r.HandleFunc("/leave-lobby", handleLeaveLobby).Methods("POST")
	r.HandleFunc("/player-role", handlePlayerRole).Methods("GET")
	r.HandleFunc("/ws", ws.HandleWebSocket)
	r.HandleFunc("/lobby-state", HandleLobbyState).Methods("GET")

	return r
}

func main() {
	r := newRest()
	origins := handlers.AllowedOrigins([]string{"https://gohangman.vercel.app", "http://localhost:3000"})
	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	methods := handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS"})
	credentials := handlers.AllowCredentials()

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/").Handler(fs)

	// Start server
	log.Println("Hangman running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(origins, methods, headers, credentials)(r)))
}

func HandleLobbyState(w http.ResponseWriter, r *http.Request) {
	lobbyID := r.URL.Query().Get("lobby")
	if lobbyID == "" {
		http.Error(w, "Missing lobby ID", http.StatusBadRequest)
		return
	}
	lobby, err := session.GetLobby(lobbyID)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}
	lobbyState := lobby.State
	player1Exists := lobby.Player1Exists
	player2Exists := lobby.Player2Exists
	player1Name := lobby.Player1
	player2Name := lobby.Player2
	player1Instruction := lobby.Player1Instruction
	player2Instruction := lobby.Player2Instruction
	player1OppInstruction := lobby.Player1OppInstruction
	player2OppInstruction := lobby.Player2OppInstruction
	player1Restarted := lobby.Player1Restarted
	player2Restarted := lobby.Player2Restarted
	player1RevealedWord := string(lobby.Game2.Revealed)
	player2RevealedWord := string(lobby.Game1.Revealed)
	player1AttemptsLeft := lobby.Game2.AttemptsLeft
	player2AttemptsLeft := lobby.Game1.AttemptsLeft
	player1Ready := lobby.Game2Ready
	player2Ready := lobby.Game1Ready
	player1Guessed := string(lobby.Game2.GuessedLetters)
	player2Guessed := string(lobby.Game1.GuessedLetters)
	json.NewEncoder(w).Encode(map[string]any{
		"state":                 string(lobbyState),
		"player1Exists":         player1Exists,
		"player2Exists":         player2Exists,
		"player1Name":           player1Name,
		"player2Name":           player2Name,
		"player1Instruction":    player1Instruction,
		"player2Instruction":    player2Instruction,
		"player1OppInstruction": player1OppInstruction,
		"player2OppInstruction": player2OppInstruction,
		"player1Restarted":      player1Restarted,
		"player2Restarted":      player2Restarted,
		"player1RevealedWord":   player1RevealedWord,
		"player2RevealedWord":   player2RevealedWord,
		"player1AttemptsLeft":   player1AttemptsLeft,
		"player2AttemptsLeft":   player2AttemptsLeft,
		"player1Ready":          player1Ready,
		"player2Ready":          player2Ready,
		"player1Guessed":        player1Guessed,
		"player2Guessed":        player2Guessed,
	})
}

func getLobbyFromCookies(r *http.Request) (*session.Lobby, string, string, error) {
	playerCookie, err := r.Cookie("player")
	if err != nil {
		return nil, "", "", fmt.Errorf("player not identified")
	}
	playerIdCookie, err := r.Cookie("id")
	if err != nil {
		return nil, "", "", fmt.Errorf("player id not identified")
	}
	lobbyCookie, err := r.Cookie("lobby")
	if err != nil {
		return nil, "", "", fmt.Errorf("lobby not identified")
	}

	lobby, err := session.GetLobby(lobbyCookie.Value)
	if err != nil {
		return nil, "", "", fmt.Errorf("lobby not found")
	}

	return lobby, playerCookie.Value, playerIdCookie.Value, nil
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
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

	// Generate Unique Player ID and save it as a cookie
	playerID := session.GenerateID()
	http.SetCookie(w, &http.Cookie{
		Name:  "id",
		Value: playerID,
	})

	// lobby is a pointer to the newly created Lobby
	lobby := session.CreateLobby(req.LobbyName)

	http.SetCookie(w, &http.Cookie{
		Name:  "lobby",
		Value: lobby.ID,
	})

	// assign player to lobby
	_, err := session.JoinLobby(lobby.ID, req.HostName, playerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Created A Lobby: %s. Host: %s", lobby.ID, playerID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":       lobby.ID,
		"playerID": playerID,
	})
}

func handleJoinLobby(w http.ResponseWriter, r *http.Request) {
	var req JoinLobbyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "player",
		Value: req.PlayerName,
	})

	// Generate Unique Player ID and save it as a cookie and set the lobby player ID
	playerID := session.GenerateID()
	http.SetCookie(w, &http.Cookie{
		Name:  "id",
		Value: playerID,
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "lobby",
		Value: req.LobbyID,
	})

	_, err := session.JoinLobby(req.LobbyID, req.PlayerName, playerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Notify lobby that a player has joined
	fmt.Println("Broadcasting a player join")
	ws.BroadcastToLobby(req.LobbyID, "join")

	json.NewEncoder(w).Encode(map[string]string{
		"playerID": playerID,
	})
}

func handleChooseWord(w http.ResponseWriter, r *http.Request) {
	var req session.WordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	word := req.Word

	if !game.ValidateWord(word) {
		return
	}

	lobby_pointer, playerName, _, err := getLobbyFromCookies(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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

	lobby_pointer, playerName, _, err := getLobbyFromCookies(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if lobby_pointer.State == session.StateWaiting {
		http.Error(w, "Lobby not ready", http.StatusUnauthorized)
		return
	}

	letter := req.Letter

	if playerName == lobby_pointer.Player1 {
		if lobby_pointer.Game2.WinOrLost() {
			return
		}
	} else if playerName == lobby_pointer.Player2 {
		if lobby_pointer.Game1.WinOrLost() {
			return
		}
	}

	if len(letter) != 1 {
		http.Error(w, "Enter a single letter.", http.StatusBadRequest)
		return
	}

	if playerName == lobby_pointer.Player1 {
		lobby_pointer.Game2.Guess(rune(letter[0]))
	} else if playerName == lobby_pointer.Player2 {
		lobby_pointer.Game1.Guess(rune(letter[0]))
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
	lobby, _, playerID, err := getLobbyFromCookies(r)
	if err != nil {
		log.Println("Error getting lobby from cookies")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	lobby_id := lobby.ID
	log.Printf("Leaving Lobby: %s and Player: %s", lobby_id, playerID)

	// Check which player is trying to leave
	if playerID == lobby.Player1ID {
		lobby.PlayerCount = "0"
		ws.BroadcastToLobby(lobby_id, "closeAll")
		return

	} else if playerID == lobby.Player2ID {
		// just have player2 leave the lobby
		lobby.Player2 = ""
		lobby.PlayerCount = "1"

		// update player2 lobby cookie and lobby player2 info
		http.SetCookie(w, &http.Cookie{
			Name:   "lobby",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		ws.BroadcastToLobby(lobby_id, "closeOne")
		return

	} else {
		// return error, since this player shouldn't be in this lobby
		http.Error(w, "You are not part of this lobby.", http.StatusUnauthorized)
		return
	}
}

func handleListGames(w http.ResponseWriter, r *http.Request) {
	var req LobbyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	lobby, err := session.GetLobby(req.LobbyID)
	if err != nil {
		http.Error(w, "lobby not found", http.StatusNotFound)
		return
	}

	resp := map[string]interface{}{
		"game1":      lobby.Game1,
		"game2":      lobby.Game2,
		"game1Ready": strconv.FormatBool(lobby.Game1Ready),
		"game2Ready": strconv.FormatBool(lobby.Game2Ready),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handlePlayerRole(w http.ResponseWriter, r *http.Request) {
	lobbyID := r.URL.Query().Get("lobby")
	playerID := r.URL.Query().Get("id")

	lobby, err := session.GetLobby(lobbyID)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	role := "guest"
	name := lobby.Player2
	opponent := lobby.Player1
	if playerID == lobby.Player1ID {
		role = "host"
		name = lobby.Player1
		opponent = lobby.Player2
	}

	json.NewEncoder(w).Encode(map[string]string{
		"role":     role,
		"name":     name,
		"opponent": opponent,
	})
}
