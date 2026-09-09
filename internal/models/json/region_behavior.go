package json

type RegionBehavior struct {
	AdjustDarknessLevel              any `json:"adjustDarknessLevel,omitempty"`
	DisplayScrollingText             any `json:"displayScrollingText,omitempty"`
	ExecuteMacro                     any `json:"executeMacro,omitempty"`
	ExecuteScript                    any `json:"executeScript,omitempty"`
	ModifyMovementCost               any `json:"modifyMovementCost,omitempty"`
	PauseGame                        any `json:"pauseGame,omitempty"`
	SuppressWeather                  any `json:"suppressWeather,omitempty"`
	TeleportToken                    any `json:"teleportToken,omitempty"`
	ToggleBehavior                   any `json:"toggleBehavior,omitempty"`
	RegionBehaviorEnvironment        any `json:"environment,omitempty"`
	RegionBehaviorEnvironmentFeature any `json:"environmentFeature,omitempty"`
}
