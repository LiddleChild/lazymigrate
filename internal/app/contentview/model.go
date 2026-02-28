package contentview

import (
	"os"
	"strconv"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	teav2 "charm.land/bubbletea/v2"
	"github.com/LiddleChild/lazymigrate/internal/app/migrationview"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
	"github.com/LiddleChild/lazymigrate/internal/log"
	"github.com/LiddleChild/lazymigrate/internal/migrator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
)

type content struct {
	name    string
	content string
}

type Model struct {
	step             migrator.MigrationStep
	isLoadingContent bool

	upContent   content
	downContent content

	viewport viewport.Model
	spinner  spinner.Model
}

func New() *Model {
	viewport := viewport.New(0, 0)
	viewport.KeyMap.Up.SetEnabled(false)
	viewport.KeyMap.Down.SetEnabled(false)

	s := spinner.New()
	s.Spinner = spinner.MiniDot

	return &Model{
		step:             migrator.MigrationStep{},
		isLoadingContent: true,
		viewport:         viewport,
		spinner:          s,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case migrationview.SelectMigrationStepMsg:
		m.step = msg.MigrationStep

		if !m.isLoadingContent {
			cmds = append(cmds, brownsugar.ToCmdV1(m.spinner.Tick))
		}
		m.isLoadingContent = true

		cmd = tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg {
			// cursor is copied into closure (old value)
			// if current value mismatched with old value, debounce
			if m.step != msg.MigrationStep {
				return nil
			}

			return updateMigrationContentMsg{
				MigrationStep: msg.MigrationStep,
			}
		})

		cmds = append(cmds, cmd)

	case updateMigrationContentMsg:
		if msg.MigrationStep.Up != nil {
			buffer, err := os.ReadFile(msg.MigrationStep.Up.Path)
			if err != nil {
				panic(err)
			}

			m.upContent.name = msg.MigrationStep.Up.Fullname
			m.upContent.content = string(buffer)
		} else {
			m.upContent = content{}
		}

		if msg.MigrationStep.Down != nil {
			buffer, err := os.ReadFile(msg.MigrationStep.Down.Path)
			if err != nil {
				panic(err)
			}

			m.downContent.name = msg.MigrationStep.Down.Fullname
			m.downContent.content = string(buffer)
		} else {
			m.downContent = content{}
		}

		m.isLoadingContent = false

	case spinner.TickMsg:
		if !m.isLoadingContent {
			return m, nil
		}

		var cmd teav2.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, brownsugar.ToCmdV1(cmd)

	default:
		return m, nil
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) Render(ctx brownsugar.RenderContext) string {
	var (
		border = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

		width  = ctx.Width - border.GetHorizontalFrameSize()
		height = ctx.Height - border.GetVerticalFrameSize()
	)

	spinner := lipgloss.NewStyle().
		Width(width).
		Height(height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(m.spinner.View())

	m.viewport.Width = width
	m.viewport.Height = height

	filename := lipgloss.NewStyle().
		Foreground(brownsugar.ColorBlack).
		Background(brownsugar.ColorBrightWhite)

	if m.isLoadingContent {
		m.viewport.SetContent(spinner)
	} else {
		m.viewport.SetContent(
			lipgloss.JoinVertical(lipgloss.Top,
				filename.Render(m.upContent.name),
				lipgloss.JoinHorizontal(lipgloss.Left,
					m.lineNumber(m.upContent.content),
					m.upContent.content,
				),
				"",
				filename.Render(m.downContent.name),
				lipgloss.JoinHorizontal(lipgloss.Left,
					m.lineNumber(m.upContent.content),
					m.downContent.content,
				),
			),
		)
	}

	return border.Render(m.viewport.View())
}

func (m *Model) lineNumber(s string) string {
	count := strings.Count(s, "\n") + 1

	mx := len(strconv.FormatInt(int64(count), 10))
	style := lipgloss.NewStyle().
		Foreground(brownsugar.ColorBrightBlack).
		BorderRight(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(brownsugar.ColorBrightBlack).
		MarginRight(1).
		MarginLeft(1).
		Width(mx).
		Align(lipgloss.Right)

	arr := make([]string, 0, count)
	for i := range count {
		arr = append(arr, style.Render(strconv.FormatInt(int64(i+1), 10)))
	}

	spew.Fdump(log.Entry, arr)

	return strings.Join(arr, "\n")
}
