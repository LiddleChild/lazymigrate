package migrationview

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
	"github.com/LiddleChild/lazymigrate/internal/components/list"
	"github.com/LiddleChild/lazymigrate/internal/migrator"
)

var _ list.Item = (*item)(nil)

type item struct {
	step           migrator.MigrationStep
	currentVersion uint
	isDirty        bool
	isApplied      bool
}

func (i item) Render(ctx list.Context) string {
	var (
		isDirtyVersion = i.step.Version == i.currentVersion && i.isDirty
		isMigrated     = i.step.Version <= i.currentVersion
	)

	cursor := " "
	if i.step.Version == i.currentVersion {
		if isDirtyVersion {
			cursor = "✗"
		} else {
			cursor = "▶"
		}
	}

	symbol := " "
	if i.step.Version == 0 {
		symbol = "○"
	} else if i.isApplied {
		symbol = "✔"
	} else {
		symbol = "○"
	}

	line := fmt.Sprintf("%s %s %d %s",
		cursor,
		symbol,
		i.step.Version,
		i.step.Identifier,
	)

	style := lipgloss.NewStyle().
		BorderLeft(false).
		BorderStyle(lipgloss.OuterHalfBlockBorder())

	if isDirtyVersion {
		style = style.
			Bold(false).
			Background(brownsugar.ColorRed).
			Foreground(brownsugar.ColorBlack)
	}

	if !isMigrated {
		style = style.
			Bold(false).
			Foreground(brownsugar.ColorBrightBlack)
	}

	if ctx.Selected {
		style = style.
			Bold(true).
			Foreground(brownsugar.ColorCyan).
			Background(brownsugar.ColorBrightBlack).
			BorderLeft(true)

		if isDirtyVersion {
			style = style.
				Bold(true).
				Background(brownsugar.ColorRed).
				Foreground(brownsugar.ColorWhite)
		}
	}

	style = style.
		BorderForeground(style.GetForeground()).
		BorderBackground(style.GetBackground()).
		PaddingLeft(1 - style.GetBorderLeftSize())

	lineWidth := len(line) + style.GetBorderLeftSize() + style.GetPaddingLeft()

	return style.Width(max(ctx.Width, lineWidth)).Render(line)
}
