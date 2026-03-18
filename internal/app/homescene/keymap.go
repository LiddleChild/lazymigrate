package homescene

import "charm.land/bubbles/v2/key"

type keyMap struct {
	NewMigration     key.Binding
	ConnectionSource key.Binding
}

var KeyMap = keyMap{
	NewMigration: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new migration"),
	),
	ConnectionSource: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "connection"),
	),
}
