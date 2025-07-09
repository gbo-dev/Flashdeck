// data/deck.go
package data

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Card struct {
	Question string   `yaml:"question"`
	Answer   string   `yaml:"answer"`
	Tags     []string `yaml:"tags"`
}

type Deck struct {
	ID        uuid.UUID `yaml:"id"`
	Name      string    `yaml:"name"`
	Cards     []Card    `yaml:"cards"`
	CurrentID int       `yaml:"current_id"`
}

func NewCard(q, a string, tags []string) Card {
	return Card{
		Question: q,
		Answer:   a,
		Tags:     tags,
	}
}

func NewDeck(name string) *Deck {
	return &Deck{
		ID:        uuid.New(),
		Name:      name,
		Cards:     []Card{},
		CurrentID: 0,
	}
}

func (d *Deck) AddCard(card Card) {
	d.Cards = append(d.Cards, card)
}

func (d *Deck) RemoveCard(index int) {
	if index < 0 || index >= len(d.Cards) {
		return
	}

	d.Cards = append(d.Cards[:index], d.Cards[index+1:]...)

	if len(d.Cards) == 0 {
		// No cards left, reset CurrentID
		d.CurrentID = 0
	} else if d.CurrentID >= len(d.Cards) {
		d.CurrentID = len(d.Cards) - 1
	}
}

func (d *Deck) NextCard() error {
	if len(d.Cards) == 0 {
		return nil
	}

	d.CurrentID = (d.CurrentID + 1) % len(d.Cards)
	return SaveDeck(*d)
}

func (d *Deck) PrevCard() error {
	if len(d.Cards) == 0 {
		return nil
	}

	d.CurrentID = (d.CurrentID - 1 + len(d.Cards)) % len(d.Cards)
	return SaveDeck(*d)

}

func (d *Deck) CurrentCard() *Card {

	if d.CurrentID < 0 || d.CurrentID >= len(d.Cards) {
		d.CurrentID = 0
	}

	if len(d.Cards) == 0 {
		return nil
	}

	return &d.Cards[d.CurrentID]
}

// ShuffleCards randomizes the order of cards in the deck
func (d *Deck) ShuffleCards() {
	if len(d.Cards) <= 1 {
		return
	}
	
	// Use current time as seed for randomization
	rand.Seed(time.Now().UnixNano())
	
	// Fisher-Yates shuffle algorithm
	for i := len(d.Cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}
	
	// Reset current position to start
	d.CurrentID = 0
}
