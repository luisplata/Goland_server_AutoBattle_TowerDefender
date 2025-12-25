package command

// EndTurnCommand creates a command to end the current player's turn
func EndTurnCommand(gameID, playerID int) Command {
	return Command{
		GameID:   gameID,
		PlayerID: playerID,
		Type:     CommandEndTurn,
		Data:     nil,
	}
}
