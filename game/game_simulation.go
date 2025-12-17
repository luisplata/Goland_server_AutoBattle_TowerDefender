package game

import (
	"autobattle-server/command"
	"log/slog"
)

type GameSimulation struct {
	state *GameState
	game  *Game
}

func NewGameSimulation(state *GameState) *GameSimulation {
	return &GameSimulation{
		state: state,
	}
}

func (s *GameSimulation) BindGame(game *Game) {
	s.game = game
}

func (s *GameSimulation) ProcessTick() {
	s.state.AdvanceTick()

	// 1️⃣ Aplicar comandos del tick
	commands := s.game.Commands.Drain()
	for _, cmd := range commands {
		s.ApplyCommand(cmd)
	}

	// 2️⃣ Lógica del juego (orden SAGRADO)
	s.Produce()
	s.Move()
	s.Block()
	s.Attack()
	s.Cleanup()
}

// =======================
// Comandos
// =======================

func (s *GameSimulation) ApplyCommand(cmd command.Command) {
	slog.Info("Applying command", "cmd", cmd, "playerId", cmd.PlayerID)
	switch cmd.Type {
	case command.CommandDummy:
		// placeholder
	}
}

// =======================
// Fases del Tick
// =======================

func (s *GameSimulation) Produce() {}
func (s *GameSimulation) Move()    {}
func (s *GameSimulation) Block()   {}
func (s *GameSimulation) Attack()  {}
func (s *GameSimulation) Cleanup() {}
