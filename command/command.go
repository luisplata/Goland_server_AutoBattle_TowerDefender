package command

type CommandType string

const (
	CommandJoin  CommandType = "join"
	CommandDummy CommandType = "dummy"
)

type Command struct {
	PlayerID int         `json:"playerId"`
	Type     CommandType `json:"type"`
}
