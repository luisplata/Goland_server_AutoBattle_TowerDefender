package game

type UnitMove struct {
	ID int `json:"id"`
	X  int `json:"x"`
	Y  int `json:"y"`
}

type Delta struct {
	Type    string       `json:"type"`
	Tick    int          `json:"tick"`
	Spawned []*UnitState `json:"spawned,omitempty"`
	Moved   []UnitMove   `json:"moved,omitempty"`
	Dead    []int        `json:"dead,omitempty"`
}

func BuildDelta(prev, curr Snapshot) Delta {
	delta := Delta{
		Type:    "delta",
		Tick:    curr.Tick,
		Spawned: []*UnitState{},
		Moved:   []UnitMove{},
		Dead:    []int{},
	}

	// Detectar unidades nuevas (spawned)
	for id, currUnit := range curr.Units {
		if _, exists := prev.Units[id]; !exists {
			delta.Spawned = append(delta.Spawned, currUnit)
		}
	}

	// Detectar movimientos
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
	}

	// Detectar unidades muertas
	for id := range prev.Units {
		if _, exists := curr.Units[id]; !exists {
			delta.Dead = append(delta.Dead, id)
		}
	}

	return delta
}
