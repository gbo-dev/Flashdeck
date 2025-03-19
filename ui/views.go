// ui/views.go
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// View renders the UI
func (m model) View() string {
	var content string

	switch m.mode {
	case ModeSettings:
		content = m.ViewSettings()
	case ModeDeckList:
		content = m.ViewDeckList()
	case ModeViewCard, ModeConfirmRemoveCard:
		content = m.ViewCard()
	case ModeCreateDeck:
		content = m.ViewCreateDeck()
	case ModeCreateCard:
		content = m.ViewCreateCard()
	case ModeConfirmDelete:
		content = m.ViewConfirmDelete()
	}

	return AppStyle.Render(content)
}

func (m model) ViewSettings() string {

	title := TitleStyle.
		MarginLeft(2).
		Render("Settings")

	// Set up container styles for settings items with left margin but no top padding
	containerStyle := SettingsContainer
	var settingsContent strings.Builder

	items := []string{
		fmt.Sprintf("Chaos Mode: %s", formatBoolSetting(m.settings.ChaosMode)),
		fmt.Sprintf("Show Timer: %s", formatBoolSetting(m.settings.ShowTimer)),
		fmt.Sprintf("Sound Effects: %s", formatBoolSetting(m.settings.Audio)),
	}

	numSettings := len(items)
	for i, item := range items {
		if m.list.Cursor() == i {
			settingsContent.WriteString(SelectedSettingStyle.
				Render("➤ " + item))
		} else {
			settingsContent.WriteString(SettingItemStyle.
				Render("  " + item))
		}
		if i < numSettings-1 {
			settingsContent.WriteString("\n")
		}
	}

	settingsBox := containerStyle.Render(settingsContent.String())

	leftMargin := lipgloss.NewStyle().MarginLeft(2)

	helpView := m.getHelpView()
	helpWithMargin := leftMargin.Render(helpView)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		settingsBox,
		helpWithMargin,
	)

	return content
}

func (m model) ViewDeckList() string {

		if m.deckManager.GetNumDecks() == 0 {

			title := TitleStyle.MarginLeft(2).Render("Terminal Flashcards")
			message := AppStyle.Render("No decks available.\n\nCreate your first deck by pressing 'n'.")

			helpContent := m.getHelpView()

			return lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				message,
				helpContent,
			)

		} else {
			helpHeight := lipgloss.Height(m.getHelpView())
			listHeight := m.height - helpHeight - 4 // Account for padding and margins

			m.list.SetHeight(listHeight)

			listContent := m.list.View()
			helpContent := m.getHelpView()

			return lipgloss.JoinVertical(
				lipgloss.Left,
				listContent,
				helpContent,
			)
		}
}

func (m model) ViewCard() string {
	helpContent := m.getHelpView()

	if m.currentDeck == nil || len(m.currentDeck.Cards) == 0 {
		// Custom view for empty decks
		counterView := CardCounterView(0, 0)
		title := TitleStyle.Render(m.currentDeck.Name)
		emptyMessage := CardStyle.Render("This deck has no cards yet.")
		instructions := Instructions.Render("Press 'c' to create your first card")

		content := lipgloss.JoinVertical(
			lipgloss.Center,
			counterView,
			title,
			emptyMessage,
			instructions,
			helpContent,
		)
		return content
	}

	card := m.currentDeck.CurrentCard()
	if card == nil {
		return "No cards in this deck."
	}

	counterView := CardCounterView(m.currentDeck.CurrentID+1, len(m.currentDeck.Cards))
	deckTitle := TitleStyle.Render(m.currentDeck.Name)

	if m.mode == ModeConfirmRemoveCard {
		// Create a warning box with thick red borders
		warningStyle := WarningStyleContainer.
			Padding(1, 2).
			MarginLeft(2).
			MarginTop(1).
			Width(60)

		question := card.Question
		if len(question) > 30 {
			question = question[:27] + "..."
		}

		buttons := lipgloss.JoinHorizontal(lipgloss.Center, YesButton, "  ", NoButton)

		confirmBox := warningStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				RedMessageStyle.Render("Are you sure you want to delete this card?\n"),
				question,
				RedMessageStyle.Render("This action cannot be undone.\n"),
				buttons,
			),
		)

		helpContent = confirmBox

	}

	cardContent := CardContentView(card.Question, card.Answer, card.Tags, m.showAnswer)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		counterView,
		deckTitle,
		cardContent,
		helpContent,
	)
	return content
}

func (m model) ViewCreateDeck() string {

	title := TitleStyle.MarginLeft(2).Render("Deck creation")
	helpView := m.getHelpView()

	leftMargin := lipgloss.NewStyle().MarginLeft(2)
	prompt := leftMargin.Render("Enter deck name:")
	inputField := leftMargin.Render(m.newDeckInput.View())

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"\n",
		prompt,
		inputField,
		"\n",
		helpView,
	)

	return content
}

func (m model) ViewCreateCard() string {
	title := TitleStyle.MarginLeft(2).Render("Add card to " + m.currentDeck.Name)
	helpView := m.getHelpView()

	leftMargin := lipgloss.NewStyle().PaddingLeft(2)

	questionLabel := leftMargin.Render("Question:")
	questionInput := leftMargin.Render(m.questionInput.View())

	answerLabel := leftMargin.Render("Answer:")
	answerInput := leftMargin.Render(m.answerInput.View())

	tagsLabel := leftMargin.Render("Tags (comma-separated):")
	tagsInput := leftMargin.Render(m.tagsInput.View())

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"\n",
		questionLabel,
		questionInput,
		"\n",
		answerLabel,
		answerInput,
		"\n",
		tagsLabel,
		tagsInput,
		"\n",
		helpView,
	)

	return content
}

func (m model) ViewConfirmDelete() string {
	title := TitleStyle.MarginLeft(2).Render("Confirm Deletion")

	var confirmBox string
	if m.deckToDelete != nil {
		confirmBox = WarningStyleContainer.Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				RedMessageStyle.Render("Are you sure you want to delete this deck?\n"),
				m.deckToDelete.Name,
				RedMessageStyle.Render("This action cannot be undone.\n"),
			),
		)
	} else {
		confirmBox = WarningStyleContainer.Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				RedMessageStyle.Render("Are you sure you want to delete this deck?"),
				RedMessageStyle.Render("This action cannot be undone.\n"),
			),
		)
	}

	prompt := lipgloss.NewStyle().Render("Type 'delete' to confirm:")
	inputField := InputWarningStyle.Render(m.confirmInput.View())

	helpView := m.getHelpView()

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		confirmBox,
		prompt,
		inputField,
		"\n",
		helpView,
	)

	return content
}

// CustomHelpView creates a neatly organized help view with keybindings
func CustomHelpView(keys []key.Binding) string {

	if len(keys) == 0 {
		return ""
	}

	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))

	regularKeys := keys
	if len(keys) > 1 {
		regularKeys = keys[:len(keys)-2] // Exclude last 2 keys (Back and Quit)
	}

	// Format key bindings in rows of 5
	var rows []string
	for i := 0; i < len(regularKeys); i += 5 {
		var row []string
		end := i + 5
		if end > len(regularKeys) {
			end = len(regularKeys)
		}

		for j := i; j < end; j++ {
			if regularKeys[j].Help().Key == "" {
				continue
			}

			key := keyStyle.Render(regularKeys[j].Help().Key)
			desc := descStyle.Render(regularKeys[j].Help().Desc)

			// Ensure consistent spacing
			keyDesc := fmt.Sprintf("%-10s %-15s", key, desc)
			row = append(row, keyDesc)
		}

		if len(row) > 0 {
			rows = append(rows, strings.Join(row, "   "))
		}
	}

	// Add Back and Quit on their own row if they exist
	if len(keys) > 1 {
		var systemRow []string

		// Get last two keys (Back and Quit), but check them in case the order changed
		lastKeys := keys[len(keys)-2:]

		for _, k := range lastKeys {
			if k.Help().Key == "" {
				continue
			}

			key := keyStyle.Render(k.Help().Key)
			desc := descStyle.Render(k.Help().Desc)

			// Ensure consistent spacing
			keyDesc := fmt.Sprintf("%-10s %-15s", key, desc)
			systemRow = append(systemRow, keyDesc)
		}

		if len(systemRow) > 0 {
			rows = append(rows, strings.Join(systemRow, "   "))
		}
	}

	content := strings.Join(rows, "\n")

	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		BorderForeground(lipgloss.Color("#555555")).
		Render(content)
}

func (m model) getHelpView() string {
	return CustomHelpView(m.getKeysForMode())
}
