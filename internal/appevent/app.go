package appevent

import (
	tea "charm.land/bubbletea/v2"
	"github.com/LiddleChild/lazymigrate/internal/appscene"
)

type UpdateMigrationRequestMsg struct{}

func NewUpdateMigrationRequestMsg() tea.Msg {
	return UpdateMigrationRequestMsg{}
}

type SwitchSceneMsg struct {
	Scene appscene.Scene
}

func NewSwitchSceneMsg(scene appscene.Scene) tea.Msg {
	return SwitchSceneMsg{
		Scene: scene,
	}
}
