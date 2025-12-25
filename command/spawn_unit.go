package command

type SpawnUnitData struct {
	Type     CommandType `json:"type"`
	UnitType string      `json:"unitType"` // "tower", "land_generator", "naval_generator", "wall", etc.
	X        int         `json:"x"`
	Y        int         `json:"y"`
}

// Helper functions para crear comandos de spawn de estructuras

// NewSpawnTowerCommand crea un comando para spawnear una torre
func NewSpawnTowerCommand(playerID, gameID, x, y int) Command {
	return Command{
		PlayerID: playerID,
		GameID:   gameID,
		Type:     CommandSpawnUnit,
		Data: SpawnUnitData{
			Type:     CommandSpawnUnit,
			UnitType: "tower",
			X:        x,
			Y:        y,
		},
	}
}

// NewSpawnLandGeneratorCommand crea un comando para spawnear un generador terrestre
func NewSpawnLandGeneratorCommand(playerID, gameID, x, y int) Command {
	return Command{
		PlayerID: playerID,
		GameID:   gameID,
		Type:     CommandSpawnUnit,
		Data: SpawnUnitData{
			Type:     CommandSpawnUnit,
			UnitType: "land_generator",
			X:        x,
			Y:        y,
		},
	}
}

// NewSpawnNavalGeneratorCommand crea un comando para spawnear un generador naval
func NewSpawnNavalGeneratorCommand(playerID, gameID, x, y int) Command {
	return Command{
		PlayerID: playerID,
		GameID:   gameID,
		Type:     CommandSpawnUnit,
		Data: SpawnUnitData{
			Type:     CommandSpawnUnit,
			UnitType: "naval_generator",
			X:        x,
			Y:        y,
		},
	}
}

// NewSpawnWallCommand crea un comando para spawnear una muralla
func NewSpawnWallCommand(playerID, gameID, x, y int) Command {
	return Command{
		PlayerID: playerID,
		GameID:   gameID,
		Type:     CommandSpawnUnit,
		Data: SpawnUnitData{
			Type:     CommandSpawnUnit,
			UnitType: "wall",
			X:        x,
			Y:        y,
		},
	}
}

// NewSpawnWarriorCommand crea un comando para spawnear un guerrero (legacy, para testing)
func NewSpawnWarriorCommand(playerID, gameID, x, y int) Command {
	return Command{
		PlayerID: playerID,
		GameID:   gameID,
		Type:     CommandSpawnUnit,
		Data: SpawnUnitData{
			Type:     CommandSpawnUnit,
			UnitType: "warrior",
			X:        x,
			Y:        y,
		},
	}
}
