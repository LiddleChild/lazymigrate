package app

import "charm.land/bubbles/v2/key"

type keyMap struct {
	Quit key.Binding
}

var KeyMap = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
