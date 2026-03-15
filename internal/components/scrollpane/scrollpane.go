package scrollpane

import (
	"math"
	"slices"

	"charm.land/lipgloss/v2"
)

type Scrollpane struct {
	foreground lipgloss.ANSIColor

	borderStyle lipgloss.Border
	cursorStyle lipgloss.Border

	visibleLine int
	totalLine   int
	currentLine int

	width  int
	height int
}

func New() Scrollpane {
	return Scrollpane{}
}

func (p Scrollpane) SetTotalLine(totalLine int) Scrollpane {
	p.totalLine = totalLine
	return p
}

func (p Scrollpane) SetCurrentLine(currentLine int) Scrollpane {
	p.currentLine = currentLine
	return p
}

func (p Scrollpane) Foreground(color lipgloss.ANSIColor) Scrollpane {
	p.foreground = color
	return p
}

func (p Scrollpane) BorderStyle(b lipgloss.Border) Scrollpane {
	p.borderStyle = b
	return p
}

func (p Scrollpane) CursorStyle(b lipgloss.Border) Scrollpane {
	p.cursorStyle = b
	return p
}

func (p Scrollpane) SetWidth(width int) Scrollpane {
	p.width = width
	return p
}

func (p Scrollpane) SetHeight(height int) Scrollpane {
	p.height = height
	p.visibleLine = height - p.GetVerticalBorderSize()
	return p
}

func (p Scrollpane) Render(content string) string {
	scrollbars := make([]string, 0, p.visibleLine)
	if p.totalLine > p.visibleLine {
		var (
			fCurrentLine = float64(p.currentLine)
			fVisibleLine = float64(p.visibleLine) - 1
			fTotalLine   = float64(p.totalLine) - 1
		)

		var (
			scrollbarTop  = math.Round(fVisibleLine / fTotalLine * fCurrentLine)
			scrollbarSize = math.Round(fVisibleLine / fTotalLine * fVisibleLine)
		)

		for i := float64(0); i < float64(p.visibleLine); i++ {
			if i >= scrollbarTop && i <= scrollbarTop+scrollbarSize {
				scrollbars = append(scrollbars, p.cursorStyle.Right)
			} else {
				scrollbars = append(scrollbars, p.borderStyle.Right)
			}
		}
	} else {
		for i := 0; i < p.visibleLine; i++ {
			scrollbars = append(scrollbars, p.borderStyle.Right)
		}
	}

	rightBorders := slices.Concat(
		[]string{p.borderStyle.TopRight},
		scrollbars,
		[]string{p.borderStyle.BottomRight},
	)

	scrollbar := lipgloss.NewStyle().
		Foreground(p.foreground).
		Render(lipgloss.JoinVertical(lipgloss.Top, rightBorders...))

	return lipgloss.JoinHorizontal(lipgloss.Left,
		lipgloss.NewStyle().
			Width(p.width-1).
			Height(p.height).
			Border(p.borderStyle).
			BorderRight(false).
			BorderForeground(p.foreground).
			Render(content),
		scrollbar,
	)
}

// TODO: properly calculate this
func (p Scrollpane) GetHorizontalBorderSize() int {
	return 2
}

// TODO: properly calculate this
func (p Scrollpane) GetVerticalBorderSize() int {
	return 2
}
