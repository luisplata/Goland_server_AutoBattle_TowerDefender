package game

// HandUpdateEvent se envía cuando la mano de un jugador cambia
type HandUpdateEvent struct {
	Type      string   `json:"type"` // "hand_updated"
	PlayerID  int      `json:"playerId"`
	Hand      []string `json:"hand"`
	DeckCount int      `json:"deckCount"`
}

// BuildHandUpdateEvent crea un evento de actualización de mano
func BuildHandUpdateEvent(playerID int, hand []string, deckCount int) HandUpdateEvent {
	return HandUpdateEvent{
		Type:      "hand_updated",
		PlayerID:  playerID,
		Hand:      append([]string{}, hand...), // copia
		DeckCount: deckCount,
	}
}
