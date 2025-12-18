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
	switch cmd.Type {

	case command.CommandSpawnUnit:
		data, ok := cmd.Data.(map[string]any)
		if !ok {
			slog.Warn("Invalid spawn data")
			return
		}

		slog.Warn("SpawnUnit Command Data", "data", data)

		unitType := data["unitType"].(string)
		x_position := int(data["x"].(float64))
		y_position := int(data["y"].(float64))

		s.spawnUnit(cmd.GameID, cmd.PlayerID, unitType, x_position, y_position)
	}
}

func (s *GameSimulation) spawnUnit(gameId int, playerId int, unitType string, x_position int, y_position int) {
	slog.Info("Spawning unit", "tick", s.state.Tick, "gameId", gameId, "playerId", playerId, "unitType", unitType, "x", x_position, "y", y_position)

	// Crear la unidad en el estado del juego
	unit := s.state.SpawnUnit(playerId, unitType, x_position, y_position)

	if unit == nil {
		slog.Warn("Failed to spawn unit - invalid position", "tick", s.state.Tick, "x", x_position, "y", y_position)
		return
	}

	slog.Info("Unit spawned successfully", "tick", s.state.Tick, "unitId", unit.ID, "playerId", playerId, "x", x_position, "y", y_position)
}

// =======================
// Fases del Tick
// =======================

func (s *GameSimulation) Produce() {}
func (s *GameSimulation) Move()    {}
func (s *GameSimulation) Block()   {}
func (s *GameSimulation) Attack()  {}
func (s *GameSimulation) Cleanup() {}
