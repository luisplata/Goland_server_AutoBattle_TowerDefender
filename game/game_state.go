package game

import (
	"math/rand"
	"sync"
	"time"
)

// GamePhase representa las diferentes fases del juego
type GamePhase string

const (
	PhaseBaseSelection GamePhase = "base_selection" // Selección de base al inicio
	PhaseTurnStart     GamePhase = "turn_start"     // Inicio de turno (efectos gráficos)
	PhasePreparation   GamePhase = "preparation"    // Preparación (spawn unidades/edificios)
	PhaseBattle        GamePhase = "battle"         // Batalla automática
	PhaseTurnEnd       GamePhase = "turn_end"       // Fin de turno
)

// PhaseConfig define la duración de cada fase en ticks
type PhaseConfig struct {
	TurnStartDuration        int `json:"turnStartDuration"`        // Ticks para fase turn_start
	PreparationDuration      int `json:"preparationDuration"`      // Ticks máximo para preparation (o hasta que estén listos)
	BattleDuration           int `json:"battleDuration"`           // Ticks para batalla (placeholder)
	TurnEndDuration          int `json:"turnEndDuration"`          // Ticks para turn_end
	AIReadyDelay             int `json:"aiReadyDelay"`             // Ticks que espera la IA para marcarse lista
	DisconnectTimeoutSeconds int `json:"disconnectTimeoutSeconds"` // Segundos antes de terminar juego por desconexión
	CardsPerTurn             int `json:"cardsPerTurn"`             // Cantidad de cartas a robar al inicio de cada turno
	InitialCardsPerHand      int `json:"initialCardsPerHand"`      // Cantidad de cartas iniciales en la mano
}

// DefaultPhaseConfig retorna la configuración por defecto
func DefaultPhaseConfig() PhaseConfig {
	return PhaseConfig{
		TurnStartDuration:        15,  // ~3 segundos
		PreparationDuration:      150, // ~30 segundos
		BattleDuration:           300, // ~60 segundos
		TurnEndDuration:          15,  // ~3 segundos
		AIReadyDelay:             5,   // ~1 segundo
		DisconnectTimeoutSeconds: 30,  // 30 segundos de timeout
		CardsPerTurn:             1,   // 1 carta al inicio de cada turno
		InitialCardsPerHand:      3,   // 3 cartas iniciales en la mano
	}
}

type GameState struct {
	mu sync.Mutex

	nextPlayerID int
	nextUnitID   int
	Tick         int                `json:"tick"`
	Players      map[int]*Player    `json:"players"`
	Units        map[int]*UnitState `json:"units"`
	Map          *GameMap           `json:"map"`

	// Phase-based system
	CurrentPhase         GamePhase   `json:"currentPhase"`  // Fase actual del juego
	TurnNumber           int         `json:"turnNumber"`    // Número de turno actual
	PhaseStartTick       int         `json:"-"`             // Tick en el que empezó la fase actual
	PhaseChangedThisTick bool        `json:"-"`             // Flag para indicar si la fase cambió este tick
	AIPlayerID           int         `json:"aiPlayerId"`    // ID del jugador AI
	HumanPlayerID        int         `json:"humanPlayerId"` // ID del jugador humano
	Config               PhaseConfig `json:"config"`        // Configuración de duración de fases

	// Preparation phase flags
	HumanPlayerReady bool `json:"humanPlayerReady"` // Si el jugador humano está listo
	AIPlayerReady    bool `json:"aiPlayerReady"`    // Si la IA está lista

	// Base IDs - para verificar que ambas bases existen
	HumanBaseID int `json:"humanBaseId"` // ID de la unidad base del jugador humano
	AIBaseID    int `json:"aiBaseId"`    // ID de la unidad base de la IA

	// Hand tracking
	HandUpdatedPlayers []int `json:"-"` // IDs de jugadores cuya mano cambió este tick

	// Timing config
	TicksPerSecond int `json:"-"` // Para convertir DPS en ticks
}

// defaultDeck devuelve un mazo básico con todas las cartas disponibles.
// Cada tipo de carta tiene múltiples copias para asegurar disponibilidad.
func defaultDeck() []string {
	deck := []string{}
	// 15 copias de cada tipo de estructura
	for i := 0; i < 15; i++ {
		deck = append(deck, TypeTower)
		deck = append(deck, TypeLandGenerator)
		deck = append(deck, TypeNavalGenerator)
	}
	// 20 copias de muros (más comunes)
	for i := 0; i < 20; i++ {
		deck = append(deck, TypeWall)
	}
	// 30 copias de warriors (unidad básica)
	for i := 0; i < 30; i++ {
		deck = append(deck, TypeWarrior)
	}
	return deck
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func shuffleCards(cards []string) {
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
}

type UnitState struct {
	ID       int    `json:"id"`
	PlayerID int    `json:"playerId"`
	UnitType string `json:"unitType"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	HP       int    `json:"hp"`
	MaxHP    int    `json:"maxHp"`

	// Combat properties
	AttackDamage        int     `json:"attackDamage"`
	AttackRange         int     `json:"attackRange"`
	AttackIntervalTicks int     `json:"-"`
	NextAttackTick      int     `json:"-"`
	AttackDPS           float64 `json:"attackDps"`

	// Movement control (not serialized)
	TargetX           int  `json:"-"`
	TargetY           int  `json:"-"`
	MoveIntervalTicks int  `json:"-"`
	NextMoveTick      int  `json:"-"`
	CanMove           bool `json:"-"`
	BlockedTicks      int  `json:"-"` // Contador de ticks bloqueado (para detectar deadlocks)

	// Detection
	DetectionRange int `json:"detectionRange"`

	// Activity/animation state
	Status string `json:"status"` // idle, moving, waiting, blocked, attacking

	// Generator properties
	IsGenerator        bool   `json:"isGenerator"`
	GeneratedUnitType  string `json:"generatedUnitType,omitempty"`
	GenerationInterval int    `json:"-"`
	NextGenerationTick int    `json:"-"`
	UnitsGenerated     int    `json:"-"`
	MaxUnitsGenerated  int    `json:"-"`

	// Blocking
	IsBlocker bool `json:"isBlocker"`

	// Category for pathfinding
	Category UnitCategory `json:"category"`

	// Build Range - área que esta estructura expande
	BuildRange int `json:"buildRange"` // Radio de construcción que proporciona
}

func NewGameState() *GameState {
	// Usar tiempo actual como seed por defecto
	return NewGameStateWithSeed(time.Now().UnixNano())
}

func NewGameStateWithSeed(seed int64) *GameState {
	return &GameState{
		Players:      make(map[int]*Player),
		nextPlayerID: 1,
		nextUnitID:   1,
		Units:        make(map[int]*UnitState),
		Map:          NewGameMap(seed),
		CurrentPhase: PhaseBaseSelection,   // Empezar en fase de selección de base
		TurnNumber:   0,                    // El turno 1 empieza después de colocar bases
		Config:       DefaultPhaseConfig(), // Usar configuración por defecto
	}
}

// NewGameStateWithConfig crea un nuevo estado de juego con configuración personalizada
func NewGameStateWithConfig(config PhaseConfig) *GameState {
	state := NewGameState()
	state.Config = config
	return state
}

func (g *GameState) AdvanceTick() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.Tick++
}

// SOLO para red (lectura)
func (g *GameState) GetSnapshot() GameState {
	g.mu.Lock()
	defer g.mu.Unlock()

	playersCopy := make(map[int]*Player)
	for id, p := range g.Players {
		playersCopy[id] = &Player{
			ID:        p.ID,
			IsAI:      p.IsAI,
			Hand:      append([]string{}, p.Hand...),
			DeckCount: p.DeckCount,
			Connected: p.Connected,
		}
	}

	return GameState{
		Tick:             g.Tick,
		Players:          playersCopy,
		Units:            g.Units,
		Map:              g.Map,
		CurrentPhase:     g.CurrentPhase,
		TurnNumber:       g.TurnNumber,
		AIPlayerID:       g.AIPlayerID,
		HumanPlayerID:    g.HumanPlayerID,
		HumanPlayerReady: g.HumanPlayerReady,
		AIPlayerReady:    g.AIPlayerReady,
		Config:           g.Config,
	}
}

// SOLO para /join
func (g *GameState) AddPlayer() *Player {
	g.mu.Lock()
	defer g.mu.Unlock()

	player := &Player{
		ID: g.nextPlayerID,
	}
	player.Deck = defaultDeck()
	shuffleCards(player.Deck)
	player.DeckCount = len(player.Deck)

	// Dibujar mano inicial (cantidad según config)
	for i := 0; i < g.Config.InitialCardsPerHand && len(player.Deck) > 0; i++ {
		if card, ok := g.drawCardLocked(player); ok {
			player.Hand = append(player.Hand, card)
		}
	}

	g.Players[player.ID] = player

	// Si es el primer jugador (humano), crear también el jugador AI
	if g.nextPlayerID == 1 {
		g.HumanPlayerID = player.ID

		// Crear jugador AI
		g.nextPlayerID++
		aiPlayer := &Player{
			ID:   g.nextPlayerID,
			IsAI: true,
		}
		aiPlayer.Deck = defaultDeck()
		shuffleCards(aiPlayer.Deck)
		aiPlayer.DeckCount = len(aiPlayer.Deck)

		// Dibujar mano inicial para IA (cantidad según config)
		for i := 0; i < g.Config.InitialCardsPerHand && len(aiPlayer.Deck) > 0; i++ {
			if card, ok := g.drawCardLocked(aiPlayer); ok {
				aiPlayer.Hand = append(aiPlayer.Hand, card)
			}
		}

		g.Players[aiPlayer.ID] = aiPlayer
		g.AIPlayerID = aiPlayer.ID
	}

	g.nextPlayerID++
	return player
}

// SetPlayerConnected marca el estado de conexión de un jugador
func (g *GameState) SetPlayerConnected(playerID int, connected bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if p, ok := g.Players[playerID]; ok {
		p.Connected = connected
	}
}

// IsPlayerConnected retorna true si el jugador está marcado como conectado
func (g *GameState) IsPlayerConnected(playerID int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if p, ok := g.Players[playerID]; ok {
		return p.Connected
	}
	return false
}

// drawCardLocked roba una carta del mazo del jugador (requiere lock tomado).
// Si el mazo se vacía, se recrea y baraja automáticamente (mazo infinito).
func (g *GameState) drawCardLocked(p *Player) (string, bool) {
	if len(p.Deck) == 0 {
		// Recrear y barajar el mazo cuando se acabe
		p.Deck = defaultDeck()
		shuffleCards(p.Deck)
	}
	card := p.Deck[0]
	p.Deck = p.Deck[1:]
	p.Hand = append(p.Hand, card)
	p.DeckCount = len(p.Deck)
	return card, true
}

// DrawCard roba una carta para el jugador.
func (g *GameState) DrawCard(playerID int) (string, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	p, ok := g.Players[playerID]
	if !ok {
		return "", false
	}
	return g.drawCardLocked(p)
}

// HasCardInHand verifica si un jugador tiene una carta específica en su mano.
func (g *GameState) HasCardInHand(playerID int, unitType string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	p, ok := g.Players[playerID]
	if !ok {
		return false
	}
	for _, c := range p.Hand {
		if c == unitType {
			return true
		}
	}
	return false
}

// ConsumeCardFromHand remueve una carta específica de la mano del jugador.
func (g *GameState) ConsumeCardFromHand(playerID int, unitType string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	p, ok := g.Players[playerID]
	if !ok {
		return false
	}
	for i, c := range p.Hand {
		if c == unitType {
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
			g.HandUpdatedPlayers = append(g.HandUpdatedPlayers, playerID)
			return true
		}
	}
	return false
}

// drawForAllPlayersLocked roba cartas para cada jugador según CardsPerTurn (requiere lock tomado).
// Retorna lista de playerIDs que robaron cartas.
func (g *GameState) drawForAllPlayersLocked() []int {
	updated := []int{}
	for _, p := range g.Players {
		for i := 0; i < g.Config.CardsPerTurn; i++ {
			if _, ok := g.drawCardLocked(p); ok {
				// Solo agregar el playerID una vez, aunque haya robado múltiples cartas
				if i == 0 {
					updated = append(updated, p.ID)
				}
			}
		}
	}
	return updated
}

// AdvancePhase avanza a la siguiente fase del juego
func (g *GameState) AdvancePhase() {
	g.mu.Lock()
	defer g.mu.Unlock()

	switch g.CurrentPhase {
	case PhaseBaseSelection:
		g.CurrentPhase = PhaseTurnStart
		// No resetear ready flags aquí, se hace en TurnStart

	case PhaseTurnStart:
		// Ya en turno, solo pasar a preparación después de la animación breve
		g.CurrentPhase = PhasePreparation
		g.HumanPlayerReady = false
		g.AIPlayerReady = false

	case PhasePreparation:
		g.CurrentPhase = PhaseBattle

	case PhaseBattle:
		g.CurrentPhase = PhaseTurnEnd

	case PhaseTurnEnd:
		// Nuevo turno: robar carta al entrar en turn_start
		updated := g.drawForAllPlayersLocked()
		g.HandUpdatedPlayers = updated
		g.CurrentPhase = PhaseTurnStart
		g.TurnNumber++
	}

	g.PhaseStartTick = g.Tick
	g.PhaseChangedThisTick = true
}

// SetPlayerReady marca al jugador como listo en la fase de preparación
func (g *GameState) SetPlayerReady(playerID int, ready bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if playerID == g.HumanPlayerID {
		g.HumanPlayerReady = ready
	} else if playerID == g.AIPlayerID {
		g.AIPlayerReady = ready
	}
}

// AreBothPlayersReady verifica si ambos jugadores están listos
func (g *GameState) AreBothPlayersReady() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.HumanPlayerReady && g.AIPlayerReady
}

// CanPlayerAct verifica si un jugador puede realizar acciones en la fase actual
func (g *GameState) CanPlayerAct(playerID int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Solo se pueden realizar acciones en la fase de preparación
	return g.CurrentPhase == PhasePreparation
}

// HasPlayerPlacedBase verifica si un jugador ya colocó su base
func (g *GameState) HasPlayerPlacedBase(playerID int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if playerID == g.HumanPlayerID {
		return g.HumanBaseID > 0 && g.Units[g.HumanBaseID] != nil
	} else if playerID == g.AIPlayerID {
		return g.AIBaseID > 0 && g.Units[g.AIBaseID] != nil
	}
	return false
}

// BothBasesPlaced verifica si ambos jugadores colocaron sus bases
func (g *GameState) BothBasesPlaced() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Verificar que ambas bases existan como unidades en el mapa
	humanBaseExists := g.HumanBaseID > 0 && g.Units[g.HumanBaseID] != nil
	aiBaseExists := g.AIBaseID > 0 && g.Units[g.AIBaseID] != nil

	return humanBaseExists && aiBaseExists
}

// MarkBasePlaced marca que un jugador colocó su base
func (g *GameState) MarkBasePlaced(playerID int, baseID int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if playerID == g.HumanPlayerID {
		g.HumanBaseID = baseID
	} else if playerID == g.AIPlayerID {
		g.AIBaseID = baseID
	}
}

// StartFirstTurn inicia el turno 1 después de colocar las bases
func (g *GameState) StartFirstTurn() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.TurnNumber = 1
}

// GetCurrentPhase retorna la fase actual del juego
func (g *GameState) GetCurrentPhase() GamePhase {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.CurrentPhase
}

// DidPhaseChange verifica si la fase cambió este tick y resetea el flag
func (g *GameState) DidPhaseChange() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	changed := g.PhaseChangedThisTick
	g.PhaseChangedThisTick = false // Reset para el próximo tick
	return changed
}

// DrainHandUpdates retorna los playerIDs con manos modificadas y resetea la lista.
func (g *GameState) DrainHandUpdates() []int {
	g.mu.Lock()
	defer g.mu.Unlock()

	updated := g.HandUpdatedPlayers
	g.HandUpdatedPlayers = nil
	return updated
}

// SpawnUnit crea una nueva unidad en el juego
// Requiere que la posición esté dentro del área controlada del jugador
// (excepto cuando aún no colocó base, donde se permite para la base inicial).
func (g *GameState) SpawnUnit(playerID int, unitType string, x, y int) *UnitState {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Validar posición (terreno + ocupación)
	if !g.canUnitTypeEnter(unitType, -1, x, y) {
		return nil
	}

	// Validar que esté dentro del área controlada por el jugador (todas las unidades)
	if !g.isWithinControlledArea(playerID, x, y) {
		return nil
	}

	unit := &UnitState{
		ID:       g.nextUnitID,
		PlayerID: playerID,
		UnitType: unitType,
		X:        x,
		Y:        y,
		HP:       100,
	}

	// Initialize stats based on unit type
	g.applyUnitStats(unit)

	// Set default target: enemy base position (if available), otherwise current position
	if unit.CanMove {
		// Determine enemy base ID depending on the spawning player's side
		enemyBaseID := 0
		if playerID == g.HumanPlayerID {
			enemyBaseID = g.AIBaseID
		} else if playerID == g.AIPlayerID {
			enemyBaseID = g.HumanBaseID
		}

		if enemyBaseID > 0 {
			if enemyBase, ok := g.Units[enemyBaseID]; ok {
				unit.TargetX = enemyBase.X
				unit.TargetY = enemyBase.Y
			} else {
				// Fallback to current position if enemy base not found yet
				unit.TargetX = unit.X
				unit.TargetY = unit.Y
			}
		} else {
			// Fallback to current position if enemy base not assigned yet
			unit.TargetX = unit.X
			unit.TargetY = unit.Y
		}
	} else {
		// Non-movable units keep target at their own tile
		unit.TargetX = unit.X
		unit.TargetY = unit.Y
	}
	g.Units[unit.ID] = unit
	g.nextUnitID++

	return unit
}

// MoveUnit attempts to move a unit to a target position if it's walkable
// and belongs to the requesting player. Returns true on success.
func (g *GameState) MoveUnit(playerID, unitID, x, y int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	unit, ok := g.Units[unitID]
	if !ok {
		return false
	}
	if unit.PlayerID != playerID {
		return false
	}
	// Immediate move (legacy). Prefer SetUnitDestination.
	if !g.isTileAllowedForUnit(unit, x, y) {
		return false
	}

	unit.X = x
	unit.Y = y
	return true
}

// SetUnitDestination sets a target position; unit will step over time.
func (g *GameState) SetUnitDestination(playerID, unitID, x, y int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	unit, ok := g.Units[unitID]
	if !ok {
		return false
	}
	if unit.PlayerID != playerID {
		return false
	}
	if !unit.CanMove {
		return false
	}
	// Destination can be any tile; step validation happens each move tick
	unit.TargetX = x
	unit.TargetY = y
	return true
}

// applyUnitStats assigns movement properties based on UnitType.
func (g *GameState) applyUnitStats(unit *UnitState) {
	stats := GetUnitStats(unit.UnitType)

	// Aplicar stats básicas
	unit.HP = stats.HP
	unit.MaxHP = stats.HP
	unit.Category = stats.Category

	// Aplicar propiedades de movimiento
	unit.CanMove = stats.CanMove
	unit.MoveIntervalTicks = stats.MoveIntervalTicks
	if stats.CanMove {
		unit.NextMoveTick = g.Tick + stats.MoveIntervalTicks
	} else {
		unit.NextMoveTick = 0
	}

	// Aplicar propiedades de detección
	unit.DetectionRange = stats.DetectionRange
	if unit.DetectionRange <= 0 {
		unit.DetectionRange = 5
	}

	// Aplicar propiedades de combate
	unit.AttackDamage = stats.AttackDamage
	unit.AttackRange = stats.AttackRange
	unit.AttackDPS = stats.AttackDPS
	// Si hay DPS configurado, calcular intervalo por ticks en base a AttackDamage
	if unit.AttackDPS > 0 && unit.AttackDamage > 0 {
		tps := g.TicksPerSecond
		if tps <= 0 {
			tps = 5
		}
		// segundos entre ataques = daño por ataque / DPS
		secondsBetween := float64(unit.AttackDamage) / unit.AttackDPS
		intervalTicks := int(secondsBetween*float64(tps) + 0.5) // redondeo
		if intervalTicks < 1 {
			intervalTicks = 1
		}
		unit.AttackIntervalTicks = intervalTicks
	} else {
		unit.AttackIntervalTicks = stats.AttackIntervalTicks
	}
	if unit.AttackDamage > 0 {
		unit.NextAttackTick = g.Tick + unit.AttackIntervalTicks
	} else {
		unit.NextAttackTick = 0
	}

	// Aplicar propiedades de generador
	unit.IsGenerator = stats.IsGenerator
	unit.GeneratedUnitType = stats.GeneratedUnitType
	unit.GenerationInterval = stats.GenerationInterval
	unit.MaxUnitsGenerated = stats.MaxUnitsGenerated
	if stats.IsGenerator {
		unit.NextGenerationTick = g.Tick + stats.GenerationInterval
	} else {
		unit.NextGenerationTick = 0
	}

	// Aplicar propiedades de bloqueo
	unit.IsBlocker = stats.IsBlocker

	// Aplicar rango de construcción
	unit.BuildRange = stats.BuildRange

	// Estado inicial
	unit.Status = "idle"
}

// canUnitTypeEnter checks if a unit of unitType can enter tile (x,y).
// skipUnitID allows ignoring a specific unit occupying that tile (useful for movement of that unit).
func (g *GameState) canUnitTypeEnter(unitType string, skipUnitID int, x, y int) bool {
	// Bounds & terrain
	tile, ok := g.Map.GetTile(x, y)
	if !ok {
		return false
	}

	stats := GetUnitStats(unitType)

	switch stats.Category {
	case CategoryNavalUnit:
		// Navales solo en agua
		if tile.TerrainID != TerrainWater {
			return false
		}
	default:
		// Estructuras y terrestres solo en tiles walkable (no agua)
		if !tile.Walkable {
			return false
		}
	}

	// Ocupación: no permitir dos unidades en el mismo tile y respetar bloqueadores
	for _, other := range g.Units {
		if other.ID == skipUnitID {
			continue
		}
		if other.X == x && other.Y == y {
			return false
		}
	}

	return true
}

// isWithinControlledArea verifica si una posición está dentro del área controlada por un jugador.
// El área controlada está determinada por la base principal y las estructuras con BuildRange > 0.
func (g *GameState) isWithinControlledArea(playerID int, x, y int) bool {
	// Si el jugador no tiene base aún, permitir spawneo libre (para colocar la base inicial)
	baseID := 0
	if playerID == g.HumanPlayerID {
		baseID = g.HumanBaseID
	} else if playerID == g.AIPlayerID {
		baseID = g.AIBaseID
	}

	if baseID == 0 {
		return true // Permite colocar la base inicial en cualquier lugar
	}

	// Verificar si está dentro del rango de alguna estructura del jugador
	for _, unit := range g.Units {
		if unit.PlayerID != playerID {
			continue
		}
		if unit.HP <= 0 {
			continue
		}
		if unit.BuildRange <= 0 {
			continue
		}

		// Calcular distancia Manhattan
		dx := x - unit.X
		if dx < 0 {
			dx = -dx
		}
		dy := y - unit.Y
		if dy < 0 {
			dy = -dy
		}
		dist := dx + dy

		if dist <= unit.BuildRange {
			return true
		}
	}

	return false
}

// isTileAllowedForUnit checks terrain and blocking constraints for the given unit.
func (g *GameState) isTileAllowedForUnit(unit *UnitState, x, y int) bool {
	return g.canUnitTypeEnter(unit.UnitType, unit.ID, x, y)
}

// findSpawnPosition intenta encontrar una posición válida para un unitType.
// Valida terreno, ocupación y área controlada (todas las unidades).
func (g *GameState) findSpawnPosition(unitType string, playerID int, attempts int) (int, int, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for i := 0; i < attempts; i++ {
		x := rand.Intn(g.Map.Width)
		y := rand.Intn(g.Map.Height)

		// Validar terreno y ocupación
		if !g.canUnitTypeEnter(unitType, -1, x, y) {
			continue
		}

		// Validar área controlada para todos los spawns (excepto base inicial, manejado en isWithinControlledArea)
		if !g.isWithinControlledArea(playerID, x, y) {
			continue
		}

		return x, y, true
	}
	return 0, 0, false
}
