package brownsugar

import (
	teav2 "charm.land/bubbletea/v2"
	tea "github.com/charmbracelet/bubbletea"
)

func ToCmdV1(cmd teav2.Cmd) tea.Cmd {
	return func() tea.Msg {
		return cmd()
	}
}
