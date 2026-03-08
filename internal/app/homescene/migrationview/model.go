package migrationview

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LiddleChild/lazymigrate/internal/appevent"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
	"github.com/LiddleChild/lazymigrate/internal/components/focus"
	"github.com/LiddleChild/lazymigrate/internal/migrator"
)

var (
	Keyj     = key.NewBinding(key.WithKeys("j", "down"))
	Keyk     = key.NewBinding(key.WithKeys("k", "up"))
	Keyg     = key.NewBinding(key.WithKeys("g"))
	KeyG     = key.NewBinding(key.WithKeys("G"))
	KeySpace = key.NewBinding(key.WithKeys("space"))
	Keyf     = key.NewBinding(key.WithKeys("f"))
)

type Model struct {
	focus.Model

	migration migrator.Migration
	cursor    int
	isLocked  *atomic.Bool

	viewport viewport.Model
}

func New() *Model {
	viewport := viewport.New()
	viewport.KeyMap.Up.SetEnabled(false)
	viewport.KeyMap.Down.SetEnabled(false)

	isLocked := new(atomic.Bool)
	isLocked.Store(false)

	return &Model{
		Model: focus.New(),
		migration: migrator.Migration{
			Steps:          make([]migrator.MigrationStep, 0),
			CurrentVersion: 0,
			IsDirty:        false,
		},
		cursor:   0,
		isLocked: isLocked,
		viewport: viewport,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var (
		cmd tea.Cmd
		agg brownsugar.CmdAggregator
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.IsFocused() {
			break
		}

		switch {
		case key.Matches(msg, Keyj) && !m.isLocked.Load():
			cmd = m.SetCursor(m.cursor + 1)
			agg.Add(cmd)

		case key.Matches(msg, Keyk) && !m.isLocked.Load():
			cmd = m.SetCursor(m.cursor - 1)
			agg.Add(cmd)

		case key.Matches(msg, KeyG) && !m.isLocked.Load():
			cmd = m.SetCursor(m.numberOfSteps() - 1)
			agg.Add(cmd)

		case key.Matches(msg, Keyg) && !m.isLocked.Load():
			cmd = m.SetCursor(0)
			agg.Add(cmd)

		case key.Matches(msg, KeySpace) && m.lock():
			version := m.GetSelectedMigrationStep().Version
			agg.Add(brownsugar.Cmd(appevent.NewMigrateMsg(version)))

		case key.Matches(msg, Keyf) && m.lock():
			if m.GetSelectedMigrationStep().Version > 0 {
				version := m.GetSelectedMigrationStep().Version
				agg.Add(brownsugar.Cmd(appevent.NewForceMigrateMsg(version)))
			} else {
				slog.Error("cannot force migrate to version zero")
			}
		}

	case appevent.UpdateMigrationMsg:
		steps := append([]migrator.MigrationStep{
			{
				Version:    0,
				Identifier: "ROOT (no migration applied)",
			},
		}, msg.Steps...)

		m.migration = migrator.Migration{
			Steps:            steps,
			AppliedMigration: msg.AppliedMigration,
			CurrentVersion:   msg.CurrentVersion,
			IsDirty:          msg.IsDirty,
		}

		agg.Add(m.SetCursor(m.indexOfVersion(m.migration.CurrentVersion)))

		// unlock here as we always update migration state after every operation
		m.unlock()
	}

	if m.IsFocused() {
		m.viewport, cmd = m.viewport.Update(msg)
		agg.Add(cmd)
	}

	return m, tea.Batch(agg...)
}

func (m *Model) Render(ctx brownsugar.Context) string {
	var (
		border = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

		width  = ctx.Width - border.GetHorizontalFrameSize()
		height = ctx.Height - border.GetVerticalFrameSize()
	)

	var longestVersionLength int
	for _, migration := range m.migration.Steps {
		longestVersionLength = max(longestVersionLength, len(strconv.FormatInt(int64(migration.Version), 10)))
	}

	contents := []string{}
	for i, migration := range m.migration.Steps {
		var (
			isDirtyVersion = migration.Version == m.migration.CurrentVersion && m.migration.IsDirty
			isMigrated     = migration.Version <= m.migration.CurrentVersion
			isSelected     = m.cursor == i
		)

		cursor := " "
		if migration.Version == m.migration.CurrentVersion {
			if isDirtyVersion {
				cursor = "✗"
			} else {
				cursor = "▶"
			}
		}

		symbol := " "
		if migration.Version == 0 {
			symbol = "○"
		} else if _, ok := m.migration.AppliedMigration[migration.Signature]; ok {
			symbol = "✔"
		} else {
			symbol = "○"
		}

		// var symbol string
		// switch {
		// case isDirtyVersion:
		// 	symbol = "✗"
		// case migration.Version == 0:
		// 	symbol = "○"
		// case isMigrated:
		// 	symbol = "✔"
		// default:
		// 	symbol = "○"
		// }

		line := fmt.Sprintf("%s %s %d %s",
			cursor,
			symbol,
			migration.Version,
			migration.Identifier,
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

		if isSelected {
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

		contents = append(
			contents,
			style.Width(max(ctx.Width, lineWidth)).Render(line),
		)
	}

	m.viewport.SetWidth(width)
	m.viewport.SetHeight(height)
	m.viewport.SetContent(strings.Join(contents, "\n"))
	m.focusAtCursor()

	return border.
		BorderForeground(m.borderColor()).
		Render(m.viewport.View())
}

func (m *Model) SetCursor(cursor int) tea.Cmd {
	if cursor < 0 {
		cursor = 0
	} else if cursor >= m.numberOfSteps() {
		cursor = m.numberOfSteps() - 1
	}

	if m.cursor != cursor {
		m.cursor = cursor
		return brownsugar.Cmd(appevent.NewSelectMigrationStepMsg(m.migration.Steps[m.cursor]))
	}

	return nil
}

func (m *Model) GetSelectedMigrationStep() migrator.MigrationStep {
	return m.migration.Steps[m.cursor]
}

func (m *Model) borderColor() lipgloss.ANSIColor {
	if m.IsFocused() {
		return brownsugar.ColorYellow
	} else {
		return brownsugar.ColorWhite
	}
}

func (m *Model) indexOfVersion(version uint) int {
	return slices.IndexFunc(m.migration.Steps, func(migration migrator.MigrationStep) bool {
		return migration.Version == version
	})
}

func (m *Model) focusAtCursor() {
	var (
		offset = m.cursor - m.viewport.Height()/2

		topBound    = m.cursor - m.viewport.Height()/2
		bottomBound = m.cursor + m.viewport.Height()/2
	)

	if topBound < 0 {
		offset = 0
	} else if bottomBound >= m.numberOfSteps() {
		offset = topBound
	}

	m.viewport.SetYOffset(offset)
}

func (m *Model) numberOfSteps() int {
	return len(m.migration.Steps)
}

func (m *Model) lock() bool {
	return m.isLocked.CompareAndSwap(false, true)
}

func (m *Model) unlock() {
	m.isLocked.Store(false)
}
