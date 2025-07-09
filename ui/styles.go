// ui/styles.go
package ui

import (
	"fmt"
	"strings"
	"time"

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

	// Code formatting styles
	CodeBlockStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#98FB98")).
			Background(lipgloss.Color("#2F2F2F")).
			Padding(0, 1).
			MarginTop(1).
			MarginBottom(1)

	InlineCodeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB347")).
			Background(lipgloss.Color("#2F2F2F")).
			Padding(0, 1)
)

// GetCardCounterView renders the card counter.
func CardCounterView(current, total int) string {
	counter := fmt.Sprintf("Card %d of %d", current, total)
	return CounterStyle.Render(counter)
}

// GetTimerView renders the timer display.
func GetTimerView(elapsed time.Duration) string {
	seconds := int(elapsed.Seconds())
	minutes := seconds / 60
	seconds = seconds % 60
	
	timerText := fmt.Sprintf("⏱️  %02d:%02d", minutes, seconds)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#88D4E6")).
		Bold(true).
		Align(lipgloss.Center).
		Render(timerText)
}

// formatCodeBlocks processes text to highlight code blocks and inline code
func formatCodeBlocks(text string) string {
	// Handle code blocks (```code```)
	for strings.Contains(text, "```") {
		start := strings.Index(text, "```")
		if start == -1 {
			break
		}
		end := strings.Index(text[start+3:], "```")
		if end == -1 {
			break
		}
		end += start + 3
		
		codeContent := text[start+3 : end]
		formattedCode := CodeBlockStyle.Render(codeContent)
		text = text[:start] + formattedCode + text[end+3:]
	}
	
	// Handle inline code (`code`)
	for strings.Contains(text, "`") {
		start := strings.Index(text, "`")
		if start == -1 {
			break
		}
		end := strings.Index(text[start+1:], "`")
		if end == -1 {
			break
		}
		end += start + 1
		
		codeContent := text[start+1 : end]
		formattedCode := InlineCodeStyle.Render(codeContent)
		text = text[:start] + formattedCode + text[end+1:]
	}
	
	return text
}

// GetCardContentView renders the flashcard content (question, answer, and tags).
func CardContentView(question, answer string, tags []string, showAnswer bool) string {
	// Format code in question
	formattedQuestion := formatCodeBlocks(question)
	content := QuestionStyle.Render(formattedQuestion)
	
	if showAnswer {
		// Format code in answer
		formattedAnswer := formatCodeBlocks(answer)
		content += "\n\n" + AnswerStyle.Render(formattedAnswer)
	} else {
		content += "\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Render("💭 [Press SPACE to reveal the answer]")
	}

	if len(tags) > 0 {
		tagText := "🏷️  Tags: " + strings.Join(tags, ", ")
		content += "\n\n" + TagStyle.Render(tagText)
	}

	return CardStyle.Render(content)
}
