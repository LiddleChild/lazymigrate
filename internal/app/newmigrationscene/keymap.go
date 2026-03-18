package newmigrationscene

import "charm.land/bubbles/v2/key"

type keyMap struct {
	Back   key.Binding
	Create key.Binding
}

var KeyMap = keyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Create: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "create"),
	),
}
