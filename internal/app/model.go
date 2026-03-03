package app

import (
	"log/slog"
	"math"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LiddleChild/lazymigrate/internal/app/contentview"
	"github.com/LiddleChild/lazymigrate/internal/app/logsview"
	"github.com/LiddleChild/lazymigrate/internal/app/migrationview"
	"github.com/LiddleChild/lazymigrate/internal/app/newmigrationview"
	"github.com/LiddleChild/lazymigrate/internal/appevent"
	"github.com/LiddleChild/lazymigrate/internal/appscene"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
	"github.com/LiddleChild/lazymigrate/internal/log"
	"github.com/LiddleChild/lazymigrate/internal/migrator"
	"github.com/davecgh/go-spew/spew"
)

type FocusedPane int

const (
	FocusPaneMigration FocusedPane = iota
	FocusPaneContent
	focusPaneEnd
)

func (fp *FocusedPane) Next() {
	*fp = (*fp + 1) % focusPaneEnd
}

var (
	Keyq     = key.NewBinding(key.WithKeys("q"))
	KeyCtrlC = key.NewBinding(key.WithKeys("ctrl+c"))
	KeyEnter = key.NewBinding(key.WithKeys("enter"))
	KeyEsc   = key.NewBinding(key.WithKeys("esc"))
	Keyn     = key.NewBinding(key.WithKeys("n"))
)

var _ tea.Model = (*model)(nil)

type model struct {
	migrator *migrator.Migrator

	width       int
	height      int
	focusedPane FocusedPane
	scene       appscene.Scene

	migrationview    *migrationview.Model
	contentview      *contentview.Model
	logsview         *logsview.Model
	newmigrationview *newmigrationview.Model
}

func New(migrator *migrator.Migrator) tea.Model {
	migrationview := migrationview.New()
	migrationview.Focus()

	contentview := contentview.New()

	logsview := logsview.New()

	newmigrationview := newmigrationview.New()

	return &model{
		migrator:         migrator,
		width:            0,
		height:           0,
		focusedPane:      FocusPaneMigration,
		scene:            appscene.SceneHome,
		migrationview:    migrationview,
		contentview:      contentview,
		logsview:         logsview,
		newmigrationview: newmigrationview,
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Sequence(
		tea.Batch(
			m.migrationview.Init(),
			m.contentview.Init(),
			m.logsview.Init(),
		),
		brownsugar.Cmd(appevent.NewUpdateMigrationRequestMsg()),
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg,
		tea.KeyPressMsg,
		cursor.BlinkMsg,
		appevent.UpdateMigrationMsg:

	default:
		spew.Fdump(log.Entry, msg)
	}

	var cmd tea.Cmd
	agg := brownsugar.CmdAggregator{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keyq) || key.Matches(msg, KeyCtrlC):
			m.migrator.Stop()
			return m, tea.Quit

		case m.scene == appscene.SceneHome && key.Matches(msg, KeyEnter):
			m.focusedPane = FocusPaneContent
			m.updateFocusedPane()

		case m.scene == appscene.SceneHome && key.Matches(msg, KeyEsc):
			m.focusedPane = FocusPaneMigration
			m.updateFocusedPane()

		case m.scene == appscene.SceneHome && key.Matches(msg, Keyn):
			return m, brownsugar.Cmd(appevent.NewSwitchSceneMsg(appscene.SceneNewMigration))
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case appevent.UpdateMigrationRequestMsg:
		migration, err := m.migrator.GetMigration()
		if err != nil {
			return m, brownsugar.Cmd(appevent.NewErrMsg(err))
		}

		return m, brownsugar.Cmd(appevent.NewUpdateMigrationMsg(migration))

	case appevent.MigrateMsg:
		agg.Add(tea.Sequence(
			m.migrateToVersionCmd(msg.Version),
			brownsugar.Cmd(appevent.NewUpdateMigrationRequestMsg()),
		))

	case appevent.ForceMigrateMsg:
		agg.Add(tea.Sequence(
			m.forceMigrateToVersionCmd(msg.Version),
			brownsugar.Cmd(appevent.NewUpdateMigrationRequestMsg()),
		))

	case appevent.CreateMigrationMsg:
		if err := m.migrator.CreateMigration(msg.Name); err != nil {
			return m, brownsugar.Cmd(appevent.NewErrMsg(err))
		}

		return m, tea.Sequence(
			brownsugar.Cmd(appevent.NewMigrationCreatedMsg()),
			brownsugar.Cmd(appevent.NewUpdateMigrationRequestMsg()),
		)

	case appevent.SwitchSceneMsg:
		m.scene = msg.Scene
		return m, nil

	case appevent.ErrMsg:
		slog.Error(msg.Err.Error())
		return m, nil
	}

	switch m.scene {
	case appscene.SceneHome:
		m.migrationview, cmd = m.migrationview.Update(msg)
		agg.Add(cmd)

		m.contentview, cmd = m.contentview.Update(msg)
		agg.Add(cmd)

	case appscene.SceneNewMigration:
		m.newmigrationview, cmd = m.newmigrationview.Update(msg)
		agg.Add(cmd)
	}

	m.logsview, cmd = m.logsview.Update(msg)
	agg.Add(cmd)

	return m, tea.Batch(agg...)
}

func (m *model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("")
	}

	topH := int(math.Round(float64(m.height) * 2 / 3))
	bottomH := m.height - topH

	borderPane := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

	top := lipgloss.NewStyle().
		Width(m.width).
		Height(topH)

	migrationPane := lipgloss.NewStyle().
		Width(int(math.Round(float64(top.GetWidth()) / 3))).
		Height(top.GetHeight())

	contentPane := borderPane.
		Width(top.GetWidth() - migrationPane.GetWidth()).
		Height(top.GetHeight())

	bottomPane := borderPane.
		Width(m.width).
		Height(bottomH)

	return tea.View{
		AltScreen: true,
		Content: lipgloss.JoinVertical(lipgloss.Top,
			brownsugar.RenderWithCondition(m.scene == appscene.SceneNewMigration,
				func() string {
					return m.newmigrationview.Render(brownsugar.RenderContext{
						Width:  top.GetWidth(),
						Height: top.GetHeight(),
					})
				},
			),
			brownsugar.RenderWithCondition(m.scene == appscene.SceneHome,
				func() string {
					return top.Render(
						lipgloss.JoinHorizontal(lipgloss.Left,
							m.migrationview.Render(brownsugar.RenderContext{
								Width:  migrationPane.GetWidth(),
								Height: migrationPane.GetHeight(),
							}),
							m.contentview.Render(brownsugar.RenderContext{
								Width:  contentPane.GetWidth(),
								Height: contentPane.GetHeight(),
							}),
						),
					)
				},
			),
			m.logsview.Render(brownsugar.RenderContext{
				Width:  bottomPane.GetWidth(),
				Height: bottomPane.GetHeight(),
			}),
		),
	}
}

func (m *model) migrateToVersionCmd(version uint) tea.Cmd {
	return func() tea.Msg {
		if err := m.migrator.MigrateToVersion(version); err != nil {
			return appevent.ErrMsg{Err: err}
		}

		return nil
	}
}

func (m *model) forceMigrateToVersionCmd(version uint) tea.Cmd {
	return func() tea.Msg {
		if err := m.migrator.ForceMigrateToVersion(version); err != nil {
			return appevent.ErrMsg{Err: err}
		}

		return nil
	}
}

func (m *model) updateFocusedPane() {
	m.migrationview.Blur()
	m.contentview.Blur()

	switch m.focusedPane {
	case FocusPaneMigration:
		m.migrationview.Focus()

	case FocusPaneContent:
		m.contentview.Focus()
	}
}
