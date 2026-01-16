package network

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"autobattle-server/command"
	"autobattle-server/game"

	"log/slog"

	"github.com/gorilla/websocket"
)

type HttpServer struct {
	manager *game.GameManager
	wsHub   *WsHub
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // dev only
	},
}

func NewHttpServer(manager *game.GameManager, hub *WsHub) *HttpServer {
	return &HttpServer{
		manager: manager,
		wsHub:   hub,
	}
}

// enableCORS adds CORS headers to allow requests from browsers
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func (s *HttpServer) Start() {
	http.HandleFunc("/game/create", s.handleCreateGame)
	http.HandleFunc("/game/join", s.handleJoin)
	http.HandleFunc("/game/state", s.handleGameState)
	http.HandleFunc("/command/send", s.handleSendCommand)
	http.HandleFunc("/unit-stats", s.handleUnitStats)
	http.HandleFunc("/ws", s.handleWebSocket)
	http.HandleFunc("/openapi.yml", s.handleOpenAPI)
	http.HandleFunc("/docs", s.handleSwaggerUI)

	http.ListenAndServe(":8080", nil)
}

func (s *HttpServer) handleOpenAPI(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	http.ServeFile(w, r, "openapi.yml")
}

func (s *HttpServer) handleSwaggerUI(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<!DOCTYPE html>
<html>
  <head>
    <title>Swagger UI</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.onload = () => {
        SwaggerUIBundle({
          url: '/openapi.yml',
          dom_id: '#swagger-ui',
        });
      };
    </script>
  </body>
</html>`))
}

func (s *HttpServer) handleSendCommand(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		GameID   int                 `json:"gameId"`
		PlayerID int                 `json:"playerId"`
		Type     command.CommandType `json:"type"`
		Data     any                 `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	game, ok := s.manager.GetGame(payload.GameID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cmd := command.Command{
		PlayerID: payload.PlayerID,
		Type:     payload.Type,
		Data:     payload.Data,
		GameID:   payload.GameID,
	}

	game.Commands.Enqueue(cmd)
	w.WriteHeader(http.StatusAccepted)
}

func (s *HttpServer) handleCreateGame(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Intentar leer configuración del body (opcional)
	var requestBody struct {
		Config *game.PhaseConfig `json:"config"`
	}

	var createdGame *game.Game

	// Si hay body, intentar parsearlo
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err == nil && requestBody.Config != nil {
			// Crear juego con configuración personalizada
			createdGame = s.manager.CreateGameWithConfig(*requestBody.Config)
		} else {
			// Si hay error o no hay config, usar valores por defecto
			createdGame = s.manager.CreateGame()
		}
	} else {
		// Sin body, usar configuración por defecto
		createdGame = s.manager.CreateGame()
	}

	snapshot := createdGame.State.GetSnapshot()

	response := map[string]interface{}{
		"gameId":   createdGame.ID,
		"snapshot": snapshot,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *HttpServer) handleGameState(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	gameIDStr := r.URL.Query().Get("gameId")
	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	game, ok := s.manager.GetGame(gameID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	snapshot := game.State.GetSnapshot()
	json.NewEncoder(w).Encode(snapshot)
}

func (s *HttpServer) handleJoin(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	gameIDStr := r.URL.Query().Get("gameId")
	gameID, err := strconv.Atoi(gameIDStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	game, ok := s.manager.GetGame(gameID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	player := game.State.AddPlayer()
	json.NewEncoder(w).Encode(player)
}

func (s *HttpServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	gameIDStr := r.URL.Query().Get("gameId")
	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Optional: identify the player for connection tracking
	var playerID int
	if pStr := r.URL.Query().Get("playerId"); pStr != "" {
		if p, convErr := strconv.Atoi(pStr); convErr == nil {
			playerID = p
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &WsClient{
		conn:     conn,
		gameID:   gameID,
		playerID: playerID,
	}

	s.wsHub.Add(client)

	// Mark player connected if identified
	if playerID > 0 {
		if g, ok := s.manager.GetGame(gameID); ok {
			g.State.SetPlayerConnected(playerID, true)
			slog.Info("Player connected", "gameId", gameID, "playerId", playerID)
		}
	}

	// lectura pasiva (mantiene viva la conexión)
	go func() {
		defer s.wsHub.Remove(client)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				// On disconnect mark player as disconnected and start timeout
				if client.playerID > 0 {
					if g, ok := s.manager.GetGame(client.gameID); ok {
						g.State.SetPlayerConnected(client.playerID, false)
						slog.Info("Player disconnected", "gameId", client.gameID, "playerId", client.playerID)

						// Start timeout from config to end game if still offline
						go func(gid, pid int, timeoutSeconds int) {
							time.Sleep(time.Duration(timeoutSeconds) * time.Second)
							if g2, ok2 := s.manager.GetGame(gid); ok2 {
								if !g2.State.IsPlayerConnected(pid) {
									slog.Info("Disconnect timeout reached; ending game", "gameId", gid, "playerId", pid, "timeoutSeconds", timeoutSeconds)
									s.manager.EndGame(gid, pid, "disconnect_timeout")
								}
							}
						}(client.gameID, client.playerID, g.State.Config.DisconnectTimeoutSeconds)
					}
				}
				return
			}
		}
	}()
}

func (s *HttpServer) handleUnitStats(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Retornar estadísticas de todos los tipos de unidades
	unitTypes := []string{
		game.TypeWarrior,
		game.TypeTower,
		game.TypeWall,
		game.TypeLandGenerator,
		game.TypeNavalGenerator,
	}

	stats := make(map[string]game.UnitStats)
	for _, unitType := range unitTypes {
		stats[unitType] = game.GetUnitStats(unitType)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
