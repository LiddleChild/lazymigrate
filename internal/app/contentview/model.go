package contentview

import (
	"os"
	"strconv"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LiddleChild/lazymigrate/internal/appevent"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
	"github.com/LiddleChild/lazymigrate/internal/components/focus"
	"github.com/LiddleChild/lazymigrate/internal/migrator"
)

var (
	Keyj = key.NewBinding(key.WithKeys("j"))
	Keyk = key.NewBinding(key.WithKeys("k"))
	Keyg = key.NewBinding(key.WithKeys("g"))
	KeyG = key.NewBinding(key.WithKeys("G"))
)

type content struct {
	name    string
	content string
}

type Model struct {
	focus.Model

	step             migrator.MigrationStep
	isLoadingContent bool
	isZeroVersion    bool

	upContent   content
	downContent content

	viewport viewport.Model
	spinner  spinner.Model
}

func New() *Model {
	viewport := viewport.New()

	s := spinner.New()
	s.Spinner = spinner.MiniDot

	return &Model{
		Model:            focus.New(),
		step:             migrator.MigrationStep{},
		isLoadingContent: true,
		isZeroVersion:    false,
		upContent:        content{},
		downContent:      content{},
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
	case appevent.SelectMigrationStepMsg:
		m.step = msg.MigrationStep

		if !m.isLoadingContent {
			cmds = append(cmds, m.spinner.Tick)
		}
		m.isLoadingContent = true

		cmd = tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg {
			// cursor is copied into closure (old value)
			// if current value mismatched with old value, debounce
			if m.step != msg.MigrationStep {
				return nil
			}

			return appevent.UpdateMigrationContentMsg{
				MigrationStep: msg.MigrationStep,
			}
		})

		cmds = append(cmds, cmd)

	case appevent.UpdateMigrationContentMsg:
		if msg.MigrationStep.Up != nil {
			buffer, err := os.ReadFile(msg.MigrationStep.Up.Path)
			if err != nil {
				return m, appevent.ErrCmd(err)
			}

			m.upContent.name = msg.MigrationStep.Up.Fullname
			m.upContent.content = string(buffer)
		} else {
			m.upContent = content{}
		}

		if msg.MigrationStep.Down != nil {
			buffer, err := os.ReadFile(msg.MigrationStep.Down.Path)
			if err != nil {
				return m, appevent.ErrCmd(err)
			}

			m.downContent.name = msg.MigrationStep.Down.Fullname
			m.downContent.content = string(buffer)
		} else {
			m.downContent = content{}
		}

		m.isZeroVersion = msg.MigrationStep.Version == 0
		m.isLoadingContent = false

	case spinner.TickMsg:
		if !m.isLoadingContent {
			return m, nil
		}

		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	if m.IsFocused() {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

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

	m.viewport.SetWidth(width)
	m.viewport.SetHeight(height)

	filename := lipgloss.NewStyle().
		Foreground(brownsugar.ColorBlack).
		Background(brownsugar.ColorBrightWhite)

	switch {
	case m.isLoadingContent:
		m.viewport.SetContent(spinner)

	case !m.isLoadingContent && m.isZeroVersion:
		m.viewport.SetContent("")

	case !m.isLoadingContent && !m.isZeroVersion:
		m.viewport.SetContent(
			lipgloss.JoinVertical(lipgloss.Top,
				filename.Render(m.upContent.name),
				m.renderWithLineNumber(m.upContent.content),
				"",
				filename.Render(m.downContent.name),
				m.renderWithLineNumber(m.downContent.content),
			),
		)
	}

	return border.
		BorderForeground(m.borderColor()).
		Render(m.viewport.View())
}

func (m *Model) borderColor() lipgloss.ANSIColor {
	if m.IsFocused() {
		return brownsugar.ColorYellow
	} else {
		return brownsugar.ColorWhite
	}
}

func (m *Model) renderWithLineNumber(s string) string {
	count := strings.Count(s, "\n") + 1

	mx := len(strconv.FormatInt(int64(count), 10))
	style := lipgloss.NewStyle().
		Foreground(brownsugar.ColorBrightBlack).
		BorderRight(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(brownsugar.ColorBrightBlack).
		MarginRight(1).
		MarginLeft(1).
		Align(lipgloss.Right)

	style = style.
		Width(mx + style.GetBorderRightSize())

	arr := make([]string, 0, count)
	for i := range count {
		arr = append(arr, style.Render(strconv.FormatInt(int64(i+1), 10)))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left,
		strings.Join(arr, "\n"),
		s,
	)
}
