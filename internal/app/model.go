package app

import (
	"math"

	"charm.land/bubbles/v2/spinner"
	"github.com/LiddleChild/lazymigrate/internal/app/contentview"
	"github.com/LiddleChild/lazymigrate/internal/app/migrationview"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
	"github.com/LiddleChild/lazymigrate/internal/log"
	"github.com/LiddleChild/lazymigrate/internal/migrator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
)

var _ tea.Model = (*model)(nil)

type model struct {
	migrator *migrator.Migrator

	width  int
	height int

	migrationview *migrationview.Model
	contentview   *contentview.Model
}

func New(migrator *migrator.Migrator) tea.Model {
	migrationview := migrationview.New()
	contentview := contentview.New()

	return &model{
		migrator:      migrator,
		width:         0,
		height:        0,
		migrationview: migrationview,
		contentview:   contentview,
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		migrationview.UpdateMigrationsCmd(m.migrator.GetMigration()),
		m.migrationview.Init(),
		m.contentview.Init(),
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
	default:
		spew.Fdump(log.Entry, msg)
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case migrationview.SelectMigrationStepMsg:
	}

	m.migrationview, cmd = m.migrationview.Update(msg)
	cmds = append(cmds, cmd)

	m.contentview, cmd = m.contentview.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
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
		Width(m.width - borderPane.GetHorizontalFrameSize()).
		Height(bottomH - borderPane.GetVerticalFrameSize())

	return lipgloss.JoinVertical(lipgloss.Top,
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
	)
}
