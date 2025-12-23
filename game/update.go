package game

// UpdateMessage provides a unified payload shape for snapshots and deltas.
type UpdateMessage struct {
	Type    string             `json:"type"`
	Tick    int                `json:"tick"`
	Units   map[int]*UnitState `json:"units,omitempty"`
	Spawned []*UnitState       `json:"spawned,omitempty"`
	Moved   []UnitMove         `json:"moved,omitempty"`
	Dead    []int              `json:"dead,omitempty"`
}

func SnapshotToUpdate(s Snapshot) UpdateMessage {
	return UpdateMessage{
		Type:  s.Type,
		Tick:  s.Tick,
		Units: s.Units,
	}
}

func DeltaToUpdate(d Delta) UpdateMessage {
	return UpdateMessage{
		Type:    d.Type,
		Tick:    d.Tick,
		Spawned: d.Spawned,
		Moved:   d.Moved,
		Dead:    d.Dead,
	}
}
