package appevent

import (
	tea "charm.land/bubbletea/v2"
	"github.com/LiddleChild/lazymigrate/internal/source"
)

type UpdateSourcesMsg struct {
	CurrentSourceIndex int
	Sources            []source.Source
}

func NewUpdateSourcesMsg(index int, sources []source.Source) tea.Msg {
	return UpdateSourcesMsg{
		CurrentSourceIndex: index,
		Sources:            sources,
	}
}

type ChangeMigratorSourceMsg source.Source

func NewChangeMigratorSourceMsg(source source.Source) tea.Msg {
	return ChangeMigratorSourceMsg(source)
}
