package components

/*
import tea "github.com/charmbracelet/bubbletea"

type Window struct {
	Header     tea.Model
	Body       tea.Model
	Footer     tea.Model
	Help       tea.Model
	helpOpened bool
}

func NewWindow() Window {
	return Window{}
}

func (w Window) Update(msg tea.Msg) (Window, tea.Cmd) {
	var cmd, cmd1, cmd2, cmd3 tea.Cmd

	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		w.Body.Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - w.Header.Height() - w.Footer.Height(),
		})
	}

	w.Header, cmd1 = w.Header.Update(msg)

	if w.helpOpened {
		w.Help, cmd2 = w.Help.Update(msg)
	} else {
		w.Body, cmd2 = w.Body.Update(msg)
	}

	w.Footer, cmd3 = w.Footer.Update(msg)

	return w, tea.Batch(cmd, cmd1, cmd2, cmd3)
}

func (w Window) View() string {

}
*/
