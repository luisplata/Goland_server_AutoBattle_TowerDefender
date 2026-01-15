package game

type Snapshot struct {
	Type              string             `json:"type"`
	Tick              int                `json:"tick"`
	Units             map[int]*UnitState `json:"units"`
	Players           map[int]*Player    `json:"players"`
	Map               *GameMap           `json:"map"`
	CurrentPhase      GamePhase          `json:"currentPhase"`
	TurnNumber        int                `json:"turnNumber"`
	HumanPlayerID     int                `json:"humanPlayerId"`
	AIPlayerID        int                `json:"aiPlayerId"`
	HumanPlayerReady  bool               `json:"humanPlayerReady"`
	AIPlayerReady     bool               `json:"aiPlayerReady"`
	HumanBaseID       int                `json:"humanBaseId"`
	AIBaseID          int                `json:"aiBaseId"`
	Config            PhaseConfig        `json:"config"`            // Configuración de fases
	CurrentPlayerTurn int                `json:"currentPlayerTurn"` // ID del jugador cuyo turno es
}

func BuildSnapshot(state *GameState) Snapshot {
	state.mu.Lock()
	defer state.mu.Unlock()

	// Copiar unidades para evitar race conditions
	unitsCopy := make(map[int]*UnitState, len(state.Units))
	for id, unit := range state.Units {
		unitsCopy[id] = &UnitState{
			ID:                unit.ID,
			PlayerID:          unit.PlayerID,
			UnitType:          unit.UnitType,
			X:                 unit.X,
			Y:                 unit.Y,
			HP:                unit.HP,
			MaxHP:             unit.MaxHP,
			AttackDamage:      unit.AttackDamage,
			AttackRange:       unit.AttackRange,
			DetectionRange:    unit.DetectionRange,
			Status:            unit.Status,
			BuildRange:        unit.BuildRange,
			Category:          unit.Category,
			IsGenerator:       unit.IsGenerator,
			GeneratedUnitType: unit.GeneratedUnitType,
		}
	}

	// Copiar jugadores para evitar race conditions
	playersCopy := make(map[int]*Player, len(state.Players))
	for id, player := range state.Players {
		handCopy := make([]string, len(player.Hand))
		copy(handCopy, player.Hand)
		playersCopy[id] = &Player{
			ID:        player.ID,
			IsAI:      player.IsAI,
			Hand:      handCopy,
			DeckCount: player.DeckCount,
			Connected: player.Connected,
		}
	}

	// Determinar quién es el jugador actual del turno
	currentPlayerTurn := state.HumanPlayerID
	if state.AIPlayerID > 0 && state.TurnNumber%2 == 0 {
		currentPlayerTurn = state.AIPlayerID
	}

	return Snapshot{
		Type:              "snapshot",
		Tick:              state.Tick,
		Units:             unitsCopy,
		Players:           playersCopy,
		Map:               state.Map,
		CurrentPhase:      state.CurrentPhase,
		TurnNumber:        state.TurnNumber,
		HumanPlayerID:     state.HumanPlayerID,
		AIPlayerID:        state.AIPlayerID,
		HumanPlayerReady:  state.HumanPlayerReady,
		AIPlayerReady:     state.AIPlayerReady,
		HumanBaseID:       state.HumanBaseID,
		AIBaseID:          state.AIBaseID,
		Config:            state.Config,
		CurrentPlayerTurn: currentPlayerTurn,
	}
}
