package app

import (
	"log/slog"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/LiddleChild/lazymigrate/internal/app/homescene"
	"github.com/LiddleChild/lazymigrate/internal/app/newmigrationscene"
	"github.com/LiddleChild/lazymigrate/internal/app/sourcesscene"
	"github.com/LiddleChild/lazymigrate/internal/appevent"
	"github.com/LiddleChild/lazymigrate/internal/appscene"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
	"github.com/LiddleChild/lazymigrate/internal/log"
	"github.com/LiddleChild/lazymigrate/internal/migrator"
	"github.com/LiddleChild/lazymigrate/internal/source"
	"github.com/davecgh/go-spew/spew"
)

var (
	Keyq     = key.NewBinding(key.WithKeys("q"))
	KeyCtrlC = key.NewBinding(key.WithKeys("ctrl+c"))
)

var _ tea.Model = (*model)(nil)

type model struct {
	migrator      *migrator.Migrator
	sourceManager *source.Manager

	width  int
	height int

	sceneManager brownsugar.ViewModel
}

func New(migrator *migrator.Migrator, sourceManager *source.Manager) tea.Model {
	sceneManager := brownsugar.NewSceneManager(
		appscene.SceneHome,
		homescene.New(),
		newmigrationscene.New(),
		sourcesscene.New(),
	)

	return &model{
		migrator:      migrator,
		sourceManager: sourceManager,
		width:         0,
		height:        0,
		sceneManager:  sceneManager,
	}
}

func (m *model) Init() tea.Cmd {
	return m.sceneManager.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg,
		tea.KeyPressMsg,
		cursor.BlinkMsg,
		appevent.UpdateMigrationMsg:

	default:
		// TODO: use slog here
		spew.Fdump(log.Entry, msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keyq) || key.Matches(msg, KeyCtrlC):
			m.migrator.Stop()
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m, nil

	case appevent.UpdateMigrationRequestMsg:
		migration, err := m.migrator.GetMigration()
		if err != nil {
			return m, brownsugar.Cmd(appevent.NewErrMsg(err))
		}

		return m, brownsugar.Cmd(appevent.NewUpdateMigrationMsg(migration))

	case appevent.UpdateSourcesRequestMsg:
		return m, brownsugar.Cmd(appevent.NewUpdateSourcesMsg(
			m.sourceManager.GetCurrentSourceIndex(),
			m.sourceManager.ListSources(),
		))

	case appevent.MigrateMsg:
		return m, tea.Sequence(
			m.migrateToVersionCmd(msg.Version),
			brownsugar.Cmd(appevent.NewUpdateMigrationRequestMsg()),
		)

	case appevent.ForceMigrateMsg:
		return m, tea.Sequence(
			m.forceMigrateToVersionCmd(msg.Version),
			brownsugar.Cmd(appevent.NewUpdateMigrationRequestMsg()),
		)

	case appevent.CreateMigrationMsg:
		if err := m.migrator.CreateMigration(msg.Name); err != nil {
			return m, brownsugar.Cmd(appevent.NewErrMsg(err))
		}

		return m, tea.Sequence(
			brownsugar.Cmd(appevent.NewMigrationCreatedMsg()),
			brownsugar.Cmd(appevent.NewUpdateMigrationRequestMsg()),
		)

	case appevent.ChangeMigratorSourceMsg:
		source := source.Source(msg)

		if err := m.migrator.Open(source); err != nil {
			return m, brownsugar.Cmd(appevent.NewErrMsg(err))
		}

		m.sourceManager.SetCurrentSource(source)

		return m, brownsugar.Cmd(appevent.NewUpdateSourcesRequestMsg())

	case appevent.ErrMsg:
		slog.Error(msg.Err.Error())
		return m, nil
	}

	var cmd tea.Cmd
	m.sceneManager, cmd = m.sceneManager.Update(msg)

	return m, cmd
}

func (m *model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("")
	}

	return tea.View{
		AltScreen: true,
		Content: m.sceneManager.Render(brownsugar.Context{
			Width:  m.width,
			Height: m.height,
		}),
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
