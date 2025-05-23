package common

import tea "github.com/charmbracelet/bubbletea"

// QuitGameMsg signals that the program to close the game. You can send a [QuitGameMsg] with
// [QuitGame].
type QuitGameMsg struct{}

// QuitGame is a special command that tells the program to close the current game.
func QuitGame() tea.Msg {
	return QuitGameMsg{}
}

type SaveGameMsg struct {
	GameName string
	Data     []byte
}

func SaveGame(gameName string, data []byte) tea.Cmd {
	return func() tea.Msg {
		return SaveGameMsg{
			GameName: gameName,
			Data:     data,
		}
	}
}

type LoadGameMsg struct {
	GameName string
}

func LoadGame(gameName string) tea.Cmd {
	return func() tea.Msg {
		return LoadGameMsg{
			GameName: gameName,
		}
	}
}

type LoadGameResponseMsg struct {
	GameName string
	Data     []byte
}

func LoadGameResponse(gameName string, data []byte) tea.Cmd {
	return func() tea.Msg {
		return LoadGameResponseMsg{
			GameName: gameName,
			Data:     data,
		}
	}
}
