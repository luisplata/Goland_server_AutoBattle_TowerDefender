package main

import (
	"autobattle-server/game"
	"autobattle-server/network"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"reflect"
	"time"
)

func main() {
	go func() {
		// pprof en localhost:6060
		http.ListenAndServe("localhost:6060", nil)
	}()
	if err := InitDB(); err != nil {
		slog.Error("No se pudo conectar a la base de datos", "error", err)
		return
	}

	gameManager := game.NewGameManager()
	wsHub := network.NewWsHub()

	httpServer := network.NewHttpServer(gameManager, wsHub)
	go httpServer.Start()

	lastSnapshots := make(map[int]*game.Snapshot)

	ticker := time.NewTicker(33 * time.Millisecond) // 30 FPS
	defer ticker.Stop()
	for {
		games := gameManager.GetAllGames()

		if len(games) == 0 {
			// Si no hay juegos activos, espera hasta el siguiente tick
			<-ticker.C
			continue
		}

		for _, g := range games {
			if g.Clock.ShouldTick() {
				// Si el fin de juego ya fue confirmado, terminar inmediatamente
				if g.State.GameEnd != nil && g.State.GameEnd.Confirmed {
					currentSnapshot := game.BuildSnapshot(g.State)
					if last, ok := lastSnapshots[g.ID]; !ok || !reflect.DeepEqual(*last, currentSnapshot) {
						wsHub.Broadcast(g.ID, game.SnapshotToUpdate(currentSnapshot))
						lastSnapshots[g.ID] = &currentSnapshot
					}
					gameManager.EndGame(g.ID, g.State.GameEnd.LoserID, g.State.GameEnd.Reason)
					// Limpiar memoria de snapshots y otros recursos del juego terminado
					delete(lastSnapshots, g.ID)
					continue
				}
				// Si hay fin de juego pendiente, no avanzar simulación; solo emitir snapshot si cambió
				if g.State.IsGameEndPending() {
					currentSnapshot := game.BuildSnapshot(g.State)
					if last, ok := lastSnapshots[g.ID]; !ok || !reflect.DeepEqual(*last, currentSnapshot) {
						wsHub.Broadcast(g.ID, game.SnapshotToUpdate(currentSnapshot))
						lastSnapshots[g.ID] = &currentSnapshot
					}
					continue
				}

				// Guardar la fase anterior antes de procesar
				previousPhase := g.State.GetCurrentPhase()

				g.Simulation.ProcessTick()

				// Verificar condiciones de victoria/derrota
				if gameOver, loserID, reason := g.Simulation.CheckVictoryConditions(); gameOver {
					slog.Info("Game ended - victory condition met (pending confirmation)", "gameId", g.ID, "loserId", loserID, "reason", reason)
					g.State.SetPendingEnd(loserID, reason)
					currentSnapshot := game.BuildSnapshot(g.State)
					if last, ok := lastSnapshots[g.ID]; !ok || !reflect.DeepEqual(*last, currentSnapshot) {
						wsHub.Broadcast(g.ID, game.SnapshotToUpdate(currentSnapshot))
						lastSnapshots[g.ID] = &currentSnapshot
					}
					continue // Skip further processing for this game
				}

				currentSnapshot := game.BuildSnapshot(g.State)

				// Detectar cambio de fase y enviar evento especial
				if g.State.DidPhaseChange() {
					phaseEvent := game.BuildPhaseChangeEvent(g.State, previousPhase)
					wsHub.Broadcast(g.ID, phaseEvent)
				}

				// Detectar cambios en manos y enviar eventos hand_updated
				updatedPlayers := g.State.DrainHandUpdates()
				if len(updatedPlayers) > 0 {
					for _, pID := range updatedPlayers {
						if player, ok := currentSnapshot.Players[pID]; ok {
							handEvent := game.BuildHandUpdateEvent(pID, player.Hand, player.DeckCount)
							wsHub.Broadcast(g.ID, handEvent)
						}
					}
				}

				// Solo enviar snapshot si cambió respecto al anterior
				if last, ok := lastSnapshots[g.ID]; !ok || !reflect.DeepEqual(*last, currentSnapshot) {
					wsHub.Broadcast(g.ID, game.SnapshotToUpdate(currentSnapshot))
					lastSnapshots[g.ID] = &currentSnapshot
				}
			}
		}

		// Esperar al siguiente tick (30 FPS)
		<-ticker.C
	}
}
