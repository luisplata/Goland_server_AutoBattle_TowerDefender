package game

// UnitCategory define las categorías de unidades/edificios
type UnitCategory string

const (
	CategoryStructure  UnitCategory = "structure"  // Torres, generadores, murallas
	CategoryLandUnit   UnitCategory = "land_unit"  // Unidades terrestres
	CategoryNavalUnit  UnitCategory = "naval_unit" // Unidades navales
	CategoryProjectile UnitCategory = "projectile" // Proyectiles (futuro)
)

// UnitType define los tipos específicos de unidades
const (
	// Base principal
	TypeMainBase = "main_base" // Base principal de cada jugador

	// Estructuras
	TypeTower          = "tower"           // Torre de defensa
	TypeLandGenerator  = "land_generator"  // Generador de unidades terrestres
	TypeNavalGenerator = "naval_generator" // Generador de unidades navales
	TypeWall           = "wall"            // Muralla bloqueadora

	// Unidades terrestres (generadas por land_generator)
	TypeLandSoldier = "land_soldier" // Soldado terrestre básico

	// Unidades navales (generadas por naval_generator)
	TypeNavalShip = "naval_ship" // Barco básico

	// Legacy
	TypeWarrior = "warrior" // Warrior antiguo (deprecated)
)

// UnitStats define las estadísticas base de cada tipo de unidad
type UnitStats struct {
	Category          UnitCategory `json:"category"`
	HP                int          `json:"hp"`
	CanMove           bool         `json:"canMove"`
	MoveIntervalTicks int          `json:"moveIntervalTicks"` // Ticks entre movimientos

	// Combat stats
	AttackDamage        int `json:"attackDamage"`        // Daño por ataque
	AttackRange         int `json:"attackRange"`         // Rango de ataque (en tiles)
	AttackIntervalTicks int `json:"attackIntervalTicks"` // Ticks entre ataques

	// Generator stats
	IsGenerator        bool   `json:"isGenerator"`        // Si genera unidades
	GeneratedUnitType  string `json:"generatedUnitType"`  // Tipo de unidad que genera
	GenerationInterval int    `json:"generationInterval"` // Ticks entre generaciones
	MaxUnitsGenerated  int    `json:"maxUnitsGenerated"`  // Máximo de unidades que puede generar (-1 = infinito)

	// Blocking
	IsBlocker bool `json:"isBlocker"` // Si bloquea el paso
}

// GetUnitStats retorna las estadísticas para un tipo de unidad
func GetUnitStats(unitType string) UnitStats {
	stats := map[string]UnitStats{
		// Base Principal
		TypeMainBase: {
			Category:           CategoryStructure,
			HP:                 1000, // Alta vida
			CanMove:            false,
			AttackDamage:       0, // No ataca
			IsGenerator:        true,
			GeneratedUnitType:  TypeWarrior,
			GenerationInterval: 20, // Genera cada 4 segundos
			MaxUnitsGenerated:  -1, // Infinitas unidades
			IsBlocker:          true,
		},

		// Torres
		TypeTower: {
			Category:            CategoryStructure,
			HP:                  500,
			CanMove:             false,
			AttackDamage:        25,
			AttackRange:         5,  // 5 tiles de rango
			AttackIntervalTicks: 10, // Ataca cada 2 segundos
			IsBlocker:           true,
		},

		// Generador de unidades terrestres
		TypeLandGenerator: {
			Category:           CategoryStructure,
			HP:                 300,
			CanMove:            false,
			IsGenerator:        true,
			GeneratedUnitType:  TypeLandSoldier,
			GenerationInterval: 25, // Genera cada 5 segundos
			MaxUnitsGenerated:  -1, // Infinitas unidades
			IsBlocker:          true,
		},

		// Generador de unidades navales
		TypeNavalGenerator: {
			Category:           CategoryStructure,
			HP:                 300,
			CanMove:            false,
			IsGenerator:        true,
			GeneratedUnitType:  TypeNavalShip,
			GenerationInterval: 30, // Genera cada 6 segundos
			MaxUnitsGenerated:  -1, // Infinitas unidades
			IsBlocker:          true,
		},

		// Muralla
		TypeWall: {
			Category:  CategoryStructure,
			HP:        200,
			CanMove:   false,
			IsBlocker: true,
		},

		// Soldado terrestre
		TypeLandSoldier: {
			Category:            CategoryLandUnit,
			HP:                  100,
			CanMove:             true,
			MoveIntervalTicks:   5, // Se mueve cada segundo
			AttackDamage:        15,
			AttackRange:         1, // Cuerpo a cuerpo
			AttackIntervalTicks: 8, // Ataca cada 1.6 segundos
		},

		// Barco naval
		TypeNavalShip: {
			Category:            CategoryNavalUnit,
			HP:                  150,
			CanMove:             true,
			MoveIntervalTicks:   6, // Más lento que unidades terrestres
			AttackDamage:        20,
			AttackRange:         2,  // Rango medio
			AttackIntervalTicks: 10, // Ataca cada 2 segundos
		},

		// Legacy warrior
		TypeWarrior: {
			Category:            CategoryLandUnit,
			HP:                  100,
			CanMove:             true,
			MoveIntervalTicks:   5,
			AttackDamage:        10,
			AttackRange:         1,
			AttackIntervalTicks: 10,
		},
	}

	if s, ok := stats[unitType]; ok {
		return s
	}

	// Default stats si no se encuentra el tipo
	return UnitStats{
		Category:          CategoryLandUnit,
		HP:                100,
		CanMove:           true,
		MoveIntervalTicks: 5,
	}
}
