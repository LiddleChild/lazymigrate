package migrationview

import (
	"log/slog"
	"slices"
	"sync/atomic"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LiddleChild/lazymigrate/internal/appevent"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
	"github.com/LiddleChild/lazymigrate/internal/components/focus"
	"github.com/LiddleChild/lazymigrate/internal/components/list"
	"github.com/LiddleChild/lazymigrate/internal/migrator"
)

type Model struct {
	focus.Model

	isLocked *atomic.Bool

	list *list.Model
}

func New() *Model {
	isLocked := new(atomic.Bool)
	isLocked.Store(false)

	return &Model{
		Model:    focus.New(),
		isLocked: isLocked,
		list:     list.New(),
	}
}

func (m *Model) Init() tea.Cmd {
	return brownsugar.Cmd(appevent.NewUpdateMigrationRequestMsg())
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
		case key.Matches(msg, KeyMap.Down) && !m.isLocked.Load():
			agg.Add(m.moveCursorCmd(m.list.Down()))

		case key.Matches(msg, KeyMap.Up) && !m.isLocked.Load():
			agg.Add(m.moveCursorCmd(m.list.Up()))

		case key.Matches(msg, KeyMap.Bottom) && !m.isLocked.Load():
			agg.Add(m.moveCursorCmd(m.list.Bottom()))

		case key.Matches(msg, KeyMap.Top) && !m.isLocked.Load():
			agg.Add(m.moveCursorCmd(m.list.Top()))

		case key.Matches(msg, KeyMap.Migrate) && m.lock():
			version := m.getSelectedMigrationStep().Version
			agg.Add(brownsugar.Cmd(appevent.NewMigrateMsg(version)))

		case key.Matches(msg, KeyMap.ForceMigrate) && m.lock():
			if m.getSelectedMigrationStep().Version > 0 {
				version := m.getSelectedMigrationStep().Version
				agg.Add(brownsugar.Cmd(appevent.NewForceMigrateMsg(version)))
			} else {
				slog.Error("cannot force migrate to version zero")
			}
		}

	case appevent.UpdateMigrationMsg:
		steps := slices.Clone(msg.Steps)
		steps = append([]migrator.MigrationStep{
			{
				Version:    0,
				Identifier: "ROOT (no migration applied)",
			},
		}, steps...)

		items := make([]list.Item, 0, len(steps))
		for _, migration := range steps {
			_, ok := msg.AppliedMigration[migration.Signature]

			items = append(items, item{
				step:           migration,
				currentVersion: msg.CurrentVersion,
				isDirty:        msg.IsDirty,
				isApplied:      ok,
			})
		}

		m.list.SetItems(items)
		m.list.SetCursor(m.indexOfVersion(msg.CurrentVersion))

		agg.Add(brownsugar.Cmd(appevent.NewSelectMigrationStepMsg(m.getSelectedMigrationStep())))

		// unlock here as we always update migration state after every operation
		m.unlock()
	}

	m.list, cmd = m.list.Update(msg)
	agg.Add(cmd)

	return m, tea.Batch(agg...)
}

func (m *Model) Render(ctx brownsugar.Context) string {
	m.list.FocusAtCursor()
	m.list.SetBorderForegroundColor(m.borderColor())

	return m.list.Render(brownsugar.Context{
		Width:  ctx.Width,
		Height: ctx.Height,
	})
}

func (m *Model) HelpMenuBindings() []key.Binding {
	return []key.Binding{
		KeyMap.Up,
		KeyMap.Down,
		KeyMap.Top,
		KeyMap.Bottom,
		KeyMap.Migrate,
		KeyMap.ForceMigrate,
		KeyMap.View,
	}
}

func (m *Model) getSelectedMigrationStep() migrator.MigrationStep {
	return m.list.GetSelectedItem().(item).step
}

func (m *Model) moveCursorCmd(moved int) tea.Cmd {
	return func() tea.Msg {
		if moved != 0 {
			return appevent.NewSelectMigrationStepMsg(m.getSelectedMigrationStep())
		}
		return nil
	}
}

func (m *Model) borderColor() lipgloss.ANSIColor {
	if m.IsFocused() {
		return brownsugar.ColorYellow
	} else {
		return brownsugar.ColorWhite
	}
}

func (m *Model) indexOfVersion(version uint) int {
	return slices.IndexFunc(m.list.GetItems(), func(i list.Item) bool {
		return i.(item).step.Version == version
	})
}

func (m *Model) lock() bool {
	return m.isLocked.CompareAndSwap(false, true)
}

func (m *Model) unlock() {
	m.isLocked.Store(false)
}
