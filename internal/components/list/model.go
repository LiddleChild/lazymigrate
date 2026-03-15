package list

import (
	"slices"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
	"github.com/LiddleChild/lazymigrate/internal/components/focus"
	"github.com/LiddleChild/lazymigrate/internal/components/scrollpane"
)

type Model struct {
	focus.Model

	items  []Item
	cursor int

	borderForegroundColor lipgloss.ANSIColor

	viewport viewport.Model
}

func New() *Model {
	viewport := viewport.New()
	viewport.KeyMap.Up.SetEnabled(false)
	viewport.KeyMap.Down.SetEnabled(false)

	return &Model{
		Model:                 focus.New(),
		items:                 make([]Item, 0),
		cursor:                0,
		borderForegroundColor: brownsugar.ColorWhite,
		viewport:              viewport,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	return m, nil
}

func (m *Model) Render(ctx brownsugar.Context) string {
	var (
		scrollpane = scrollpane.New().
				SetWidth(ctx.Width).
				SetHeight(ctx.Height).
				Foreground(m.borderForegroundColor).
				BorderStyle(lipgloss.RoundedBorder()).
				CursorStyle(lipgloss.OuterHalfBlockBorder())

		width  = ctx.Width - scrollpane.GetHorizontalBorderSize()
		height = ctx.Height - scrollpane.GetVerticalBorderSize()
	)

	contents := []string{}
	for i, item := range m.items {
		contents = append(contents,
			item.Render(Context{
				Index:    i,
				Width:    ctx.Width,
				Selected: i == m.cursor,
			}),
		)
	}

	m.viewport.SetWidth(width)
	m.viewport.SetHeight(height)
	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Top, contents...))

	return scrollpane.
		SetTotalLine(len(m.items)).
		SetCurrentLine(m.viewport.YOffset()).
		Render(m.viewport.View())
}

func (m *Model) Up() int {
	return m.SetCursor(m.cursor - 1)
}

func (m *Model) Down() int {
	return m.SetCursor(m.cursor + 1)
}

func (m *Model) Top() int {
	return m.SetCursor(0)
}

func (m *Model) Bottom() int {
	return m.SetCursor(len(m.items) - 1)
}

func (m *Model) GetCursor() int {
	return m.cursor
}

func (m *Model) SetCursor(cursor int) int {
	if cursor < 0 {
		cursor = 0
	} else if cursor >= len(m.items) {
		cursor = len(m.items) - 1
	}

	diff := cursor - m.cursor
	m.cursor = cursor

	return diff
}

func (m *Model) FocusAtCursor() {
	var (
		offset = m.cursor - m.viewport.Height()/2

		topBound    = m.cursor - m.viewport.Height()/2
		bottomBound = m.cursor + m.viewport.Height()/2
	)

	if topBound < 0 {
		offset = 0
	} else if bottomBound >= len(m.items) {
		offset = topBound
	}

	m.viewport.SetYOffset(offset)
}

func (m *Model) GetItems() []Item {
	return m.items
}

func (m *Model) SetItems(items []Item) {
	m.items = slices.Clone(items)

	// reset cursor to ensure that cursor is always in bound
	m.SetCursor(m.cursor)
}

func (m *Model) GetSelectedItem() Item {
	return m.items[m.cursor]
}

func (m *Model) SetBorderForegroundColor(color lipgloss.ANSIColor) {
	m.borderForegroundColor = color
}
