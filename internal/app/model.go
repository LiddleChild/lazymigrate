package app

import (
	"math"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LiddleChild/lazymigrate/internal/app/contentview"
	"github.com/LiddleChild/lazymigrate/internal/app/migrationview"
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
	migrator    *migrator.Migrator
	focusedPane FocusedPane

	width  int
	height int

	migrationview *migrationview.Model
	contentview   *contentview.Model
}

func New(migrator *migrator.Migrator) tea.Model {
	migrationview := migrationview.New()
	migrationview.Focus()

	contentview := contentview.New()

	return &model{
		migrator:      migrator,
		focusedPane:   FocusPaneMigration,
		width:         0,
		height:        0,
		migrationview: migrationview,
		contentview:   contentview,
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Sequence(
		tea.Batch(
			m.migrationview.Init(),
			m.contentview.Init(),
		),
		updateMigrationRequestCmd,
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg,
		migrationview.UpdateMigrationMsg:

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

	case updateMigrationRequestMsg:
		migration, err := m.migrator.GetMigration()
		if err != nil {
			panic(err)
		}

		return m, migrationview.UpdateMigrationCmd(migration)

	case migrationview.MigrateMsg:
		if err := m.migrator.MigrateToVersion(msg.Version); err != nil {
			panic(err)
		}

		cmds = append(cmds, updateMigrationRequestCmd)

	case migrationview.ForceMigrateMsg:
		if err := m.migrator.ForceMigrateToVersion(msg.Version); err != nil {
			panic(err)
		}

		cmds = append(cmds, updateMigrationRequestCmd)
	}

	m.migrationview, cmd = m.migrationview.Update(msg)
	cmds = append(cmds, cmd)

	m.contentview, cmd = m.contentview.Update(msg)
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
			bottomPane.Render("Logs"),
		),
	}
}
