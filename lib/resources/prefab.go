package resources

import "github.com/x-hgg-x/goecsengine/loader"

// MenuPrefabs contains menu prefabs
type MenuPrefabs struct {
	MainMenu          loader.EntityComponentList
	PauseMenu         loader.EntityComponentList
	GameOverMenu      loader.EntityComponentList
	LevelCompleteMenu loader.EntityComponentList
}

// GamePrefabs contains game prefabs
type GamePrefabs struct {
	Background loader.EntityComponentList
	Game       loader.EntityComponentList
	Score      loader.EntityComponentList
	Life       loader.EntityComponentList
}

// Prefabs contains menu and game prefabs
type Prefabs struct {
	Menu MenuPrefabs
	Game GamePrefabs
}
