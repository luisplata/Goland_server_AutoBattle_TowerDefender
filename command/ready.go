package command

// ReadyCommand creates a command to mark player as ready
func ReadyCommand(gameID, playerID int) Command {
	return Command{
		GameID:   gameID,
		PlayerID: playerID,
		Type:     CommandReady,
		Data:     nil,
	}
}
