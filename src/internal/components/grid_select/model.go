package grid_select

import "github.com/charmbracelet/lipgloss"

type Model[T any] struct {
	items         []T
	selectedItem  int
	renderFunc    func(item T, selected bool) string
	selectFunc    func(item T)
	border        lipgloss.Border
	gapRow        int
	gapColumn     int
	itemsInRow    int
	maxItemsInRow int
	width         int
	height        int
}

func New[T any]() Model[T] {
	m := &Model[T]{
		selectedItem: 0,
	}

	return *m
}

func (m Model[T]) Data(items []T) Model[T] {
	m.items = items
	return m
}

func (m Model[T]) RenderFunc(render func(item T, selected bool) string) Model[T] {
	m.renderFunc = render
	return m
}

func (m Model[T]) SelectFunc(selectFn func(item T)) Model[T] {
	m.selectFunc = selectFn
	return m
}

func (m Model[T]) Gap(gap int) Model[T] {
	m.gapColumn = gap
	m.gapRow = gap
	return m
}

func (m Model[T]) GapRow(gap int) Model[T] {
	m.gapRow = gap
	return m
}

func (m Model[T]) GapColumn(gap int) Model[T] {
	m.gapRow = gap
	return m
}

func (m Model[T]) Border(border lipgloss.Border) Model[T] {
	m.border = border
	return m
}

func (m Model[T]) MaxItemsPerRow(max int) Model[T] {
	m.maxItemsInRow = max
	return m
}

func (m Model[T]) Width(width int) Model[T] {
	m.width = width
	return m
}

func (m Model[T]) Height(height int) Model[T] {
	m.height = height
	return m
}
