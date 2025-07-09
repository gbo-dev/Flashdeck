// data/manager.go
package data

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type DeckManager struct {
	decks map[uuid.UUID]*Deck
}

func NewDeckManager() *DeckManager {
	return &DeckManager{
		decks: make(map[uuid.UUID]*Deck),
	}
}

func (dm *DeckManager) LoadAllDecks() error {

	// decks: map[uuid.UUID]Deck
	decks, err := LoadAllDecks()
	if err != nil {
		return err
	}

	dm.decks = make(map[uuid.UUID]*Deck)
	for _, deck := range decks {
		dm.decks[deck.ID] = &deck
	}
	return nil
}

func (dm *DeckManager) GetDeckByID(id uuid.UUID) *Deck {
	return dm.decks[id]
}

func (dm *DeckManager) GetAllDecks() []*Deck {
	decks := make([]*Deck, 0, len(dm.decks))
	for _, deck := range dm.decks {
		decks = append(decks, deck)
	}
	return decks
}

func (dm *DeckManager) GetNumDecks() int {
	return len(dm.decks)
}

// AddDeck adds a deck to DeckManager and in storage
func (dm *DeckManager) AddDeck(deck *Deck) error {
	if deck.ID == uuid.Nil {
		deck.ID = uuid.New()
	}

	if err := SaveDeck(*deck); err != nil {
		return err
	}

	dm.decks[deck.ID] = deck
	return nil
}

func (dm *DeckManager) RemoveDeck(id uuid.UUID) error {
	// Check if deck exists in memory
	if _, exists := dm.decks[id]; !exists {
		return fmt.Errorf("deck not found with ID: %s", id)
	}

	if err := DeleteDeckFromStorage(id); err != nil {
		return err
	}

	// Delete reference from dm
	delete(dm.decks, id)

	return nil
}

func (dm *DeckManager) AddCardToDeck(deckID uuid.UUID, card Card) error {
	deck := dm.GetDeckByID(deckID)
	if deck == nil {
		return fmt.Errorf("deck not found with ID: %s", deckID)
	}

	deck.AddCard(card)

	return SaveDeck(*deck)
}

func (dm *DeckManager) RemoveCardFromDeck(deckID uuid.UUID, cardIndex int) error {
	deck := dm.GetDeckByID(deckID)
	if deck == nil {
		return fmt.Errorf("deck not found with ID: %s", deckID)
	}

	deck.RemoveCard(cardIndex)

	return SaveDeck(*deck)
}

func (dm *DeckManager) SaveDeckState(deckID uuid.UUID) error {
	deck := dm.GetDeckByID(deckID)
	if deck == nil {
		return fmt.Errorf("deck not found with ID: %s", deckID)
	}

	return SaveDeck(*deck)
}

func (dm *DeckManager) SortDecksAlphabetical(decks []*Deck) {
	sort.Slice(decks, func(i, j int) bool {
		return strings.ToLower(decks[i].Name) < strings.ToLower(decks[j].Name)
	})
}

// GetRandomCard returns a random card from any deck (chaos mode)
func (dm *DeckManager) GetRandomCard() *Card {
	allCards := []*Card{}
	
	// Collect all cards from all decks
	for _, deck := range dm.decks {
		for i := range deck.Cards {
			allCards = append(allCards, &deck.Cards[i])
		}
	}
	
	if len(allCards) == 0 {
		return nil
	}
	
	// Return random card
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(allCards))
	return allCards[randomIndex]
}

// GetRandomDeckWithCards returns a random deck that has cards
func (dm *DeckManager) GetRandomDeckWithCards() *Deck {
	decksWithCards := []*Deck{}
	
	// Collect all decks that have cards
	for _, deck := range dm.decks {
		if len(deck.Cards) > 0 {
			decksWithCards = append(decksWithCards, deck)
		}
	}
	
	if len(decksWithCards) == 0 {
		return nil
	}
	
	// Return random deck
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(decksWithCards))
	return decksWithCards[randomIndex]
}
