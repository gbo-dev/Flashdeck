// main.go
package main

import (
	"fmt"
	"log"

	"go-flashcards/data"
	"go-flashcards/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO: Before git:
// Fix indentation
// Maintain consistent error handling
// No empty tags

// TODO: For fun:
// Chaos mode: random card selection from all available decks
// Random answer color toggle
// Ability randomize order of cards

// Functionality:
// Sound effects: Morse code
// Optional timer
// Smart tracking of the cards you fail
// Search for decks with 's' and '/'?
// Multichoice questions
// Emoji toggle
// Customizable colors

func main() {

	deckManager := data.NewDeckManager()

	if err := deckManager.LoadAllDecks(); err != nil {
		log.Fatal(err)
	}

	model := ui.NewModel(deckManager)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
