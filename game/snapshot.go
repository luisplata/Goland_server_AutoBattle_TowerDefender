package game

type Snapshot struct {
	Type  string             `json:"type"`
	Tick  int                `json:"tick"`
	Units map[int]*UnitState `json:"units"`
}

func BuildSnapshot(state *GameState) Snapshot {
	return Snapshot{
		Type:  "snapshot",
		Tick:  state.Tick,
		Units: state.Units,
	}
}
