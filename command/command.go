package command

type CommandType string

const (
	CommandDummy     CommandType = "dummy"
	CommandSpawnUnit CommandType = "spawn_unit"
	CommandUpgrade   CommandType = "upgrade"
	CommandMoveUnit  CommandType = "move_unit"
	CommandEndTurn   CommandType = "end_turn" // Deprecated - usar ready
	CommandReady     CommandType = "ready"    // Jugador listo para pasar de fase
)

type Command struct {
	PlayerID int         `json:"playerId"`
	Type     CommandType `json:"type"`
	Data     any         `json:"data"`
	GameID   int         `json:"gameId"`
}
