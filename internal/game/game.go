package game

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

}
