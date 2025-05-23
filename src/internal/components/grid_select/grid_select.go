package grid_select

import (
	"github.com/Liuuner/go-puzzles/src/internal/common"
	"github.com/Liuuner/go-puzzles/src/internal/style"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math"
	"strings"
)

func (m Model[T]) Init() tea.Cmd {
	return nil
}

func (m Model[T]) Update(msg tea.Msg) Model[T] {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(msg, common.Hotkeys.Up):
			newSelected := m.selectedItem - m.maxItemsInRow
			if newSelected >= 0 {
				m.selectedItem = newSelected
			}
		case key.Matches(msg, common.Hotkeys.Down):
			newSelected := m.selectedItem + m.maxItemsInRow
			if newSelected < len(m.items) {
				m.selectedItem = newSelected
			}
		case key.Matches(msg, common.Hotkeys.Left):
			if m.selectedItem%3 != 0 {
				m.selectedItem = max(0, m.selectedItem-1)
			}
		case key.Matches(msg, common.Hotkeys.Right):
			if (m.selectedItem+1)%m.maxItemsInRow != 0 {
				m.selectedItem = min(len(m.items)-1, m.selectedItem+1)
			}
		case key.Matches(msg, common.Hotkeys.Select):
			m.selectFunc(m.items[m.selectedItem])
		}
	}

	return m
}

func (m Model[T]) View() string {
	sb := strings.Builder{}
	width := (common.Config.SelectionContainerWidth+2)*m.itemsInRow + m.itemsInRow - 1

	lineAmount := len(m.items) / m.itemsInRow
	if len(m.items)%m.itemsInRow != 0 {
		lineAmount++
	}

	grid := make([][]string, lineAmount)
	for i := range grid {
		grid[i] = make([]string, m.itemsInRow)
	}

	for i, p := range m.items {
		selected := i == m.selectedItem /* && !m.helpOpened*/
		//isInSelectedRow := i/m.itemsInRow == selectedRow

		lineNum := i / m.itemsInRow
		itemNum := i % m.itemsInRow
		if itemNum != 0 && m.gapColumn > 0 {
			grid[lineNum][itemNum] = lipgloss.NewStyle().MarginLeft(m.gapColumn).Render(m.renderFunc(p, selected))
		} else {
			grid[lineNum][itemNum] = m.renderFunc(p, selected)
		}
	}

	// todo add gapRow
	for _, line := range grid {
		if m.gapRow > 0 {
			sb.WriteString(nStrings("\n", m.gapRow))
		}
		sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, line...))
		sb.WriteString("\n")
	}
	selection := lipgloss.NewStyle().Width(width).Render(sb.String())

	return m.cropPuzzleSelection(selection)
}

func (m Model[T]) cropPuzzleSelection(puzzleSelection string) string {
	selectedRow := m.selectedItem / m.itemsInRow // 3
	//rowAmount := int(math.Ceil(float64(len(m.puzzles)) / float64(m.layout.selectionItemsInRow)))

	maxHeight := m.height
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

func nStrings(s string, n int) string {
	sb := strings.Builder{}
	for range n {
		sb.WriteString(s)
	}
	return sb.String()
}
