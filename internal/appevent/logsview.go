package appevent

import (
	tea "charm.land/bubbletea/v2"
	"github.com/LiddleChild/lazymigrate/internal/log"
)

type LogMessageMsg = log.Message

func NewLogMessageMsg(msg log.Message) tea.Msg {
	return LogMessageMsg(msg)
}
