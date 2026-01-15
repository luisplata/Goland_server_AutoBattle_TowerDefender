package game

import (
	"autobattle-server/command"
	"log/slog"
	"math/rand"
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

	// 1.5️⃣ Procesar fases del juego
	s.ProcessPhases()

	// 2️⃣ Lógica del juego (solo en fase de batalla)
	if s.state.GetCurrentPhase() == PhaseBattle {
		s.Produce()
		s.Move()
		s.Block()
		s.Attack()
		s.Cleanup()
	}
}

// =======================
// Comandos
// =======================

func (s *GameSimulation) ApplyCommand(cmd command.Command) {
	// Validar que el jugador puede actuar en la fase actual
	// PlaceBase solo se permite en base_selection, otros comandos en preparation
	if cmd.Type != command.CommandReady && cmd.Type != command.CommandPlaceBase && !s.state.CanPlayerAct(cmd.PlayerID) {
		slog.Warn("Command rejected: not in preparation phase", "playerId", cmd.PlayerID, "commandType", cmd.Type, "currentPhase", s.state.GetCurrentPhase())
		return
	}

	switch cmd.Type {

	case command.CommandPlaceBase:
		data, ok := cmd.Data.(map[string]any)
		if !ok {
			slog.Warn("Invalid place_base data")
			return
		}

		x := int(data["x"].(float64))
		y := int(data["y"].(float64))

		// Solo permitir colocar base en fase base_selection
		if s.state.GetCurrentPhase() != PhaseBaseSelection {
			slog.Warn("Cannot place base outside base_selection phase", "playerId", cmd.PlayerID)
			return
		}

		// Verificar que no haya colocado base ya
		if s.state.HasPlayerPlacedBase(cmd.PlayerID) {
			slog.Warn("Player already placed base", "playerId", cmd.PlayerID)
			return
		}

		// Colocar base
		base := s.state.SpawnUnit(cmd.PlayerID, TypeMainBase, x, y)
		if base == nil {
			slog.Warn("Failed to place base", "playerId", cmd.PlayerID, "x", x, "y", y)
			return
		}

		s.state.MarkBasePlaced(cmd.PlayerID, base.ID)
		slog.Info("Base placed", "playerId", cmd.PlayerID, "baseId", base.ID, "x", x, "y", y)

	case command.CommandSpawnUnit:
		data, ok := cmd.Data.(map[string]any)
		if !ok {
			slog.Warn("Invalid spawn data")
			return
		}

		slog.Info("SpawnUnit Command", "data", data)

		unitType := data["unitType"].(string)
		x_position := int(data["x"].(float64))
		y_position := int(data["y"].(float64))

		if !s.state.ConsumeCardFromHand(cmd.PlayerID, unitType) {
			slog.Warn("Spawn rejected: card not in hand", "playerId", cmd.PlayerID, "unitType", unitType)
			return
		}

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

	case command.CommandReady:
		slog.Info("Player ready", "playerId", cmd.PlayerID, "tick", s.state.Tick, "phase", s.state.GetCurrentPhase())
		s.state.SetPlayerReady(cmd.PlayerID, true)

	case command.CommandEndTurn:
		// Backward compatibility - tratar como ready
		slog.Info("Player ready (via end_turn)", "playerId", cmd.PlayerID, "tick", s.state.Tick)
		s.state.SetPlayerReady(cmd.PlayerID, true)
	}
}

// ProcessPhases maneja la transición automática entre fases
func (s *GameSimulation) ProcessPhases() {
	currentPhase := s.state.GetCurrentPhase()

	s.state.mu.Lock()
	ticksSincePhaseStart := s.state.Tick - s.state.PhaseStartTick
	config := s.state.Config
	s.state.mu.Unlock()

	switch currentPhase {
	case PhaseBaseSelection:
		// Fase de selección de base: esperar a que ambos jugadores coloquen su base
		if s.state.BothBasesPlaced() {
			slog.Info("Both bases placed, advancing to TurnStart", "tick", s.state.Tick)
			s.state.StartFirstTurn() // Iniciar el turno 1
			s.state.AdvancePhase()
		} else {
			// La IA coloca su base SOLO después de que el humano coloque la suya
			humanPlaced := s.state.HasPlayerPlacedBase(s.state.HumanPlayerID)
			aiPlaced := s.state.HasPlayerPlacedBase(s.state.AIPlayerID)

			if humanPlaced && !aiPlaced {
				s.placeAIBase()
			}
		}

	case PhaseTurnStart:
		// Fase de inicio: usa duración configurada
		if ticksSincePhaseStart >= config.TurnStartDuration {
			slog.Info("Advancing from TurnStart to Preparation", "tick", s.state.Tick, "turn", s.state.TurnNumber)
			s.state.AdvancePhase()
		}

	case PhasePreparation:
		// Fase de preparación: avanzar cuando ambos jugadores estén listos o timeout
		if s.state.AreBothPlayersReady() {
			slog.Info("Both players ready, advancing to Battle", "tick", s.state.Tick)
			s.state.AdvancePhase()
		} else if ticksSincePhaseStart >= config.PreparationDuration {
			slog.Info("Preparation timeout, advancing to Battle", "tick", s.state.Tick)
			s.state.AdvancePhase()
		} else {
			// La IA se marca como lista automáticamente después de algunos ticks
			s.ProcessAIPreparation(ticksSincePhaseStart)
		}

	case PhaseBattle:
		// Fase de batalla: usa duración configurada
		if ticksSincePhaseStart >= config.BattleDuration {
			slog.Info("Battle finished, advancing to TurnEnd", "tick", s.state.Tick)
			s.state.AdvancePhase()
		}

	case PhaseTurnEnd:
		// Fase de fin: usa duración configurada
		if ticksSincePhaseStart >= config.TurnEndDuration {
			slog.Info("Advancing from TurnEnd to TurnStart", "tick", s.state.Tick)
			s.state.AdvancePhase()
		}
	}
}

// placeAIBase coloca la base de la IA en una posición válida automáticamente
func (s *GameSimulation) placeAIBase() {
	s.state.mu.Lock()
	aiID := s.state.AIPlayerID
	s.state.mu.Unlock()

	// Buscar una posición válida aleatoria en el mapa
	mapWidth := s.state.Map.Width
	mapHeight := s.state.Map.Height

	// Intentar hasta 100 posiciones aleatorias
	for attempts := 0; attempts < 100; attempts++ {
		x := rand.Intn(mapWidth)
		y := rand.Intn(mapHeight)

		// Verificar que la posición sea válida
		if x >= 0 && x < mapWidth && y >= 0 && y < mapHeight {
			if s.state.canUnitTypeEnter(TypeMainBase, -1, x, y) {
				base := s.state.SpawnUnit(aiID, TypeMainBase, x, y)
				if base != nil {
					s.state.MarkBasePlaced(aiID, base.ID)
					slog.Info("AI base placed randomly", "aiId", aiID, "baseId", base.ID, "x", x, "y", y)
					return
				}
			}
		}
	}

	slog.Warn("AI failed to find valid position for base after 100 attempts")
}

// ProcessAIPreparation maneja la lógica de la IA en fase de preparación
func (s *GameSimulation) ProcessAIPreparation(ticksSinceStart int) {
	s.state.mu.Lock()
	aiReadyDelay := s.state.Config.AIReadyDelay
	aiPlayerID := s.state.AIPlayerID
	s.state.mu.Unlock()

	// La IA se marca como lista después del delay configurado
	if ticksSinceStart >= aiReadyDelay {
		s.playAICard()
		s.state.SetPlayerReady(aiPlayerID, true)
	}
}

// playAICard intenta jugar la primera carta en mano en una posición válida.
func (s *GameSimulation) playAICard() {
	// Obtener la primera carta de la mano del AI
	s.state.mu.Lock()
	aiID := s.state.AIPlayerID
	p, ok := s.state.Players[aiID]
	if !ok || len(p.Hand) == 0 {
		s.state.mu.Unlock()
		return
	}
	card := p.Hand[0]
	s.state.mu.Unlock()

	// Buscar posición válida
	x, y, okPos := s.state.findSpawnPosition(card, 50)
	if !okPos {
		slog.Warn("AI could not find spawn position for card", "card", card)
		return
	}

	// Consumir carta y spawnear
	if !s.state.ConsumeCardFromHand(aiID, card) {
		slog.Warn("AI failed to consume card", "card", card)
		return
	}
	s.spawnUnit(0, aiID, card, x, y)
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

func (s *GameSimulation) Produce() {
	// Recolectar spawns para ejecutarlos fuera del lock
	spawns := make([]struct {
		playerID int
		unitType string
		x        int
		y        int
		genID    int
	}, 0)

	s.state.mu.Lock()
	currentTick := s.state.Tick

	for _, unit := range s.state.Units {
		if !unit.IsGenerator {
			continue
		}
		if unit.GenerationInterval <= 0 {
			continue
		}
		if unit.MaxUnitsGenerated >= 0 && unit.UnitsGenerated >= unit.MaxUnitsGenerated {
			continue
		}
		if currentTick < unit.NextGenerationTick {
			continue
		}

		// Buscar un tile adyacente para spawnear
		candidates := [][2]int{{unit.X + 1, unit.Y}, {unit.X - 1, unit.Y}, {unit.X, unit.Y + 1}, {unit.X, unit.Y - 1}}
		spawned := false
		for _, pos := range candidates {
			x, y := pos[0], pos[1]
			if s.state.canUnitTypeEnter(unit.GeneratedUnitType, -1, x, y) {
				spawns = append(spawns, struct {
					playerID int
					unitType string
					x        int
					y        int
					genID    int
				}{playerID: unit.PlayerID, unitType: unit.GeneratedUnitType, x: x, y: y, genID: unit.ID})
				unit.UnitsGenerated++
				spawned = true
				break
			}
		}

		// Programar próximo intento (aunque no haya espacio) para evitar spam por tick
		unit.NextGenerationTick = currentTick + unit.GenerationInterval

		if !spawned {
			slog.Warn("Generator had no space to spawn", "tick", currentTick, "generatorId", unit.ID)
		}
	}

	s.state.mu.Unlock()

	// Ejecutar spawns fuera del lock principal
	for _, job := range spawns {
		spawnedUnit := s.state.SpawnUnit(job.playerID, job.unitType, job.x, job.y)
		if spawnedUnit != nil {
			slog.Info("Generator spawned unit", "tick", currentTick, "generatorId", job.genID, "unitId", spawnedUnit.ID, "type", job.unitType, "x", job.x, "y", job.y)
		} else {
			slog.Warn("Generator failed to spawn unit", "tick", currentTick, "generatorId", job.genID, "type", job.unitType, "x", job.x, "y", job.y)
		}
	}
}
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
func (s *GameSimulation) Block() {}

// Attack procesa ataques automáticos para unidades con daño y rango.
func (s *GameSimulation) Attack() {
	s.state.mu.Lock()
	currentTick := s.state.Tick

	for _, attacker := range s.state.Units {
		if attacker.AttackDamage <= 0 {
			continue
		}
		if currentTick < attacker.NextAttackTick {
			continue
		}

		var target *UnitState
		bestDist := 1_000_000
		for _, candidate := range s.state.Units {
			if candidate.PlayerID == attacker.PlayerID {
				continue
			}
			dx := abs(attacker.X - candidate.X)
			dy := abs(attacker.Y - candidate.Y)
			dist := dx + dy // Manhattan
			if dist <= attacker.AttackRange && dist < bestDist {
				bestDist = dist
				target = candidate
			}
		}

		attacker.NextAttackTick = currentTick + attacker.AttackIntervalTicks
		if target == nil {
			continue
		}

		target.HP -= attacker.AttackDamage
		slog.Info("Attack", "tick", currentTick, "attackerId", attacker.ID, "targetId", target.ID, "damage", attacker.AttackDamage, "targetHP", target.HP)
	}

	s.state.mu.Unlock()
}

// Cleanup elimina unidades con HP <= 0.
func (s *GameSimulation) Cleanup() {
	s.state.mu.Lock()
	dead := make([]int, 0)
	for id, unit := range s.state.Units {
		if unit.HP <= 0 {
			dead = append(dead, id)
		}
	}
	for _, id := range dead {
		slog.Info("Removing dead unit", "unitId", id)
		delete(s.state.Units, id)
	}
	s.state.mu.Unlock()
}
