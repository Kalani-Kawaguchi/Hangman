package game

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"unicode"
)

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

func ValidateLetter(r rune) bool {
	if !unicode.IsLetter(r) || (r > unicode.MaxASCII || (r < 'A' || (r > 'Z' && r < 'a') || r > 'z')) {
		log.Print("Must contain only english letters")
		return false
	}

	return true
}

func ValidateWord(word string) bool {
	if len(word) < 1 {
		log.Print("Word must be at least 1 letter")
		return false
	}

	for _, r := range word {
		if !ValidateLetter(r) {
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
		Letters:        setLetters(),
		GuessedLetters: make([]rune, 0, 26),
		Status:         InProgress,
	}
}

func setLetters() map[rune]bool {
	m := make(map[rune]bool, 26)
	for i := 97; i <= 122; i++ {
		m[rune(i)] = false
	}
	return m
}

func (g *Game) Guess(letter rune) bool {
	if g.Status != InProgress {
		return false
	}

	if !ValidateLetter(letter) {
		return false
	}

	letter = unicode.ToLower(letter)

	if g.Letters[letter] {
		log.Print("Letter already guessed")
		return false
	}

	g.Letters[letter] = true
	g.GuessedLetters = append(g.GuessedLetters, letter)
	sort.Slice(g.GuessedLetters, func(i, j int) bool {
		return g.GuessedLetters[i] < g.GuessedLetters[j]
	})

	if strings.ContainsRune(g.Word, letter) {
		g.updateMaskedWord(letter)
		g.checkGameStatus()
	} else {
		g.AttemptsLeft--
		g.checkGameStatus()
	}

	return true
}

func (g *Game) updateMaskedWord(letter rune) {
	for i, c := range g.Word {
		if c == letter {
			g.Revealed[i] = letter
		}
	}
}

func (g *Game) revealMaskedWord() {
	for i, c := range g.Word {
		g.Revealed[i] = c
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

func (g *Game) WinOrLost() bool {
	if g.Status == Won {
		log.Print("You win")
		return true
	}

	if g.Status == Lost {
		log.Print("You lost")
		g.revealMaskedWord()
		return true
	}

	return false
}

func (g *Game) DisplayState() {
	fmt.Printf("Word: %s\n", string(g.Revealed))
	fmt.Printf("Guesses Left: %d\n", g.AttemptsLeft)
	fmt.Printf("Guessed letters: ")
	for _, letter := range g.GuessedLetters {
		fmt.Printf("%c ", letter)
	}
}
