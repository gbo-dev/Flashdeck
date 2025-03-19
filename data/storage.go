// data/storage.go
package data

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

const (
	DataDir   = "data"
	DecksDir  = "decks"
	IndexFile = "decks_index.yaml"
)

const SettingsFile = "settings.yaml"

type Settings struct {
	ChaosMode bool `yaml:"chaos_mode"`
	ShowTimer bool `yaml:"show_timer"`
	Audio     bool `yaml:"audio"`
}

func DefaultSettings() Settings {
	return Settings{
		ChaosMode: false,
		ShowTimer: false,
		Audio:     false,
	}
}

func SaveSettings(settings Settings) error {
	data, err := yaml.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(SettingsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	return nil
}

func LoadSettings() (Settings, error) {
	data, err := os.ReadFile(SettingsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultSettings(), nil
		}
		return Settings{}, fmt.Errorf("failed to read settings file: %w", err)
	}

	var settings Settings
	if err := yaml.Unmarshal(data, &settings); err != nil {
		return Settings{}, fmt.Errorf("failed to parse settings file: %w", err)
	}

	return settings, nil
}

func EnsureDirectories() error {
	if err := os.MkdirAll(DecksDir, 0755); err != nil {
		return fmt.Errorf("failed to create decks directory: %w", err)
	}
	return nil
}

func SaveDeck(deck Deck) error {
	yamlData, err := yaml.Marshal(&deck)
	if err != nil {
		return fmt.Errorf("failed to marshal deck: %w", err)
	}

	filename := filepath.Join(DecksDir, fmt.Sprintf("%s.yaml", deck.ID))
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, _ := os.Create(filename)
		file.Close()
	}
	return os.WriteFile(filename, yamlData, 0644)
}

func LoadDeck(id uuid.UUID) (Deck, error) {
	var deck Deck

	filename := filepath.Join(DecksDir, fmt.Sprintf("%s.yaml", id))
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return deck, fmt.Errorf("failed to read deck file: %w", err)
	}

	err = yaml.Unmarshal(yamlFile, &deck)
	if err != nil {
		return deck, fmt.Errorf("failed to unmarshal deck: %w", err)
	}

	return deck, nil
}

func LoadAllDecks() (map[uuid.UUID]Deck, error) {
	decks := make(map[uuid.UUID]Deck)

	files, err := os.ReadDir(DecksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		filePath := filepath.Join(DecksDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %v", file.Name(), err)
		}

		var deck Deck
		err = yaml.Unmarshal(data, &deck)
		if err != nil {
			return nil, fmt.Errorf("failed to parse YAML in file %s: %v", file.Name(), err)
		}

		decks[deck.ID] = deck
	}

	return decks, nil
}

func DeleteDeckFromStorage(id uuid.UUID) error {

	deckFilePath := filepath.Join(DecksDir, fmt.Sprintf("%s.yaml", id))

	// Check if the file exists before delete
	if _, err := os.Stat(deckFilePath); os.IsNotExist(err) {
		return fmt.Errorf("deck file not found: %s", deckFilePath)
	}

	if err := os.Remove(deckFilePath); err != nil {
		return fmt.Errorf("failed to delete deck file: %w", err)
	}

	return nil
}
