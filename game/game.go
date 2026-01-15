package game

import "autobattle-server/command"

type Game struct {
	ID         int
	State      *GameState
	Simulation *GameSimulation
	Clock      *GameClock
	Commands   *command.CommandQueue
	Snapshot   *Snapshot
	Delta      *Delta
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
		Delta:      nil,
	}

	// Set ticks-per-second into state for DPS-to-ticks calculations
	state.TicksPerSecond = game.Clock.TicksPerSecond()
	simulation.BindGame(game)
	return game
}

// NewGameWithConfig crea un nuevo juego con configuración personalizada
func NewGameWithConfig(id int, config PhaseConfig) *Game {
	state := NewGameStateWithConfig(config)
	simulation := NewGameSimulation(state)

	game := &Game{
		ID:         id,
		State:      state,
		Simulation: simulation,
		Clock:      NewGameClock(200),
		Commands:   command.NewCommandQueue(),
		Snapshot:   nil,
		Delta:      nil,
	}

	// Set ticks-per-second into state for DPS-to-ticks calculations
	state.TicksPerSecond = game.Clock.TicksPerSecond()
	simulation.BindGame(game)
	return game
}
