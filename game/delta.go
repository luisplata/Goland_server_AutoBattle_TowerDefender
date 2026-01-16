package game

type UnitMove struct {
	ID int `json:"id"`
	X  int `json:"x"`
	Y  int `json:"y"`
}

type UnitUpdate struct {
	ID       int    `json:"id"`
	TargetID int    `json:"targetId,omitempty"`
	HP       int    `json:"hp,omitempty"`
	Status   string `json:"status,omitempty"`
}

type Delta struct {
	Type             string       `json:"type"`
	Tick             int          `json:"tick"`
	Spawned          []*UnitState `json:"spawned,omitempty"`
	Moved            []UnitMove   `json:"moved,omitempty"`
	Updated          []UnitUpdate `json:"updated,omitempty"`
	Dead             []int        `json:"dead,omitempty"`
	CurrentPhase     GamePhase    `json:"currentPhase"`
	TurnNumber       int          `json:"turnNumber"`
	HumanPlayerID    int          `json:"humanPlayerId"`
	AIPlayerID       int          `json:"aiPlayerId"`
	HumanPlayerReady bool         `json:"humanPlayerReady"`
	AIPlayerReady    bool         `json:"aiPlayerReady"`
	Config           PhaseConfig  `json:"config"` // Configuración de fases
}

func BuildDelta(prev, curr Snapshot) Delta {
	delta := Delta{
		Type:             "delta",
		Tick:             curr.Tick,
		Spawned:          []*UnitState{},
		Moved:            []UnitMove{},
		Updated:          []UnitUpdate{},
		Dead:             []int{},
		CurrentPhase:     curr.CurrentPhase,
		TurnNumber:       curr.TurnNumber,
		HumanPlayerID:    curr.HumanPlayerID,
		AIPlayerID:       curr.AIPlayerID,
		HumanPlayerReady: curr.HumanPlayerReady,
		AIPlayerReady:    curr.AIPlayerReady,
		Config:           curr.Config,
	}

	// Detectar unidades nuevas (spawned)
	for id, currUnit := range curr.Units {
		if _, exists := prev.Units[id]; !exists {
			delta.Spawned = append(delta.Spawned, currUnit)
		}
	}

	// Detectar movimientos y cambios de estado
	for id, currUnit := range curr.Units {
		prevUnit, exists := prev.Units[id]
		if !exists {
			continue
		}

		if currUnit.X != prevUnit.X || currUnit.Y != prevUnit.Y {
			delta.Moved = append(delta.Moved, UnitMove{
				ID: id,
				X:  currUnit.X,
				Y:  currUnit.Y,
			})
		}

		// Detectar cambios de estado (TargetID, HP, Status)
		if currUnit.TargetID != prevUnit.TargetID || currUnit.HP != prevUnit.HP || currUnit.Status != prevUnit.Status {
			update := UnitUpdate{
				ID: id,
			}
			if currUnit.TargetID != prevUnit.TargetID {
				update.TargetID = currUnit.TargetID
			}
			if currUnit.HP != prevUnit.HP {
				update.HP = currUnit.HP
			}
			if currUnit.Status != prevUnit.Status {
				update.Status = currUnit.Status
			}
			delta.Updated = append(delta.Updated, update)
		}
	}

	// Detectar unidades muertas
	for id := range prev.Units {
		if _, exists := curr.Units[id]; !exists {
			delta.Dead = append(delta.Dead, id)
		}
	}

	return delta
}
