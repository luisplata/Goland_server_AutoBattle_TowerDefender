package network

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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

const playgameDistPath = "frontend/dist"

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
	// Static frontend compiled by Vite, served under /playgame
	playgameFS := http.FileServer(http.Dir(playgameDistPath))
	http.Handle("/playgame/", http.StripPrefix("/playgame/", playgameFS))
	http.HandleFunc("/playgame", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/playgame" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(playgameDistPath, "index.html"))
	})

	http.HandleFunc("/game/create", s.handleCreateGame)
	http.HandleFunc("/game/join", s.handleJoin)
	http.HandleFunc("/game/state", s.handleGameState)
	http.HandleFunc("/command/send", s.handleSendCommand)
	http.HandleFunc("/unit-stats", s.handleUnitStats)
	http.HandleFunc("/ws", s.handleWebSocket)
	http.HandleFunc("/openapi.yml", s.handleOpenAPI)
	http.HandleFunc("/docs", s.handleSwaggerUI)
	http.HandleFunc("/api/docs", s.handleDocIndex)
	http.HandleFunc("/api/readme", s.handleReadme)

	http.ListenAndServe(":7070", nil)
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

// handleDocIndex sirve una página de índice de documentación
func (s *HttpServer) handleDocIndex(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<!DOCTYPE html>
<html>
  <head>
    <title>AutoBattle API Documentation</title>
    <style>
      body { font-family: Arial, sans-serif; background: #1a1a1a; color: #eee; padding: 40px; }
      .container { max-width: 800px; margin: 0 auto; }
      h1 { color: #4CAF50; }
      .doc-link { display: block; padding: 10px; margin: 10px 0; background: #333; border: 1px solid #555; border-radius: 4px; text-decoration: none; color: #4CAF50; transition: 0.2s; }
      .doc-link:hover { background: #444; border-color: #4CAF50; }
      .description { color: #aaa; font-size: 0.9em; margin-top: 5px; }
    </style>
  </head>
  <body>
    <div class="container">
      <h1>🎮 AutoBattle API Documentation</h1>
      <p>Welcome to the AutoBattle server API documentation.</p>
      
      <h2>Documentation Links</h2>
      <a href="/docs" class="doc-link">
        📊 Swagger UI (Interactive API)
        <div class="description">Interactive API explorer with all endpoints and schemas</div>
      </a>
      
      <a href="/openapi.yml" class="doc-link">
        📄 OpenAPI Specification (YAML)
        <div class="description">Raw OpenAPI 3.0 specification file</div>
      </a>
      
      <a href="/api/readme" class="doc-link">
        📖 README
        <div class="description">Project overview and getting started guide</div>
      </a>
      
      <h2>Quick Links</h2>
      <a href="/playgame" class="doc-link">
        🕹️ Play Game
        <div class="description">Launch the game interface</div>
      </a>
      
      <a href="/unit-stats" class="doc-link">
        ⚙️ Unit Statistics (JSON)
        <div class="description">Get all unit types and their stats</div>
      </a>
    </div>
  </body>
</html>`))
}

// handleReadme sirve el README como HTML
func (s *HttpServer) handleReadme(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Leer el archivo README.md
	content, err := ioutil.ReadFile("/app/Readme.md")
	if err != nil {
		// Fallback a ruta relativa para desarrollo local
		content, err = ioutil.ReadFile("Readme.md")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("<h1>README not found</h1>"))
			return
		}
	}

	// Convertir markdown simple a HTML (muy básico)
	html := `<!DOCTYPE html>
<html>
  <head>
    <title>AutoBattle README</title>
    <style>
      body { font-family: Arial, sans-serif; background: #1a1a1a; color: #eee; padding: 40px; line-height: 1.6; }
      .container { max-width: 900px; margin: 0 auto; }
      h1 { color: #4CAF50; border-bottom: 2px solid #4CAF50; padding-bottom: 10px; }
      h2, h3 { color: #4CAF50; margin-top: 20px; }
      code { background: #333; padding: 2px 6px; border-radius: 3px; color: #4CAF50; }
      pre { background: #222; padding: 15px; border-left: 3px solid #4CAF50; overflow-x: auto; }
      a { color: #4CAF50; }
      a:hover { text-decoration: underline; }
      .back-link { margin-bottom: 20px; }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="back-link"><a href="/api/docs">← Back to Documentation</a></div>
      <div class="readme">` + markdownToHTML(string(content)) + `</div>
    </div>
  </body>
</html>`
	w.Write([]byte(html))
}

// markdownToHTML hace una conversión básica de markdown a HTML
func markdownToHTML(md string) string {
	lines := strings.Split(md, "\n")
	var result strings.Builder
	inCode := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Code blocks
		if strings.HasPrefix(trimmed, "```") {
			if inCode {
				result.WriteString("</pre>")
				inCode = false
			} else {
				result.WriteString("<pre><code>")
				inCode = true
			}
			continue
		}

		if inCode {
			result.WriteString(line + "\n")
			continue
		}

		// Headers
		if strings.HasPrefix(line, "# ") {
			result.WriteString("<h1>" + strings.TrimPrefix(line, "# ") + "</h1>\n")
		} else if strings.HasPrefix(line, "## ") {
			result.WriteString("<h2>" + strings.TrimPrefix(line, "## ") + "</h2>\n")
		} else if strings.HasPrefix(line, "### ") {
			result.WriteString("<h3>" + strings.TrimPrefix(line, "### ") + "</h3>\n")
		} else if trimmed != "" {
			result.WriteString("<p>" + trimmed + "</p>\n")
		}
	}

	if inCode {
		result.WriteString("</pre>")
	}

	return result.String()
}
