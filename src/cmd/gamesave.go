package cmd

import (
	"github.com/adrg/xdg"
	"os"
	"path/filepath"
)

var (
	GoPuzzlesDataDir     = filepath.Join(xdg.DataHome, "go-puzzles")
	GoPuzzlesGameSaveDir = filepath.Join(GoPuzzlesDataDir, "saves")
)

func saveGame(gameName string, data []byte) error {
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(GoPuzzlesGameSaveDir, 0755); err != nil {
		return err
	}

	// Save the game data to a file
	filePath := filepath.Join(GoPuzzlesGameSaveDir, gameName)
	return os.WriteFile(filePath, data, 0644)
}

func loadGame(gameName string) ([]byte, error) {
	filePath := filepath.Join(GoPuzzlesGameSaveDir, gameName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, nil
	}

	return os.ReadFile(filePath)
}
