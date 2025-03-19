// ui/keys.go
package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

// keyMap defines the key bindings for the application.
type keyMap struct {
	Up         key.Binding
	Down       key.Binding
	Enter      key.Binding
	Flip       key.Binding
	Next       key.Binding
	Prev       key.Binding
	Back       key.Binding
	Quit       key.Binding
	Help       key.Binding
	Settings   key.Binding
	Toggle     key.Binding
	CreateDeck key.Binding
	CreateCard key.Binding
	DeleteDeck key.Binding
	DeleteCard key.Binding
	Yes        key.Binding
	No         key.Binding
}

// Main menu keymap
var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Flip: key.NewBinding(
		key.WithKeys(" ", "x"),
		key.WithHelp("␣/x", "flip card"),
	),
	Next: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "right"),
	),
	Prev: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "left"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Settings: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "settings"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("␣/enter", "toggle"),
	),
	CreateDeck: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new deck"),
	),
	CreateCard: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "new card"),
	),
	DeleteDeck: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "delete deck"),
	),
	DeleteCard: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete card"),
	),
	Yes: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "yes"),
	),
	No: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "no"),
	),
}

func (m model) getKeysForMode() []key.Binding {
	var keys []key.Binding

	switch m.mode {
	case ModeDeckList:
		if m.deckManager.GetNumDecks() != 0 {
			keys = []key.Binding{
				m.keys.Up,
				m.keys.Down,
				m.keys.Enter,
				m.keys.CreateDeck,
			}
		} else {
			keys = []key.Binding{
				m.keys.CreateDeck,
			}
		}
	case ModeViewCard:
		if m.currentDeck == nil || len(m.currentDeck.Cards) == 0 {
			// For empty decks
			keys = []key.Binding{
				m.keys.CreateCard,
			}
		} else {
			// For decks with cards
			keys = []key.Binding{
				m.keys.Flip,
				m.keys.Next,
				m.keys.Prev,
				m.keys.CreateCard,
				m.keys.DeleteCard,
			}
		}
	case ModeConfirmRemoveCard:
		keys = []key.Binding{
			m.keys.Yes,
			m.keys.No,
		}
	case ModeCreateDeck:
		confirmEnter := key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		)

		keys = []key.Binding{
			confirmEnter,
		}
	case ModeCreateCard:
		confirmEnter := key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		)

		keys = []key.Binding{
			confirmEnter,
		}
	case ModeSettings:
		keys = []key.Binding{
			m.keys.Up,
			m.keys.Down,
			m.keys.Toggle,
		}
	case ModeConfirmDelete:
		confirmEnter := key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		)

		keys = []key.Binding{
			confirmEnter,
		}
	}

	// Settings on deck view, back on rest
	if m.mode == ModeDeckList {
		keys = append(keys, m.keys.Settings)
	} else {
		keys = append(keys, m.keys.Back)
	}

	keys = append(keys, m.keys.Quit)

	return keys
}
