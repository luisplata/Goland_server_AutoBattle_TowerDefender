package game

type Snapshot struct {
	Type  string             `json:"type"`
	Tick  int                `json:"tick"`
	Units map[int]*UnitState `json:"units"`
}

func BuildSnapshot(state *GameState) Snapshot {
	state.mu.Lock()
	defer state.mu.Unlock()

	// Copiar unidades para evitar race conditions
	unitsCopy := make(map[int]*UnitState, len(state.Units))
	for id, unit := range state.Units {
		unitsCopy[id] = &UnitState{
			ID:       unit.ID,
			PlayerID: unit.PlayerID,
			UnitType: unit.UnitType,
			X:        unit.X,
			Y:        unit.Y,
			HP:       unit.HP,
		}
	}

	return Snapshot{
		Type:  "snapshot",
		Tick:  state.Tick,
		Units: unitsCopy,
	}
}
