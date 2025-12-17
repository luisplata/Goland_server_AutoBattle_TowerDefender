package game

import "sync"

type GameState struct {
	mu sync.Mutex

	Tick    int             `json:"tick"`
	Players map[int]*Player `json:"players"`

	nextPlayerID int
	Units        map[int]*UnitState `json:"units"`
}

type UnitState struct {
	ID int `json:"id"`
	X  int `json:"x"`
	Y  int `json:"y"`
	HP int `json:"hp"`
}

func NewGameState() *GameState {
	return &GameState{
		Players:      make(map[int]*Player),
		nextPlayerID: 1,
		Units:        make(map[int]*UnitState),
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
