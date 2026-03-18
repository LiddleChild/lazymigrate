package migrationview

import "charm.land/bubbles/v2/key"

type keyMap struct {
	Top          key.Binding
	Up           key.Binding
	Down         key.Binding
	Bottom       key.Binding
	Migrate      key.Binding
	ForceMigrate key.Binding
	View         key.Binding
}

var KeyMap = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Top: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "top"),
	),
	Bottom: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "bottom"),
	),
	Migrate: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("space", "migrate"),
	),
	ForceMigrate: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "force"),
	),
	View: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "view"),
	),
}
