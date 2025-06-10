package game

import (
	"fmt"
	"strings"
	"unicode"
)

type Game struct {
	Word           string
	Revealed       []rune
	AttemptsLeft   int
	GuessedLetters map[rune]bool
	Status         GameStatus
}

type GameStatus int

const (
	InProgress GameStatus = iota
	Won
	Lost
)

func ValidateWord(word string) bool {
	if len(word) < 5 {
		fmt.Println("The word must be at least 5 characters.")
		return false
	}

	for _, r := range word {
		if !unicode.IsLetter(r) || (r > unicode.MaxASCII || (r < 'A' || (r > 'Z' && r < 'a') || r > 'z')) {
			fmt.Println("The word must only contain letters from the English alphabet")
			return false
		}
	}

	return true
}

func NewGame(word string, maxAttempts int) *Game {
	revealed := make([]rune, len(word))
	for i := range revealed {
		revealed[i] = '_'
	}
	return &Game{
		Word:           word,
		Revealed:       revealed,
		AttemptsLeft:   maxAttempts,
		GuessedLetters: make(map[rune]bool),
		Status:         InProgress,
	}
}

func (g *Game) Guess(letter rune) bool {
	letter = unicode.ToLower(letter)

	if g.GuessedLetters[letter] {
		fmt.Println("You already guessed that letter.")
		return false
	}

	g.GuessedLetters[letter] = true

	if strings.ContainsRune(g.Word, letter) {
		g.updateMaskedWord(letter)
	} else {
		g.AttemptsLeft--
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

func (g *Game) DisplayState() {
	fmt.Println("\nWord:", string(g.Revealed))
	fmt.Printf("Guessed Left: %d\n", g.AttemptsLeft)
	fmt.Printf("Guessed letters: ")
	for letter := range g.GuessedLetters {
		fmt.Printf("%c ", letter)
	}
	fmt.Println()
}
