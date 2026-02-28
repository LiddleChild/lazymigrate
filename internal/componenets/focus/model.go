package focus

type Model struct {
	isFocused bool
}

func New() Model {
	return Model{
		isFocused: false,
	}
}

func (m *Model) Focus() {
	m.isFocused = true
}

func (m *Model) Blur() {
	m.isFocused = false
}

func (m *Model) IsFocused() bool {
	return m.isFocused
}
