package game

type Player struct {
	ID        int      `json:"id"`
	IsAI      bool     `json:"isAi"`
	Deck      []string `json:"-"`         // Oculto en JSON
	Hand      []string `json:"hand"`      // Mano visible para cliente (debug)
	DeckCount int      `json:"deckCount"` // Tamaño del mazo restante (para UI)
	Connected bool     `json:"connected"` // Estado de conexión del jugador
}
