package game

import (
	"fmt"
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

func ValidateLetter(r rune) bool {
	if !unicode.IsLetter(r) || (r > unicode.MaxASCII || (r < 'A' || (r > 'Z' && r < 'a') || r > 'z')) {
		return false
	}

	return true
}

func ValidateWord(word string) bool {
	if len(word) < 5 {
		return false
	}

	for _, r := range word {
		ValidateLetter(r)
	}

	return true
}

func NewGame(word string, maxAttempts int) *Game {
	revealed := make([]rune, len(word))
	for i := range revealed {
		revealed[i] = '_'
	}
	return &Game{
		Word:           strings.ToLower(word),
		Revealed:       revealed,
		AttemptsLeft:   maxAttempts,
		Letters:        letters,
		GuessedLetters: make([]rune, 0, 26),
		Status:         InProgress,
	}
}

func (g *Game) Guess(letter rune) bool {
	if !ValidateLetter(letter) {
		return false
	}

	letter = unicode.ToLower(letter)

	if g.Letters[letter] {
		fmt.Println("You already guessed that letter.")
		return false
	}

	g.Letters[letter] = true
	g.GuessedLetters = append(g.GuessedLetters, letter)
	sort.Slice(g.GuessedLetters, func(i, j int) bool {
		return g.GuessedLetters[i] < g.GuessedLetters[j]
	})

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
	fmt.Printf("Guesses Left: %d\n", g.AttemptsLeft)
	fmt.Printf("Guessed letters: ")
	for _, letter := range g.GuessedLetters {
		fmt.Printf("%c ", letter)
	}
	fmt.Println()
}
