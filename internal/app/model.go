package app

import (
	"log/slog"
	"math"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LiddleChild/lazymigrate/internal/app/contentview"
	"github.com/LiddleChild/lazymigrate/internal/app/logsview"
	"github.com/LiddleChild/lazymigrate/internal/app/migrationview"
	"github.com/LiddleChild/lazymigrate/internal/appevent"
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
)

var _ tea.Model = (*model)(nil)

type model struct {
	migrator *migrator.Migrator

	width       int
	height      int
	focusedPane FocusedPane

	migrationview *migrationview.Model
	contentview   *contentview.Model
	logsview      *logsview.Model
}

func New(migrator *migrator.Migrator) tea.Model {
	migrationview := migrationview.New()
	migrationview.Focus()

	contentview := contentview.New()

	logsview := logsview.New()

	return &model{
		migrator:      migrator,
		width:         0,
		height:        0,
		focusedPane:   FocusPaneMigration,
		migrationview: migrationview,
		contentview:   contentview,
		logsview:      logsview,
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Sequence(
		tea.Batch(
			m.migrationview.Init(),
			m.contentview.Init(),
			m.logsview.Init(),
		),
		appevent.UpdateMigrationRequestCmd,
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg,
		appevent.UpdateMigrationMsg:

	default:
		spew.Fdump(log.Entry, msg)
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keyq) || key.Matches(msg, KeyCtrlC):
			m.migrator.Stop()
			return m, tea.Quit

		case key.Matches(msg, KeyEnter):
			m.focusedPane = FocusPaneContent

		case key.Matches(msg, KeyEsc):
			m.focusedPane = FocusPaneMigration
		}

		m.migrationview.Blur()
		m.contentview.Blur()

		switch m.focusedPane {
		case FocusPaneMigration:
			m.migrationview.Focus()

		case FocusPaneContent:
			m.contentview.Focus()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case appevent.UpdateMigrationRequestMsg:
		migration, err := m.migrator.GetMigration()
		if err != nil {
			return m, appevent.ErrCmd(err)
		}

		return m, appevent.UpdateMigrationCmd(migration)

	case appevent.MigrateMsg:
		cmds = append(cmds,
			tea.Sequence(
				m.migrateToVersionCmd(msg.Version),
				appevent.UpdateMigrationRequestCmd,
			),
		)

	case appevent.ForceMigrateMsg:
		cmds = append(cmds,
			tea.Sequence(
				m.forceMigrateToVersionCmd(msg.Version),
				appevent.UpdateMigrationRequestCmd,
			),
		)
	case appevent.ErrMsg:
		slog.Error(msg.Err.Error())
	}

	m.migrationview, cmd = m.migrationview.Update(msg)
	cmds = append(cmds, cmd)

	m.contentview, cmd = m.contentview.Update(msg)
	cmds = append(cmds, cmd)

	m.logsview, cmd = m.logsview.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
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
			top.Render(
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
