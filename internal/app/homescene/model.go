package homescene

import (
	"math"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LiddleChild/lazymigrate/internal/app/homescene/contentview"
	"github.com/LiddleChild/lazymigrate/internal/app/homescene/logsview"
	"github.com/LiddleChild/lazymigrate/internal/app/homescene/migrationview"
	"github.com/LiddleChild/lazymigrate/internal/appscene"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
)

type FocusedPane int

const (
	FocusPaneMigration FocusedPane = iota
	FocusPaneContent
	FocusPaneLog
)

var (
	KeyEnter = key.NewBinding(key.WithKeys("enter"))
	KeyEsc   = key.NewBinding(key.WithKeys("esc"))
	Keyn     = key.NewBinding(key.WithKeys("n"))
	Keyl     = key.NewBinding(key.WithKeys("l"))
)

var _ brownsugar.SceneModel = (*Model)(nil)

type Model struct {
	focusedPane FocusedPane

	migrationview *migrationview.Model
	contentview   *contentview.Model
	logsview      *logsview.Model
}

func New() brownsugar.SceneModel {
	migrationview := migrationview.New()
	migrationview.Focus()

	contentview := contentview.New()

	logsview := logsview.New()

	return &Model{
		focusedPane:   FocusPaneMigration,
		migrationview: migrationview,
		contentview:   contentview,
		logsview:      logsview,
	}
}

func (m *Model) Scene() string {
	return appscene.SceneHome
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.migrationview.Init(),
		m.contentview.Init(),
		m.logsview.Init(),
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
		case key.Matches(msg, KeyEnter):
			m.focusedPane = FocusPaneContent
			m.updateFocusedPane()

		case key.Matches(msg, Keyl):
			m.focusedPane = FocusPaneLog
			m.updateFocusedPane()

		case key.Matches(msg, KeyEsc):
			m.focusedPane = FocusPaneMigration
			m.updateFocusedPane()

		case key.Matches(msg, Keyn):
			return m, brownsugar.Cmd(brownsugar.NewSwitchSceneMsg(appscene.SceneNewMigration))
		}
	}

	m.migrationview, cmd = m.migrationview.Update(msg)
	agg.Add(cmd)

	m.contentview, cmd = m.contentview.Update(msg)
	agg.Add(cmd)

	m.logsview, cmd = m.logsview.Update(msg)
	agg.Add(cmd)

	return m, tea.Batch(agg...)
}

func (m *Model) Render(ctx brownsugar.Context) string {
	var (
		border = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

		topHeight    = int(math.Round(float64(ctx.Height) * 2 / 3))
		bottomHeight = ctx.Height - topHeight
	)

	topPane := lipgloss.NewStyle().
		Width(ctx.Width).
		Height(topHeight)

	bottomPane := lipgloss.NewStyle().
		Width(ctx.Width).
		Height(bottomHeight)

	migrationPane := lipgloss.NewStyle().
		Width(int(math.Round(float64(topPane.GetWidth()) / 3))).
		Height(topPane.GetHeight())

	contentPane := border.
		Width(topPane.GetWidth() - migrationPane.GetWidth()).
		Height(topPane.GetHeight())

	return lipgloss.JoinVertical(lipgloss.Top,
		topPane.Render(
			lipgloss.JoinHorizontal(lipgloss.Left,
				m.migrationview.Render(brownsugar.Context{
					Width:  migrationPane.GetWidth(),
					Height: migrationPane.GetHeight(),
				}),
				m.contentview.Render(brownsugar.Context{
					Width:  contentPane.GetWidth(),
					Height: contentPane.GetHeight(),
				}),
			),
			m.logsview.Render(brownsugar.Context{
				Width:  bottomPane.GetWidth(),
				Height: bottomPane.GetHeight(),
			}),
		),
	)
}

func (m *Model) updateFocusedPane() {
	m.migrationview.Blur()
	m.contentview.Blur()
	m.logsview.Blur()

	switch m.focusedPane {
	case FocusPaneMigration:
		m.migrationview.Focus()

	case FocusPaneContent:
		m.contentview.Focus()

	case FocusPaneLog:
		m.logsview.Focus()
	}
}
