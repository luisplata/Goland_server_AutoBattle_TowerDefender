package game

import (
	"log/slog"
	"sync"
)

type GameManager struct {
	mu     sync.Mutex
	games  map[int]*Game
	nextID int
}

func NewGameManager() *GameManager {
	return &GameManager{
		games:  make(map[int]*Game),
		nextID: 1,
	}
}

func (gm *GameManager) CreateGame() *Game {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game := NewGame(gm.nextID)
	gm.games[game.ID] = game
	gm.nextID++

	return game
}

// CreateGameWithConfig crea un juego con configuración personalizada
func (gm *GameManager) CreateGameWithConfig(config PhaseConfig) *Game {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game := NewGameWithConfig(gm.nextID, config)
	gm.games[game.ID] = game
	gm.nextID++

	return game
}

func (gm *GameManager) GetGame(id int) (*Game, bool) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game, ok := gm.games[id]
	return game, ok
}

func (gm *GameManager) GetAllGames() []*Game {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	list := make([]*Game, 0, len(gm.games))
	for _, g := range gm.games {
		list = append(list, g)
	}
	return list
}

// EndGame elimina el juego y registra el motivo/derrota
func (gm *GameManager) EndGame(id int, loserID int, reason string) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	if g, ok := gm.games[id]; ok {
		slog.Info("Ending game due to condition", "gameId", id, "loserId", loserID, "reason", reason, "tick", g.State.Tick, "turn", g.State.TurnNumber)
		delete(gm.games, id)
	}
}
