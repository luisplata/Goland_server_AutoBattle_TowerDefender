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

func (s *HttpServer) Start() {
	http.HandleFunc("/game/create", s.handleCreateGame)
	http.HandleFunc("/game/join", s.handleJoin)
	http.HandleFunc("/game/state", s.handleGameState)
	http.HandleFunc("/command/send", s.handleSendCommand)
	http.HandleFunc("/ws", s.handleWebSocket)

	http.ListenAndServe(":8080", nil)
}

func (s *HttpServer) handleSendCommand(w http.ResponseWriter, r *http.Request) {
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	game := s.manager.CreateGame()
	snapshot := game.State.GetSnapshot()

	response := map[string]interface{}{
		"gameId":   game.ID,
		"snapshot": snapshot,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *HttpServer) handleGameState(w http.ResponseWriter, r *http.Request) {
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
