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

		unitType := data["unitType"].(string)
		lane := int(data["lane"].(float64))

		s.spawnUnit(cmd.GameID, cmd.PlayerID, unitType, lane)
	}
}

func (s *GameSimulation) spawnUnit(gameId int, playerId int, unitType string, lane int) {

}

// =======================
// Fases del Tick
// =======================

func (s *GameSimulation) Produce() {}
func (s *GameSimulation) Move()    {}
func (s *GameSimulation) Block()   {}
func (s *GameSimulation) Attack()  {}
func (s *GameSimulation) Cleanup() {}
