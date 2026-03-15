package sourcesscene

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
	"github.com/LiddleChild/lazymigrate/internal/components/list"
	"github.com/LiddleChild/lazymigrate/internal/source"
)

var _ list.Item = (*item)(nil)

type item struct {
	source.Source
	current bool
}

func (i item) Render(ctx list.Context) string {
	var (
		style = lipgloss.NewStyle().
			BorderLeft(ctx.Selected).
			BorderStyle(lipgloss.OuterHalfBlockBorder()).
			BorderForeground(brownsugar.ColorCyan)

		name = i.Name
	)

	style = style.
		PaddingLeft(2 - style.GetBorderLeftSize())

	if ctx.Selected {
		style = style.
			Background(brownsugar.ColorBrightBlack).
			BorderBackground(brownsugar.ColorBrightBlack)
	}

	if i.current {
		style = style.BorderForeground(brownsugar.ColorYellow)
		name = fmt.Sprintf("%s (connected)", name)
	}

	return lipgloss.JoinVertical(lipgloss.Top,
		i.renderLine(ctx, style.Bold(true), name),
		i.renderLine(ctx, style.Foreground(brownsugar.ColorBrightWhite), i.Path),
		i.renderLine(ctx, style.Foreground(brownsugar.ColorBrightWhite), i.DatabaseURL.Redacted()),
	)
}

func (i item) Height() int {
	return 3
}

func (i item) renderLine(ctx list.Context, style lipgloss.Style, line string) string {
	if ctx.Selected {
		style = style.Foreground(brownsugar.ColorCyan)
	}

	if i.current {
		style = style.Foreground(brownsugar.ColorYellow)
	}

	return style.
		Width(max(ctx.Width, lipgloss.Width(line)) + style.GetHorizontalFrameSize()).
		Render(line)
}
