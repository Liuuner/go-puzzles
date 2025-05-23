package cmd

import (
	"context"
	"github.com/Liuuner/go-puzzles/src/internal/common"
	"github.com/Liuuner/go-puzzles/src/internal/components"
	"github.com/Liuuner/go-puzzles/src/internal/puzzles"
	"github.com/Liuuner/go-puzzles/src/internal/puzzles/minesweeper"
	"github.com/Liuuner/go-puzzles/src/internal/style"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v3"
	"log"
	"math"
	"os"
	"strings"
)

func Run() {

	cmd := &cli.Command{
		EnableShellCompletion: true,
		Name:                  "go-puzzles",
		Version:               "v0.0.1-dev",
		HideVersion:           true,
		HideHelpCommand:       true,
		Description:           "A collection of terminal puzzles",
		Commands: []*cli.Command{
			minesweeper.Command(runStandalone),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runDefault()
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func runDefault() error {
	allPuzzles := []puzzles.Puzzle{
		minesweeper.Minesweeper{},
		puzzles.EmptyPuzzle{},
		puzzles.EmptyPuzzle{},

		puzzles.EmptyPuzzle{},
		puzzles.EmptyPuzzle{},
		minesweeper.Minesweeper{},

		puzzles.EmptyPuzzle{},
		puzzles.EmptyPuzzle{},
		puzzles.EmptyPuzzle{},

		minesweeper.Minesweeper{},
		puzzles.EmptyPuzzle{},
		minesweeper.Minesweeper{},
		puzzles.EmptyPuzzle{},
		puzzles.EmptyPuzzle{},
		minesweeper.Minesweeper{},
		puzzles.EmptyPuzzle{},
		puzzles.EmptyPuzzle{},
		puzzles.EmptyPuzzle{},
		minesweeper.Minesweeper{},
		puzzles.EmptyPuzzle{},
	}
	/*
		grid_select.New[puzzles.Puzzle]().
			Data(allPuzzles).
			RenderFunc(func(p puzzles.Puzzle, selected bool) string { return p.Preview() }).*/
	m := model{
		puzzle:       puzzles.EmptyPuzzle{},
		puzzleOpened: false,
		//puzzleSelection: puzzleSelection,
		selectedPuzzle: 0,
		/**/
		puzzles: allPuzzles,
		layout: layout{
			selectionItemsInRow: 1,
		},
	}
	/*
		puzzleSelection :=
			grid_select.New[puzzles.Puzzle]().
				Data(allPuzzles).
				RenderFunc(func(p puzzles.Puzzle, selected bool) string { return p.Preview() }).
	*/

	p := tea.NewProgram(m, tea.WithAltScreen() /*, tea.WithMouseCellMotion()*/)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func runStandalone(puzzle puzzles.Puzzle) error {
	m := model{
		puzzle:         puzzle,
		standaloneMode: true,
	}

	p := tea.NewProgram(m, tea.WithAltScreen() /*, tea.WithMouseCellMotion()*/)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.handleWindowResize(msg)
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	if _, ok := msg.(common.QuitGameMsg); ok {
		if m.standaloneMode {
			return m, tea.Quit
		}
		m.puzzleOpened = false
		m.puzzle = puzzles.EmptyPuzzle{}
	}

	// Handling save and load game messages
	{
		if msg, ok := msg.(common.SaveGameMsg); ok {
			err := saveGame(msg.GameName, msg.Data)
			// TODO not really a good way to handle errors
			if err != nil {
				panic(err)
			}
		}
		if msg, ok := msg.(common.LoadGameMsg); ok {
			data, err := loadGame(msg.GameName)
			if err != nil {
				panic(err)
			}

			return m, common.LoadGameResponse(msg.GameName, data)
		}

		if msg, ok := msg.(common.LoadGameResponseMsg); ok {
			if m.puzzle.Name() == msg.GameName {
				m.puzzle.Update(msg.Data)
			}
		}
	}

	var batch tea.Cmd

	if m.puzzleOpened || m.standaloneMode {
		puzzle, cmd := m.puzzle.Update(msg)
		m.puzzle = puzzle
		batch = tea.Batch(batch, cmd)
	} else {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			// todo handleKeyMsg will tell when a puzzle is selected with a cmd like PuzzleSelected
			m, batch = m.handleKeyMsg(keyMsg)
			if m.puzzleOpened {
				// todo when all m.puzzle.Init() commands are finished, we return the PuzzleReady Cmd or something in this fashion
				batch = tea.Batch(batch, m.puzzle.Init(), tea.WindowSize())
			}
		}
	}

	return m, batch
}

func (m model) handleKeyMsg(msg tea.KeyMsg) (model, tea.Cmd) {
	switch {
	case key.Matches(msg, common.Hotkeys.Quit):
		if m.helpOpened {
			m.helpOpened = false
		} else {
			return m, tea.Quit
		}
	case key.Matches(msg, common.Hotkeys.Up):
		newSelected := m.selectedPuzzle - m.layout.selectionItemsInRow
		if newSelected >= 0 {
			m.selectedPuzzle = newSelected
		}
	case key.Matches(msg, common.Hotkeys.Down):
		newSelected := m.selectedPuzzle + m.layout.selectionItemsInRow
		if newSelected < len(m.puzzles) {
			m.selectedPuzzle = newSelected
		}
	case key.Matches(msg, common.Hotkeys.Left):
		if m.selectedPuzzle%3 != 0 {
			m.selectedPuzzle = max(0, m.selectedPuzzle-1)
		}
	case key.Matches(msg, common.Hotkeys.Right):
		if (m.selectedPuzzle+1)%m.layout.selectionItemsInRow != 0 {
			m.selectedPuzzle = min(len(m.puzzles)-1, m.selectedPuzzle+1)
		}
	case key.Matches(msg, common.Hotkeys.Select):
		if !m.puzzleOpened {
			m.puzzle = m.puzzles[m.selectedPuzzle].New()
			m.puzzleOpened = true
		}
	case key.Matches(msg, common.Hotkeys.Help):
		m.helpOpened = !m.helpOpened
	}
	return m, nil
}

func (m model) View() string {
	if m.puzzleOpened || m.standaloneMode {
		return m.puzzle.View()
	}

	header := drawHeader(m)

	puzzleSelection := drawPuzzleSelection(m)
	puzzleSelection = cropPuzzleSelection(m, puzzleSelection)

	help := ""
	if m.helpOpened {
		help = drawHelp()
	}
	body := style.Composite(help, puzzleSelection, lipgloss.Center, lipgloss.Center, 0, 0)

	return lipgloss.JoinVertical(lipgloss.Center, header, body)
}

func drawHelp() string {
	helpContent := lipgloss.NewStyle().Bold(true).Render("Help") + "\n\n" +
		"Movement:\n" +
		hotkeyHelp(common.Hotkeys.Up) + "\n" +
		hotkeyHelp(common.Hotkeys.Down) + "\n" +
		hotkeyHelp(common.Hotkeys.Left) + "\n" +
		hotkeyHelp(common.Hotkeys.Right) + "\n\n" +
		"Other:\n" +
		hotkeyHelp(common.Hotkeys.Select) + "\n" +
		hotkeyHelp(common.Hotkeys.Quit) + "\n" +
		hotkeyHelp(common.Hotkeys.Help) + "\n\n" +
		"Press q to quit or esc to close"

	return lipgloss.NewStyle().
		Width(64).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#89b4fa")).
		PaddingLeft(1).
		Render(helpContent)
}

func hotkeyHelp(b key.Binding) string {
	return b.Help().Key + " " + b.Help().Desc
}

func (m *model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.terminalInfo.fullWidth = msg.Width
	m.terminalInfo.fullHeight = msg.Height

	m.layout.selectionHeight = m.terminalInfo.fullHeight - common.Config.HeaderHeight - 1

	m.layout.selectionItemsInRow = max(m.terminalInfo.fullWidth/(common.Config.SelectionContainerWidth+3), 1)

	if m.puzzleOpened {
		m.puzzle.Update(msg)
	}
}

func drawHeader(m model) string {
	title := lipgloss.NewStyle().Bold(true).Render("Go Puzzles")

	return components.Header(m.terminalInfo.fullWidth, "v0.0.1-dev", title, "Help: ?")
}

func drawPuzzleSelection(m model) string {
	sb := strings.Builder{}
	width := (common.Config.SelectionContainerWidth+2)*m.layout.selectionItemsInRow + m.layout.selectionItemsInRow - 1

	lineAmount := len(m.puzzles) / m.layout.selectionItemsInRow
	if len(m.puzzles)%m.layout.selectionItemsInRow != 0 {
		lineAmount++
	}

	grid := make([][]string, lineAmount)
	for i := range grid {
		grid[i] = make([]string, m.layout.selectionItemsInRow)
	}

	for i, p := range m.puzzles {
		selected := i == m.selectedPuzzle && !m.helpOpened
		//isInSelectedRow := i/m.layout.selectionItemsInRow == selectedRow

		lineNum := i / m.layout.selectionItemsInRow
		itemNum := i % m.layout.selectionItemsInRow
		if itemNum != 0 {
			grid[lineNum][itemNum] = lipgloss.NewStyle().MarginLeft(1).Render(buildPuzzleContainer(p, selected))
		} else {
			grid[lineNum][itemNum] = buildPuzzleContainer(p, selected)
		}
	}

	for _, line := range grid {
		sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, line...))
		sb.WriteString("\n")
	}
	return lipgloss.NewStyle().Width(width).Render(sb.String())
}

func cropPuzzleSelection(m model, puzzleSelection string) string {
	selectedRow := m.selectedPuzzle / m.layout.selectionItemsInRow // 3
	//rowAmount := int(math.Ceil(float64(len(m.puzzles)) / float64(m.layout.selectionItemsInRow)))

	maxHeight := m.layout.selectionHeight
	height := lipgloss.Height(puzzleSelection)

	fullyDisplayedRowsOnScreen := int(math.Floor(float64(maxHeight) / float64(common.Config.SelectionContainerHeight+2))) //3

	offsetTop := 0

	if selectedRow >= fullyDisplayedRowsOnScreen {
		offsetTop = (selectedRow + 1 - fullyDisplayedRowsOnScreen) * (common.Config.SelectionContainerHeight + 2)
	}
	if height-offsetTop < maxHeight {
		offsetTop = offsetTop - (maxHeight - (fullyDisplayedRowsOnScreen * (common.Config.SelectionContainerHeight + 2))) // 2
	}

	if height > maxHeight {
		puzzleSelection = style.TranslateYContainer(offsetTop, maxHeight, puzzleSelection)
	}

	return puzzleSelection
}

func buildPuzzleContainer(puzzle puzzles.Puzzle, selected bool) string {
	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Height(common.Config.SelectionContainerHeight).
		Width(common.Config.SelectionContainerWidth)

	if selected {
		containerStyle = containerStyle.BorderForeground(lipgloss.Color("#89b4fa"))
	}

	return containerStyle.Render(puzzle.Preview())
}
