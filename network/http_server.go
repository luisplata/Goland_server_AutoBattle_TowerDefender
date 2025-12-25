package network

import (
	"encoding/json"
	"net/http"
	"strconv"

	"autobattle-server/command"
	"autobattle-server/game"

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
	http.HandleFunc("/ws", s.handleWebSocket)

	http.ListenAndServe(":8080", nil)
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

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &WsClient{
		conn:   conn,
		gameID: gameID,
	}

	s.wsHub.Add(client)

	// lectura pasiva (mantiene viva la conexión)
	go func() {
		defer s.wsHub.Remove(client)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
}
