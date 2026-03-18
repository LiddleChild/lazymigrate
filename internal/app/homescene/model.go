package homescene

import (
	"math"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LiddleChild/lazymigrate/internal/app/homescene/contentview"
	"github.com/LiddleChild/lazymigrate/internal/app/homescene/migrationview"
	"github.com/LiddleChild/lazymigrate/internal/appscene"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
)

type FocusedPane int

const (
	FocusPaneMigration FocusedPane = iota
	FocusPaneContent
)

var (
	KeyEnter = key.NewBinding(key.WithKeys("enter"))
	KeyEsc   = key.NewBinding(key.WithKeys("esc"))
	Keyn     = key.NewBinding(key.WithKeys("n"))
	Keyc     = key.NewBinding(key.WithKeys("c"))
)

var _ brownsugar.SceneModel = (*Model)(nil)

type Model struct {
	focusedPane FocusedPane

	migrationview *migrationview.Model
	contentview   *contentview.Model
}

func New() brownsugar.SceneModel {
	migrationview := migrationview.New()
	migrationview.Focus()

	contentview := contentview.New()

	return &Model{
		focusedPane:   FocusPaneMigration,
		migrationview: migrationview,
		contentview:   contentview,
	}
}

func (m *Model) Scene() string {
	return appscene.SceneHome
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.migrationview.Init(),
		m.contentview.Init(),
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

		case key.Matches(msg, KeyEsc):
			m.focusedPane = FocusPaneMigration
			m.updateFocusedPane()

		case key.Matches(msg, Keyn):
			return m, brownsugar.Cmd(brownsugar.NewSwitchSceneMsg(appscene.SceneNewMigration))

		case key.Matches(msg, Keyc):
			return m, brownsugar.Cmd(brownsugar.NewSwitchSceneMsg(appscene.SceneSources))
		}
	}

	m.migrationview, cmd = m.migrationview.Update(msg)
	agg.Add(cmd)

	m.contentview, cmd = m.contentview.Update(msg)
	agg.Add(cmd)

	return m, tea.Batch(agg...)
}

func (m *Model) Render(ctx brownsugar.Context) string {
	migrationPane := lipgloss.NewStyle().
		Width(int(math.Round(float64(ctx.Width) / 3))).
		Height(ctx.Height)

	contentPane := lipgloss.NewStyle().
		Width(ctx.Width - migrationPane.GetWidth()).
		Height(ctx.Height)

	return lipgloss.JoinHorizontal(lipgloss.Left,
		m.migrationview.Render(brownsugar.Context{
			Width:  migrationPane.GetWidth(),
			Height: migrationPane.GetHeight(),
		}),
		m.contentview.Render(brownsugar.Context{
			Width:  contentPane.GetWidth(),
			Height: contentPane.GetHeight(),
		}),
	)
}

func (m *Model) updateFocusedPane() {
	m.migrationview.Blur()
	m.contentview.Blur()

	switch m.focusedPane {
	case FocusPaneMigration:
		m.migrationview.Focus()

	case FocusPaneContent:
		m.contentview.Focus()
	}
}
