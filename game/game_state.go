package game

import "sync"

type GameState struct {
	mu sync.Mutex

	nextPlayerID int
	nextUnitID   int
	Tick         int                `json:"tick"`
	Players      map[int]*Player    `json:"players"`
	Units        map[int]*UnitState `json:"units"`
	Map          *GameMap           `json:"map"`
}

type UnitState struct {
	ID       int    `json:"id"`
	PlayerID int    `json:"playerId"`
	UnitType string `json:"unitType"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	HP       int    `json:"hp"`
	// Server-side movement control (not serialized)
	TargetX           int  `json:"-"`
	TargetY           int  `json:"-"`
	MoveIntervalTicks int  `json:"-"`
	NextMoveTick      int  `json:"-"`
	CanMove           bool `json:"-"`
}

func NewGameState() *GameState {
	return &GameState{
		Players:      make(map[int]*Player),
		nextPlayerID: 1,
		nextUnitID:   1,
		Units:        make(map[int]*UnitState),
		Map:          NewGameMap(),
	}
}

func (g *GameState) AdvanceTick() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.Tick++
}

// SOLO para red (lectura)
func (g *GameState) GetSnapshot() GameState {
	g.mu.Lock()
	defer g.mu.Unlock()

	playersCopy := make(map[int]*Player)
	for id, p := range g.Players {
		playersCopy[id] = &Player{ID: p.ID}
	}

	return GameState{
		Tick:    g.Tick,
		Players: playersCopy,
		Units:   g.Units,
		Map:     g.Map,
	}
}

// SOLO para /join
func (g *GameState) AddPlayer() *Player {
	g.mu.Lock()
	defer g.mu.Unlock()

	player := &Player{
		ID: g.nextPlayerID,
	}
	g.Players[player.ID] = player
	g.nextPlayerID++

	return player
}

// SpawnUnit crea una nueva unidad en el juego
func (g *GameState) SpawnUnit(playerID int, unitType string, x, y int) *UnitState {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Validar posición
	if !g.Map.IsWalkable(x, y) {
		return nil
	}

	unit := &UnitState{
		ID:       g.nextUnitID,
		PlayerID: playerID,
		UnitType: unitType,
		X:        x,
		Y:        y,
		HP:       100,
	}

	// Initialize stats based on unit type
	g.applyUnitStats(unit)

	// Initialize movement target to current position to avoid drifting to (0,0)
	unit.TargetX = unit.X
	unit.TargetY = unit.Y
	g.Units[unit.ID] = unit
	g.nextUnitID++

	return unit
}

// MoveUnit attempts to move a unit to a target position if it's walkable
// and belongs to the requesting player. Returns true on success.
func (g *GameState) MoveUnit(playerID, unitID, x, y int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	unit, ok := g.Units[unitID]
	if !ok {
		return false
	}
	if unit.PlayerID != playerID {
		return false
	}
	// Immediate move (legacy). Prefer SetUnitDestination.
	if !g.isTileAllowedForUnit(unit, x, y) {
		return false
	}

	unit.X = x
	unit.Y = y
	return true
}

// SetUnitDestination sets a target position; unit will step over time.
func (g *GameState) SetUnitDestination(playerID, unitID, x, y int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	unit, ok := g.Units[unitID]
	if !ok {
		return false
	}
	if unit.PlayerID != playerID {
		return false
	}
	if !unit.CanMove {
		return false
	}
	// Destination can be any tile; step validation happens each move tick
	unit.TargetX = x
	unit.TargetY = y
	return true
}

// applyUnitStats assigns movement properties based on UnitType.
func (g *GameState) applyUnitStats(unit *UnitState) {
	switch unit.UnitType {
	case "warrior":
		unit.CanMove = true
		unit.MoveIntervalTicks = 5 // 1 tile per ~1s if tick=200ms
		unit.NextMoveTick = g.Tick
		unit.HP = 100
	case "tower":
		unit.CanMove = false
		unit.MoveIntervalTicks = 0
		unit.NextMoveTick = 0
		unit.HP = 300
	default:
		// default mobile unit
		unit.CanMove = true
		unit.MoveIntervalTicks = 5
		unit.NextMoveTick = g.Tick
	}
}

// isTileAllowedForUnit checks terrain constraints for a given unit type.
func (g *GameState) isTileAllowedForUnit(unit *UnitState, x, y int) bool {
	// For now, warriors/towers use map walkability. Extend for boats later.
	return g.Map.IsWalkable(x, y)
}
