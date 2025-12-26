package network

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WsClient struct {
	conn     *websocket.Conn
	gameID   int
	playerID int
}

type WsHub struct {
	mu      sync.Mutex
	clients map[*WsClient]struct{}
}

func NewWsHub() *WsHub {
	return &WsHub{
		clients: make(map[*WsClient]struct{}),
	}
}

func (h *WsHub) Add(client *WsClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client] = struct{}{}
}

func (h *WsHub) Remove(client *WsClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, client)
	client.conn.Close()
}

func (h *WsHub) Broadcast(gameID int, payload any) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for c := range h.clients {
		if c.gameID == gameID {
			_ = c.conn.WriteJSON(payload)
		}
	}
}
