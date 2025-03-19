// data/manager.go
package data

import (
	"fmt"
	"sort"
	"strings"

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
