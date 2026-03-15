package brownsugar

import (
	tea "charm.land/bubbletea/v2"
)

var _ ViewModel = (*SceneManager)(nil)

type SceneManager struct {
	currentScene string
	sceneMap     map[string]SceneModel
	initialized  bool
}

func NewSceneManager(currentScene string, scenes ...SceneModel) *SceneManager {
	sceneMap := make(map[string]SceneModel)
	for _, scene := range scenes {
		sceneMap[scene.Scene()] = scene
	}

	return &SceneManager{
		currentScene: currentScene,
		sceneMap:     sceneMap,
		initialized:  false,
	}
}

func (sm *SceneManager) Init() tea.Cmd {
	return tea.Sequence(
		Cmd(NewSwitchSceneMsg(sm.currentScene)),
		sm.initialize,
	)
}

func (sm *SceneManager) Update(msg tea.Msg) (ViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case SwitchSceneMsg:
		sm.currentScene = msg.Scene

		scene, ok := sm.sceneMap[sm.currentScene]
		if ok {
			return sm, scene.Init()
		}

		return sm, nil
	}

	if !sm.initialized {
		return sm, nil
	}

	scene, ok := sm.sceneMap[sm.currentScene]
	if !ok {
		return sm, nil
	}

	var cmd tea.Cmd
	sm.sceneMap[sm.currentScene], cmd = scene.Update(msg)

	return sm, cmd
}

func (sm *SceneManager) Render(ctx Context) string {
	scene, ok := sm.sceneMap[sm.currentScene]
	if !ok {
		return ""
	}

	return scene.Render(ctx)
}

func (sm *SceneManager) initialize() tea.Msg {
	sm.initialized = true
	return nil
}
