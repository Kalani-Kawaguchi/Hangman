package game

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"unicode"
)

var letters = map[rune]bool{
	97:  false, // 'a'
	98:  false, // 'b'
	99:  false, // 'c'
	100: false, // 'd'
	101: false, // 'e'
	102: false, // 'f'
	103: false, // 'g'
	104: false, // 'h'
	105: false, // 'i'
	106: false, // 'j'
	107: false, // 'k'
	108: false, // 'l'
	109: false, // 'm'
	110: false, // 'n'
	111: false, // 'o'
	112: false, // 'p'
	113: false, // 'q'
	114: false, // 'r'
	115: false, // 's'
	116: false, // 't'
	117: false, // 'u'
	118: false, // 'v'
	119: false, // 'w'
	120: false, // 'x'
	121: false, // 'y'
	122: false, // 'z'
}

type Game struct {
	Word           string
	Revealed       []rune
	AttemptsLeft   int
	Letters        map[rune]bool
	GuessedLetters []rune
	Status         GameStatus
}

type GameStatus int

const (
	InProgress GameStatus = iota
	Won
	Lost
)

func ValidateLetter(r rune, w http.ResponseWriter) bool {
	if !unicode.IsLetter(r) || (r > unicode.MaxASCII || (r < 'A' || (r > 'Z' && r < 'a') || r > 'z')) {
		http.Error(w, "Must contain only English letters.", http.StatusBadRequest)
		return false
	}

	return true
}

func ValidateWord(word string, w http.ResponseWriter) bool {
	if len(word) < 5 {
		http.Error(w, "Word must be at least 5 letters.", http.StatusBadRequest)
		return false
	}

	for _, r := range word {
		if !ValidateLetter(r, w) {
			return false
		}
	}

	return true
}

func NewGame(word string) Game {
	revealed := make([]rune, len(word))
	for i := range revealed {
		revealed[i] = '_'
	}

	return Game{
		Word:           strings.ToLower(word),
		Revealed:       revealed,
		AttemptsLeft:   6,
		Letters:        letters,
		GuessedLetters: make([]rune, 0, 26),
		Status:         InProgress,
	}
}

func (g *Game) Guess(letter rune, w http.ResponseWriter) bool {
	if !ValidateLetter(letter, w) {
		return false
	}

	letter = unicode.ToLower(letter)

	if g.Letters[letter] {
		http.Error(w, "Letter already guessed", http.StatusBadRequest)
		return false
	}

	g.Letters[letter] = true
	g.GuessedLetters = append(g.GuessedLetters, letter)
	sort.Slice(g.GuessedLetters, func(i, j int) bool {
		return g.GuessedLetters[i] < g.GuessedLetters[j]
	})

	if strings.ContainsRune(g.Word, letter) {
		g.updateMaskedWord(letter)
		g.DisplayState(w)
	} else {
		g.AttemptsLeft--
		g.DisplayState(w)
	}

	g.checkGameStatus()
	return true
}

func (g *Game) updateMaskedWord(letter rune) {
	for i, c := range g.Word {
		if c == letter {
			g.Revealed[i] = letter
		}
	}
}

func (g *Game) checkGameStatus() {
	if g.AttemptsLeft <= 0 {
		g.Status = Lost
		return
	}

	if !strings.ContainsRune(string(g.Revealed), '_') {
		g.Status = Won
	}
}

func (g *Game) DisplayState(w http.ResponseWriter) {
	fmt.Fprintf(w, "Word: %s\n", string(g.Revealed))
	fmt.Fprintf(w, "Guesses Left: %d\n", g.AttemptsLeft)
	fmt.Fprintf(w, "Guessed letters: ")
	for _, letter := range g.GuessedLetters {
		fmt.Fprintf(w, "%c ", letter)
	}
}
