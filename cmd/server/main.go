package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Kalani-Kawaguchi/Hangman/internal/game"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	word := "Gopher"
	hangmanGame := game.NewGame(word, 6)

	fmt.Println("Welcome to Hangman!")
	for hangmanGame.Status == game.InProgress {
		hangmanGame.DisplayState()
		fmt.Print("Enter a letter: ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if len(input) != 1 {
			fmt.Println("Please enter a single letter.")
			continue
		}

		hangmanGame.Guess(rune(input[0]))
	}

	hangmanGame.DisplayState()
	switch hangmanGame.Status {
	case game.Won:
		fmt.Println("You Won! The word was: ", hangmanGame.Word)
	case game.Lost:
		fmt.Println("Game Over! The word was: ", hangmanGame.Word)
	}
}
