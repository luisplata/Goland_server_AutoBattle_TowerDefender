package game

// UpdateMessage provides a unified payload shape for snapshots and deltas.
type UpdateMessage struct {
	Type              string             `json:"type"`
	Tick              int                `json:"tick"`
	Units             map[int]*UnitState `json:"units,omitempty"`
	Players           map[int]*Player    `json:"players,omitempty"`
	Map               *GameMap           `json:"map,omitempty"`
	Spawned           []*UnitState       `json:"spawned,omitempty"`
	Moved             []UnitMove         `json:"moved,omitempty"`
	Dead              []int              `json:"dead,omitempty"`
	CurrentPhase      GamePhase          `json:"currentPhase"`
	TurnNumber        int                `json:"turnNumber"`
	HumanPlayerID     int                `json:"humanPlayerId"`
	AIPlayerID        int                `json:"aiPlayerId"`
	HumanPlayerReady  bool               `json:"humanPlayerReady"`
	AIPlayerReady     bool               `json:"aiPlayerReady"`
	Config            PhaseConfig        `json:"config"` // Configuración de fases
	CurrentPlayerTurn int                `json:"currentPlayerTurn"`
	GameEnd           *GameEndInfo       `json:"gameEnd,omitempty"`
}

// PhaseChangeEvent notifica cuando cambia la fase del juego
type PhaseChangeEvent struct {
	Type          string    `json:"type"`
	Tick          int       `json:"tick"`
	CurrentPhase  GamePhase `json:"currentPhase"`
	PreviousPhase GamePhase `json:"previousPhase"`
	TurnNumber    int       `json:"turnNumber"`
	HumanPlayerID int       `json:"humanPlayerId"`
	AIPlayerID    int       `json:"aiPlayerId"`
}

func SnapshotToUpdate(s Snapshot) UpdateMessage {
	return UpdateMessage{
		Type:              s.Type,
		Tick:              s.Tick,
		Units:             s.Units,
		Players:           s.Players,
		Map:               s.Map,
		CurrentPhase:      s.CurrentPhase,
		TurnNumber:        s.TurnNumber,
		HumanPlayerID:     s.HumanPlayerID,
		AIPlayerID:        s.AIPlayerID,
		HumanPlayerReady:  s.HumanPlayerReady,
		AIPlayerReady:     s.AIPlayerReady,
		Config:            s.Config,
		CurrentPlayerTurn: s.CurrentPlayerTurn,
		GameEnd:           s.GameEnd,
	}
}

func DeltaToUpdate(d Delta) UpdateMessage {
	return UpdateMessage{
		Type:             d.Type,
		Tick:             d.Tick,
		Spawned:          d.Spawned,
		Moved:            d.Moved,
		Dead:             d.Dead,
		CurrentPhase:     d.CurrentPhase,
		TurnNumber:       d.TurnNumber,
		HumanPlayerID:    d.HumanPlayerID,
		AIPlayerID:       d.AIPlayerID,
		HumanPlayerReady: d.HumanPlayerReady,
		AIPlayerReady:    d.AIPlayerReady,
		Config:           d.Config,
	}
}

// BuildPhaseChangeEvent crea un evento de cambio de fase
func BuildPhaseChangeEvent(state *GameState, previousPhase GamePhase) PhaseChangeEvent {
	state.mu.Lock()
	defer state.mu.Unlock()

	return PhaseChangeEvent{
		Type:          "phase_changed",
		Tick:          state.Tick,
		CurrentPhase:  state.CurrentPhase,
		PreviousPhase: previousPhase,
		TurnNumber:    state.TurnNumber,
		HumanPlayerID: state.HumanPlayerID,
		AIPlayerID:    state.AIPlayerID,
	}
}
