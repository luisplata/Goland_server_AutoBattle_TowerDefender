package command

type SpawnUnitData struct {
	Type CommandType `json:"type"`
	X    int         `json:"x"`
	Y    int         `json:"y"`
}
