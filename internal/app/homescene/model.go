package homescene

import (
	"math"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LiddleChild/lazymigrate/internal/app/homescene/contentview"
	"github.com/LiddleChild/lazymigrate/internal/app/homescene/migrationview"
	"github.com/LiddleChild/lazymigrate/internal/app/homescene/sourceview"
	"github.com/LiddleChild/lazymigrate/internal/appevent"
	"github.com/LiddleChild/lazymigrate/internal/appscene"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
)

type FocusedPane int

const (
	FocusPaneMigration FocusedPane = iota
	FocusPaneContent
)

var _ brownsugar.SceneModel = (*Model)(nil)

type Model struct {
	focusedPane FocusedPane

	migrationview *migrationview.Model
	contentview   *contentview.Model
	sourceview    *sourceview.Model
}

func New() brownsugar.SceneModel {
	var (
		migrationview = migrationview.New()
		contentview   = contentview.New()
		sourceview    = sourceview.New()
	)

	return &Model{
		focusedPane:   FocusPaneMigration,
		migrationview: migrationview,
		contentview:   contentview,
		sourceview:    sourceview,
	}
}

func (m *Model) Scene() string {
	return appscene.SceneHome
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.migrationview.Init(),
		m.contentview.Init(),
		m.sourceview.Init(),
		m.updateFocusedPane(),
	)
}

func (m *Model) Update(msg tea.Msg) (brownsugar.SceneModel, tea.Cmd) {
	var (
		agg brownsugar.CmdAggregator
		cmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, migrationview.KeyMap.View):
			m.focusedPane = FocusPaneContent
			agg.Add(m.updateFocusedPane())

		case key.Matches(msg, contentview.KeyMap.Back):
			m.focusedPane = FocusPaneMigration
			agg.Add(m.updateFocusedPane())

		case key.Matches(msg, KeyMap.NewMigration):
			return m, brownsugar.Cmd(brownsugar.NewSwitchSceneMsg(appscene.SceneNewMigration))

		case key.Matches(msg, KeyMap.ConnectionSource):
			return m, brownsugar.Cmd(brownsugar.NewSwitchSceneMsg(appscene.SceneSources))
		}
	}

	m.migrationview, cmd = m.migrationview.Update(msg)
	agg.Add(cmd)

	m.contentview, cmd = m.contentview.Update(msg)
	agg.Add(cmd)

	m.sourceview, cmd = m.sourceview.Update(msg)
	agg.Add(cmd)

	return m, tea.Batch(agg...)
}

func (m *Model) Render(ctx brownsugar.Context) string {
	leftPane := lipgloss.NewStyle().
		Width(int(math.Round(float64(ctx.Width) / 3))).
		Height(ctx.Height)

	sourcePane := leftPane.
		Width(leftPane.GetWidth()).
		Height(5)

	migrationPane := leftPane.
		Width(leftPane.GetWidth()).
		Height(leftPane.GetHeight() - sourcePane.GetHeight())

	rightPane := lipgloss.NewStyle().
		Width(ctx.Width - leftPane.GetWidth()).
		Height(ctx.Height)

	return lipgloss.JoinHorizontal(lipgloss.Left,
		lipgloss.JoinVertical(lipgloss.Top,
			m.sourceview.Render(brownsugar.Context{
				Width:  sourcePane.GetWidth(),
				Height: sourcePane.GetHeight(),
			}),
			m.migrationview.Render(brownsugar.Context{
				Width:  migrationPane.GetWidth(),
				Height: migrationPane.GetHeight(),
			}),
		),
		m.contentview.Render(brownsugar.Context{
			Width:  rightPane.GetWidth(),
			Height: rightPane.GetHeight(),
		}),
	)
}

func (m *Model) HelpMenuBindings() []key.Binding {
	bindings := []key.Binding{KeyMap.NewMigration, KeyMap.ConnectionSource}

	switch m.focusedPane {
	case FocusPaneMigration:
		bindings = append(bindings, m.migrationview.HelpMenuBindings()...)
	case FocusPaneContent:
		bindings = append(bindings, m.contentview.HelpMenuBindings()...)
	}

	return bindings
}

func (m *Model) updateFocusedPane() tea.Cmd {
	m.migrationview.Blur()
	m.contentview.Blur()

	switch m.focusedPane {
	case FocusPaneMigration:
		m.migrationview.Focus()

	case FocusPaneContent:
		m.contentview.Focus()
	}

	return brownsugar.Cmd(appevent.NewUpdateHelpMenuKeysMsg(m.HelpMenuBindings()))
}
