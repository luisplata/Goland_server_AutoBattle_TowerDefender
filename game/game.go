package game

import "autobattle-server/command"

type Game struct {
	ID         int
	State      *GameState
	Simulation *GameSimulation
	Clock      *GameClock
	Commands   *command.CommandQueue
	Snapshot   *Snapshot
}

func NewGame(id int) *Game {
	state := NewGameState()
	simulation := NewGameSimulation(state)

	game := &Game{
		ID:         id,
		State:      state,
		Simulation: simulation,
		Clock:      NewGameClock(200),
		Commands:   command.NewCommandQueue(),
		Snapshot:   nil,
	}

	simulation.BindGame(game)
	return game
}
