package game

import "sync"

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
