// ui/styles.go
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Define all styles used in the application.
var (
	AppStyle = lipgloss.NewStyle().
			Padding(2, 2)

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#1E735E")).
			Bold(true).
			Padding(0, 1).
			Align(lipgloss.Left).
			MarginBottom(1)

	CardStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#1C7A61")).
			Padding(1, 2).
			Width(75).
			Align(lipgloss.Center)

	QuestionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Align(lipgloss.Center).
			Width(65)

	AnswerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6ED5B8")).
			Bold(true).
			Align(lipgloss.Center).
			Width(65)

	TagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555")).
			Italic(false).
			PaddingTop(1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	CounterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Align(lipgloss.Right)

	ListStyle = lipgloss.NewStyle().
			Margin(1, 0)

	ListItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	ListSelectedItem = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(lipgloss.AdaptiveColor{Light: "#81E3D9", Dark: "#FFFFFF"}).
				Foreground(lipgloss.AdaptiveColor{Light: "#2CA889", Dark: "#2CA889"}).
				Padding(0, 0, 0, 1)

	ListSelectedDesc = ListSelectedItem.
				Foreground(lipgloss.AdaptiveColor{Light: "#666666", Dark: "#666666"})

	SettingItemStyle = lipgloss.NewStyle().
				PaddingLeft(0).
				Foreground(lipgloss.Color("#FFFFFF"))

	SelectedSettingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#33C4A0")).
				Bold(true)

	CursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	PromptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2CA889"))

	MessageStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)

	WarningStyleContainer = lipgloss.NewStyle().
				Align(lipgloss.Center).
				BorderStyle(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("#FF0000")).
				Width(60)

	InputWarningStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#FF0000")).
				Padding(0, 1).
				MarginLeft(2).
				Width(30)

	SettingsContainer = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#1C7A61")).
				Padding(0, 0, 0, 1).
				MarginLeft(2).
				Width(30)

	SettingOnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7EBC39")).
			Bold(true).
			Render("ON")

	SettingOffStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true).
			Render("OFF")

	RedMessageStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			MarginTop(1).
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)

	ButtonStyle = lipgloss.NewStyle().
			Padding(0, 3).
			Bold(true)

	YesButton = ButtonStyle.
			Foreground(lipgloss.Color("#FF0000")).
			Render("[y] Yes")

	NoButton = ButtonStyle.
			Foreground(lipgloss.Color("#7EBC39")).
			Render("[n] No")

	Instructions = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Align(lipgloss.Center).
			Width(60)
)

// GetCardCounterView renders the card counter.
func CardCounterView(current, total int) string {
	counter := fmt.Sprintf("Card %d of %d", current, total)
	return CounterStyle.Render(counter)
}

// GetCardContentView renders the flashcard content (question, answer, and tags).
func CardContentView(question, answer string, tags []string, showAnswer bool) string {
	content := QuestionStyle.Render(question)
	if showAnswer {
		content += "\n\n" + AnswerStyle.Render(answer)
	} else {
		content += "\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Render("[Press SPACE to reveal the answer]")
	}

	if len(tags) > 0 {
		tagText := "Tags: " + strings.Join(tags, ", ")
		content += "\n\n" + TagStyle.Render(tagText)
	}

	return CardStyle.Render(content)
}
