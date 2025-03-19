// ui/model.go
package ui

import (
	"fmt"
	"log"
	"strings"

	"go-flashcards/data"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

// Mode represents the current UI mode (e.g., deck list, view card, etc.)
type Mode int

const (
	ModeDeckList Mode = iota
	ModeViewCard
	ModeSettings
	ModeCreateDeck
	ModeCreateCard
	ModeConfirmDelete
	ModeConfirmRemoveCard
)

// model represents the UI state and data
type model struct {
	currentDeck *data.Deck
	deckManager *data.DeckManager
	showAnswer  bool
	mode        Mode
	viewport    viewport.Model
	list        list.Model
	help        help.Model
	keys        keyMap
	width       int
	height      int
	settings    data.Settings

	// Deck creation
	newDeckInput textinput.Model

	// Card creation
	questionInput textinput.Model
	answerInput   textinput.Model
	tagsInput     textinput.Model
	activeInput   int

	confirmInput textinput.Model
	deckToDelete *data.Deck
}

type deckItem struct {
	id    uuid.UUID
	name  string
	count int
}

func (i deckItem) Title() string       { return i.name }
func (i deckItem) ID() uuid.UUID       { return i.id }
func (i deckItem) Description() string { return fmt.Sprintf("%d cards", i.count) }
func (i deckItem) FilterValue() string { return i.name }

func newTextInput(placeholder string, charLimit int, width int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = charLimit
	ti.Width = width
	ti.Cursor.Style = CursorStyle
	ti.PromptStyle = PromptStyle

	return ti
}

func CreateDeckItems(decks []*data.Deck) []list.Item {
	items := make([]list.Item, len(decks))
	for i, deck := range decks {
		items[i] = deckItem{
			id:    deck.ID,
			name:  deck.Name,
			count: len(deck.Cards),
		}
	}
	return items
}

func UpdateDeckList(deckManager *data.DeckManager, listModel *list.Model) {
	decks := deckManager.GetAllDecks()
	// Update the title to include deck count
	deckManager.SortDecksAlphabetical(decks)
	items := CreateDeckItems(decks)
	listModel.SetItems(items)
}

func NewModel(deckManager *data.DeckManager) model {

	decks := deckManager.GetAllDecks()
	deckManager.SortDecksAlphabetical(decks)
	items := CreateDeckItems(decks)

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = ListSelectedItem
	delegate.Styles.SelectedDesc = ListSelectedDesc

	l := list.New(items, delegate, 0, 0)
	l.Title = "Flashdeck"
	l.SetShowHelp(false)
	l.Styles.Title = TitleStyle
	l.Styles.PaginationStyle = ListItemStyle
	l.Styles.HelpStyle = HelpStyle
	l.Styles.NoItems = HelpStyle

	settings, err := data.LoadSettings()
	if err != nil {
		settings = data.DefaultSettings()
	}

	vp := viewport.New(80, 20)

	h := help.New()
	h.Width = 50
	h.ShowAll = true

	newDeckInput := newTextInput("awesome deck name", 50, 30)
	newDeckInput.Focus()
	questionInput := newTextInput("question", 200, 50)
	answerInput := newTextInput("answer", 200, 50)
	tagsInput := newTextInput("tags (comma-separated)", 100, 50)
	confirmInput := newTextInput("Type 'delete' to confirm", 10, 30)

	m := model{
		mode:          ModeDeckList,
		deckManager:   deckManager,
		list:          l,
		viewport:      vp,
		help:          h,
		keys:          keys,
		width:         80,
		height:        24,
		settings:      settings,
		newDeckInput:  newDeckInput,
		questionInput: questionInput,
		answerInput:   answerInput,
		tagsInput:     tagsInput,
		activeInput:   0,
		confirmInput:  confirmInput,
	}
	return m
}

// Init initializes the model.
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global keybindings
		switch {
		case key.Matches(msg, m.keys.Settings) && m.mode == ModeDeckList:
			m.list.Select(0)
			m.mode = ModeSettings
		case key.Matches(msg, m.keys.Quit) && (m.mode != ModeCreateDeck) && (m.mode != ModeCreateCard):
			return m, tea.Quit
		}

		// Mode-specific keybindings
		switch m.mode {
		case ModeSettings:
			switch {
			case key.Matches(msg, m.keys.Up):
				if m.list.Cursor() > 0 {
					m.list.Select(m.list.Cursor() - 1)
				}
			case key.Matches(msg, m.keys.Down):
				if m.list.Cursor() < 2 {
					m.list.Select(m.list.Cursor() + 1)
				}
			case key.Matches(msg, m.keys.Toggle):
				switch m.list.Cursor() {
				case 0:
					m.settings.ChaosMode = !m.settings.ChaosMode
				case 1:
					m.settings.ShowTimer = !m.settings.ShowTimer
				case 2:
					m.settings.Audio = !m.settings.Audio
				}
				if err := data.SaveSettings(m.settings); err != nil {
					log.Printf("Error saving settings: %v", err)
				}
			case key.Matches(msg, m.keys.Back):
				m.mode = ModeDeckList
				m.list.Select(0)
			}

		case ModeDeckList:
			switch {
			case key.Matches(msg, m.keys.Up, m.keys.Down):
				m.list, cmd = m.list.Update(msg)
				cmds = append(cmds, cmd)
			case key.Matches(msg, m.keys.CreateDeck):
				m.mode = ModeCreateDeck
				m.newDeckInput.Focus()
				m.newDeckInput.Reset()
				return m, textinput.Blink
			case key.Matches(msg, m.keys.DeleteDeck):
				// Get the selected deck
				i, ok := m.list.SelectedItem().(deckItem)
				if ok {
					m.deckToDelete = m.deckManager.GetDeckByID(i.id)
					if m.deckToDelete != nil {
						m.mode = ModeConfirmDelete
						m.confirmInput.Reset()
						m.confirmInput.Focus()
						cmds = append(cmds, textinput.Blink)
					}
				}
			case key.Matches(msg, m.keys.Enter):
				i, ok := m.list.SelectedItem().(deckItem)
				if ok {
					m.currentDeck = m.deckManager.GetDeckByID(i.id)
					if m.currentDeck != nil {
						m.mode = ModeViewCard
						m.showAnswer = false
					}
				}
			}

		case ModeViewCard:
			if m.currentDeck == nil || len(m.currentDeck.Cards) == 0 {
				switch {
				case key.Matches(msg, m.keys.Back):
					m.mode = ModeDeckList
				case key.Matches(msg, m.keys.CreateCard):
					m.mode = ModeCreateCard
					m.questionInput.Reset()
					m.answerInput.Reset()
					m.tagsInput.Reset()
					m.questionInput.Focus()
					m.activeInput = 0
					return m, textinput.Blink
				}
				break
			}

			switch {
			case key.Matches(msg, m.keys.Flip):
				m.showAnswer = !m.showAnswer
			case key.Matches(msg, m.keys.Next):
				if err := m.currentDeck.NextCard(); err != nil {
					log.Printf("Error selecting next card: %v", err)
				}
				m.showAnswer = false
			case key.Matches(msg, m.keys.Prev):
				if err := m.currentDeck.PrevCard(); err != nil {
					log.Printf("Error selecting previous card: %v", err)
				}
				m.showAnswer = false
			case key.Matches(msg, m.keys.CreateCard):
				// Switch to card creation mode
				m.mode = ModeCreateCard
				m.questionInput.Reset()
				m.answerInput.Reset()
				m.tagsInput.Reset()
				m.questionInput.Focus()
				m.activeInput = 0
				return m, textinput.Blink
			case key.Matches(msg, m.keys.DeleteCard):
				m.mode = ModeConfirmRemoveCard
			case key.Matches(msg, m.keys.Back):
				m.mode = ModeDeckList
			}

		case ModeCreateDeck:
			switch {
			case key.Matches(msg, m.keys.Enter):
				deckName := strings.TrimSpace(m.newDeckInput.Value())

				if deckName != "" {
					newDeck := data.NewDeck(deckName)

					if err := m.deckManager.AddDeck(newDeck); err != nil {
						log.Printf("Error saving deck: %v", err)
					}

					UpdateDeckList(m.deckManager, &m.list)
					data.SaveDeck(*newDeck)

					m.currentDeck = newDeck

					m.mode = ModeCreateCard
					m.newDeckInput.Reset()
					m.questionInput.Reset()
					m.answerInput.Reset()
					m.tagsInput.Reset()
					m.questionInput.Focus()
					m.activeInput = 0
					return m, textinput.Blink
				}
			case key.Matches(msg, m.keys.Back):
				m.mode = ModeDeckList
			default:
				m.newDeckInput, cmd = m.newDeckInput.Update(msg)
				cmds = append(cmds, cmd)
			}
		case ModeCreateCard:
			switch {
			case key.Matches(msg, m.keys.Back):
				// Cancel card creation and return to viewing the deck
				m.mode = ModeViewCard

			case key.Matches(msg, m.keys.Enter):
				switch m.activeInput {
				case 0: // Question field
					question := strings.TrimSpace(m.questionInput.Value())
					if question != "" {
						// Move to answer field
						m.activeInput = 1
						m.questionInput.Blur()
						m.answerInput.Focus()
						return m, textinput.Blink
					}
					// If empty, just blink
					return m, textinput.Blink

				case 1: // Answer field
					answer := strings.TrimSpace(m.answerInput.Value())
					if answer != "" {
						// Move to tags field
						m.activeInput = 2
						m.answerInput.Blur()
						m.tagsInput.Focus()
						return m, textinput.Blink
					}
					// If empty, just blink
					return m, textinput.Blink

				case 2: // Tags field
					question := strings.TrimSpace(m.questionInput.Value())
					answer := strings.TrimSpace(m.answerInput.Value())

					if question == "" || answer == "" {
						return m, textinput.Blink
					}

					tagsStr := strings.TrimSpace(m.tagsInput.Value())

					// Parse optional tags
					var tags []string
					if tagsStr != "" {
						tags = strings.Split(tagsStr, ",")
						for i, tag := range tags {
							tags[i] = strings.TrimSpace(tag)
						}
					}

					newCard := data.Card{
						Question: question,
						Answer:   answer,
						Tags:     tags,
					}

					if err := m.deckManager.AddCardToDeck(m.currentDeck.ID, newCard); err != nil {
						log.Printf("Error adding card: %v", err)
					}

					UpdateDeckList(m.deckManager, &m.list)

					m.questionInput.Reset()
					m.answerInput.Reset()
					m.tagsInput.Reset()

					m.activeInput = 0
					m.questionInput.Focus()
					m.answerInput.Blur()
					m.tagsInput.Blur()

					return m, textinput.Blink
				}

			case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
				m.activeInput = (m.activeInput + 1) % 3

				m.questionInput.Blur()
				m.answerInput.Blur()
				m.tagsInput.Blur()

				switch m.activeInput {
				case 0:
					m.questionInput.Focus()
				case 1:
					m.answerInput.Focus()
				case 2:
					m.tagsInput.Focus()
				}

				return m, textinput.Blink

			case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"))):
				m.activeInput = (m.activeInput - 1 + 3) % 3

				m.questionInput.Blur()
				m.answerInput.Blur()
				m.tagsInput.Blur()

				switch m.activeInput {
				case 0:
					m.questionInput.Focus()
				case 1:
					m.answerInput.Focus()
				case 2:
					m.tagsInput.Focus()
				}

				return m, textinput.Blink

			default:
				// Handle input updates based on which field is active

				switch m.activeInput {
				case 0:
					m.questionInput, cmd = m.questionInput.Update(msg)
				case 1:
					m.answerInput, cmd = m.answerInput.Update(msg)
				case 2:
					m.tagsInput, cmd = m.tagsInput.Update(msg)
				}

				cmds = append(cmds, cmd)
			}

		case ModeConfirmDelete:
			switch {
			case key.Matches(msg, m.keys.Back):
				m.mode = ModeDeckList
			case key.Matches(msg, m.keys.Enter):
				if strings.TrimSpace(m.confirmInput.Value()) == "delete" {
					if m.deckToDelete != nil {
						if err := m.deckManager.RemoveDeck(m.deckToDelete.ID); err != nil {
							log.Printf("Error deleting deck: %v", err)
						}

						UpdateDeckList(m.deckManager, &m.list)

						if m.list.Cursor() > 0 {
							m.list.Select(m.list.Cursor() - 1)
						}
					}
					m.mode = ModeDeckList
				}
			default:
				m.confirmInput, cmd = m.confirmInput.Update(msg)
				cmds = append(cmds, cmd)
			}
		case ModeConfirmRemoveCard:
			switch {
			case key.Matches(msg, m.keys.Yes), key.Matches(msg, key.NewBinding(key.WithKeys("y"))):
				// User confirmed card deletion
				if m.currentDeck != nil && len(m.currentDeck.Cards) > 0 {

					currentIndex := m.currentDeck.CurrentID
					if err := m.deckManager.RemoveCardFromDeck(m.currentDeck.ID, currentIndex); err != nil {
						log.Printf("Error removing card: %v", err)
					}

					UpdateDeckList(m.deckManager, &m.list)

					// If we removed the last card, we might need to adjust the current ID
					if len(m.currentDeck.Cards) == 0 {
						// No cards left, go back to the deck list
						m.mode = ModeViewCard
					} else if currentIndex >= len(m.currentDeck.Cards) {
						// We removed the last card in the deck, adjust the current ID
						m.currentDeck.CurrentID = len(m.currentDeck.Cards) - 1
						m.mode = ModeViewCard
						m.showAnswer = false
					} else {
						// Just removed a card, stay on the same index (which now points to the next card)
						m.mode = ModeViewCard
						m.showAnswer = false
					}
				}

			case key.Matches(msg, m.keys.No), key.Matches(msg, key.NewBinding(key.WithKeys("n"))), key.Matches(msg, m.keys.Back):
				// User cancelled card deletion
				m.mode = ModeViewCard
			}
		}
	case tea.WindowSizeMsg:
		// Handle window resizing globally
		m.width = msg.Width
		m.height = msg.Height

		// Update list dimensions
		h, v := AppStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 8

		m.help.Width = msg.Width - h
	}

	return m, tea.Batch(cmds...)
}

func formatBoolSetting(b bool) string {
	if b {
		return SettingOnStyle
	}
	return SettingOffStyle
}
