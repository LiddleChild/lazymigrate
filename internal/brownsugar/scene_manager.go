package brownsugar

import tea "charm.land/bubbletea/v2"

var _ ViewModel = (*SceneManager)(nil)

type SceneManager struct {
	currentScene string
	sceneMap     map[string]SceneModel
}

func NewSceneManager(currentScene string, scenes ...SceneModel) *SceneManager {
	sceneMap := make(map[string]SceneModel)
	for _, scene := range scenes {
		sceneMap[scene.Scene()] = scene
	}

	return &SceneManager{
		currentScene: currentScene,
		sceneMap:     sceneMap,
	}
}

func (sm *SceneManager) Init() tea.Cmd {
	agg := CmdAggregator{}
	for _, scene := range sm.sceneMap {
		agg.Add(scene.Init())
	}

	return tea.Batch(agg...)
}

func (sm *SceneManager) Update(msg tea.Msg) (ViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case SwitchSceneMsg:
		sm.currentScene = msg.Scene
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
