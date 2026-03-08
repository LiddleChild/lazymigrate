package brownsugar

import tea "charm.land/bubbletea/v2"

type SwitchSceneMsg struct {
	Scene string
}

func NewSwitchSceneMsg(scene string) tea.Msg {
	return SwitchSceneMsg{
		Scene: scene,
	}
}
