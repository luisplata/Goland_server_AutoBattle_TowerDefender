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

	case command.CommandMoveUnit:
		data, ok := cmd.Data.(map[string]any)
		if !ok {
			slog.Warn("Invalid move data")
			return
		}

		unitID := int(data["unitId"].(float64))
		x := int(data["x"].(float64))
		y := int(data["y"].(float64))

		okDest := s.state.SetUnitDestination(cmd.PlayerID, unitID, x, y)
		if !okDest {
			slog.Warn("SetUnitDestination failed", "tick", s.state.Tick, "playerId", cmd.PlayerID, "unitId", unitID, "x", x, "y", y)
		}
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
func (s *GameSimulation) Move() {
	// Step toward target respecting per-unit move interval
	s.state.mu.Lock()
	defer s.state.mu.Unlock()

	for _, unit := range s.state.Units {
		if !unit.CanMove {
			continue
		}
		if unit.X == unit.TargetX && unit.Y == unit.TargetY {
			continue
		}
		if s.state.Tick < unit.NextMoveTick {
			continue
		}

		// Compute one-tile step toward target (Manhattan), with axis fallback
		newX, newY := unit.X, unit.Y
		dx := unit.TargetX - unit.X
		dy := unit.TargetY - unit.Y

		tryXFirst := abs(dx) >= abs(dy)

		stepTried := false
		if tryXFirst {
			if dx > 0 {
				newX = unit.X + 1
			} else if dx < 0 {
				newX = unit.X - 1
			}
			stepTried = true
			if s.state.isTileAllowedForUnit(unit, newX, newY) {
				unit.X = newX
				unit.Y = newY
				unit.NextMoveTick = s.state.Tick + unit.MoveIntervalTicks
				continue
			}
			// fallback to Y axis
			newX = unit.X
		}

		// Try Y axis (either primary or fallback)
		if dy > 0 {
			newY = unit.Y + 1
		} else if dy < 0 {
			newY = unit.Y - 1
		} else if !stepTried {
			// If X was aligned and Y already equal, nothing to do
			continue
		}

		if s.state.isTileAllowedForUnit(unit, newX, newY) {
			unit.X = newX
			unit.Y = newY
			unit.NextMoveTick = s.state.Tick + unit.MoveIntervalTicks
		} else {
			// Both axes blocked; stop at current
			unit.TargetX = unit.X
			unit.TargetY = unit.Y
		}
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
func (s *GameSimulation) Block()   {}
func (s *GameSimulation) Attack()  {}
func (s *GameSimulation) Cleanup() {}
