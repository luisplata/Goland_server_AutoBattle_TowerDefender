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
	g.Units[unit.ID] = unit
	g.nextUnitID++

	return unit
}
