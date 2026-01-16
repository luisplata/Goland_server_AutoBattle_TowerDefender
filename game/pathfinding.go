package game

import (
	"container/heap"
)

const maxPathSearchSteps = MapWidth * MapHeight // allow full-map searches to avoid shoreline stalls

// PathNode representa un nodo en la búsqueda A*
type PathNode struct {
	X      int
	Y      int
	GCost  float64 // Costo desde el inicio
	HCost  float64 // Costo heurístico al objetivo
	FCost  float64 // GCost + HCost
	Parent *PathNode
	Index  int // Para mantener la posición en el heap
}

// NodeHeap es una priority queue de PathNodes
type NodeHeap []*PathNode

func (h NodeHeap) Len() int {
	return len(h)
}

func (h NodeHeap) Less(i, j int) bool {
	// Comparar por FCost, y si son iguales, por HCost (tie-breaker)
	if h[i].FCost != h[j].FCost {
		return h[i].FCost < h[j].FCost
	}
	return h[i].HCost < h[j].HCost
}

func (h NodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

func (h *NodeHeap) Push(x interface{}) {
	node := x.(*PathNode)
	node.Index = len(*h)
	*h = append(*h, node)
}

func (h *NodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	node := old[n-1]
	*h = old[0 : n-1]
	return node
}

// PathCache almacena paths calculados para evitar recalcular
type PathCache struct {
	paths map[string][]Point
}

func NewPathCache() *PathCache {
	return &PathCache{
		paths: make(map[string][]Point),
	}
}

func (pc *PathCache) GetKey(startX, startY, endX, endY int) string {
	return string([]byte{
		byte((startX >> 8) & 0xFF), byte(startX & 0xFF),
		byte((startY >> 8) & 0xFF), byte(startY & 0xFF),
		byte((endX >> 8) & 0xFF), byte(endX & 0xFF),
		byte((endY >> 8) & 0xFF), byte(endY & 0xFF),
	})
}

func (pc *PathCache) Get(startX, startY, endX, endY int) ([]Point, bool) {
	key := pc.GetKey(startX, startY, endX, endY)
	path, ok := pc.paths[key]
	return path, ok
}

func (pc *PathCache) Set(startX, startY, endX, endY int, path []Point) {
	key := pc.GetKey(startX, startY, endX, endY)
	pc.paths[key] = path
}

func (pc *PathCache) Clear() {
	pc.paths = make(map[string][]Point)
}

// Point representa una coordenada simple
type Point struct {
	X int
	Y int
}

// PathFinder es el motor de búsqueda A*
type PathFinder struct {
	cache *PathCache
}

func NewPathFinder() *PathFinder {
	return &PathFinder{
		cache: NewPathCache(),
	}
}

// FindPath encuentra un camino desde (startX, startY) a (endX, endY) usando A*
// Returns lista de pasos o nil si no hay camino
func (pf *PathFinder) FindPath(state *GameState, unit *UnitState, startX, startY, endX, endY int, maxSteps int) []Point {
	// Si ya estamos en el destino
	if startX == endX && startY == endY {
		return []Point{{X: startX, Y: startY}}
	}

	// Verificar cache
	if cached, ok := pf.cache.Get(startX, startY, endX, endY); ok {
		return cached
	}

	// Inicializar búsqueda
	openSet := &NodeHeap{}
	heap.Init(openSet)

	closedSet := make(map[string]bool)
	openMap := make(map[string]*PathNode)

	startNode := &PathNode{
		X:     startX,
		Y:     startY,
		GCost: 0,
		HCost: pf.heuristic(startX, startY, endX, endY),
	}
	startNode.FCost = startNode.GCost + startNode.HCost

	heap.Push(openSet, startNode)
	openMap[pf.nodeKey(startX, startY)] = startNode

	// Direcciones: arriba, abajo, izquierda, derecha (evita diagonales que complican más)
	directions := []Point{
		{X: 0, Y: 1},  // arriba
		{X: 0, Y: -1}, // abajo
		{X: 1, Y: 0},  // derecha
		{X: -1, Y: 0}, // izquierda
	}

	steps := 0

	for openSet.Len() > 0 && steps < maxSteps {
		steps++

		current := heap.Pop(openSet).(*PathNode)
		delete(openMap, pf.nodeKey(current.X, current.Y))

		// Alcanzamos el objetivo
		if current.X == endX && current.Y == endY {
			return pf.reconstructPath(current, startX, startY, endX, endY)
		}

		closedSet[pf.nodeKey(current.X, current.Y)] = true

		// Explorar vecinos
		for _, dir := range directions {
			neighborX := current.X + dir.X
			neighborY := current.Y + dir.Y

			neighborKey := pf.nodeKey(neighborX, neighborY)

			// Si está en closed set, ignorar
			if closedSet[neighborKey] {
				continue
			}

			// Verificar si el tile es válido y accesible
			if !state.isTileAllowedForUnit(unit, neighborX, neighborY) {
				continue
			}

			// Calcular costos
			tentativeG := current.GCost + 1.0

			neighbor, exists := openMap[neighborKey]
			if !exists {
				// Nuevo nodo
				h := pf.heuristic(neighborX, neighborY, endX, endY)
				neighbor = &PathNode{
					X:      neighborX,
					Y:      neighborY,
					GCost:  tentativeG,
					HCost:  h,
					FCost:  tentativeG + h,
					Parent: current,
				}
				openMap[neighborKey] = neighbor
				heap.Push(openSet, neighbor)
			} else if tentativeG < neighbor.GCost {
				// Encontramos un camino mejor
				neighbor.GCost = tentativeG
				neighbor.FCost = neighbor.GCost + neighbor.HCost
				neighbor.Parent = current
				// Actualizar posición en el heap
				heap.Fix(openSet, neighbor.Index)
			}
		}
	}

	// No hay camino
	return nil
}

// heuristic calcula la distancia Manhattan
func (pf *PathFinder) heuristic(x, y, goalX, goalY int) float64 {
	dx := abs(goalX - x)
	dy := abs(goalY - y)
	return float64(dx + dy)
}

// sign retorna -1, 0, 1 según el valor
func sign(v int) int {
	if v > 0 {
		return 1
	}
	if v < 0 {
		return -1
	}
	return 0
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// nodeKey genera una clave única para un nodo
func (pf *PathFinder) nodeKey(x, y int) string {
	return string([]byte{
		byte((x >> 8) & 0xFF), byte(x & 0xFF),
		byte((y >> 8) & 0xFF), byte(y & 0xFF),
	})
}

// reconstructPath reconstruye el camino desde el nodo final
func (pf *PathFinder) reconstructPath(node *PathNode, startX, startY, endX, endY int) []Point {
	path := []Point{}

	for node != nil {
		path = append([]Point{{X: node.X, Y: node.Y}}, path...)
		node = node.Parent
	}

	// Cachear el resultado
	if len(path) > 0 {
		pf.cache.Set(startX, startY, endX, endY, path)
	}

	return path
}

// GetNextStep retorna el siguiente paso en el camino
// Si no hay camino, intenta movimiento directo (fallback)
func (pf *PathFinder) GetNextStep(state *GameState, unit *UnitState, targetX, targetY int) (int, int, bool) {
	goalX, goalY, ok := pf.selectGoalTile(state, unit, targetX, targetY)
	if !ok {
		return unit.X, unit.Y, false
	}

	path := pf.FindPath(state, unit, unit.X, unit.Y, goalX, goalY, maxPathSearchSteps)

	// Si se encontró camino
	if len(path) > 1 {
		nextStep := path[1]
		// Si el siguiente paso está bloqueado ahora (otro unit), invalidar cache y recomputar una vez
		if !state.isTileAllowedForUnit(unit, nextStep.X, nextStep.Y) {
			pf.InvalidatePath(unit.X, unit.Y, goalX, goalY)
			path = pf.FindPath(state, unit, unit.X, unit.Y, goalX, goalY, maxPathSearchSteps)
			if len(path) > 1 {
				nextStep = path[1]
				if state.isTileAllowedForUnit(unit, nextStep.X, nextStep.Y) {
					return nextStep.X, nextStep.Y, true
				}
			}
		} else {
			return nextStep.X, nextStep.Y, true
		}
	}

	// Fallback: probar todas las 8 direcciones ordenadas por heurística al objetivo
	type cand struct {
		x, y     int
		distance int // Manhattan distance to target
	}
	dx := targetX - unit.X
	dy := targetY - unit.Y

	candidates := []cand{
		// Preferencia 1: hacia el eje dominante (2 direcciones ortogonales)
		{unit.X + sign(dx), unit.Y, abs(targetX-(unit.X+sign(dx))) + abs(targetY-unit.Y)},
		{unit.X, unit.Y + sign(dy), abs(targetX-unit.X) + abs(targetY-(unit.Y+sign(dy)))},
		// Preferencia 2: diagonales hacia el objetivo
		{unit.X + sign(dx), unit.Y + sign(dy), abs(targetX-(unit.X+sign(dx))) + abs(targetY-(unit.Y+sign(dy)))},
		// Preferencia 3: direcciones ortogonales perpendiculares
		{unit.X - sign(dx), unit.Y, abs(targetX-(unit.X-sign(dx))) + abs(targetY-unit.Y)},
		{unit.X, unit.Y - sign(dy), abs(targetX-unit.X) + abs(targetY-(unit.Y-sign(dy)))},
		// Preferencia 4: diagonal opuesta
		{unit.X - sign(dx), unit.Y - sign(dy), abs(targetX-(unit.X-sign(dx))) + abs(targetY-(unit.Y-sign(dy)))},
		// Preferencia 5: diagonal perpendicular 1
		{unit.X + sign(dx), unit.Y - sign(dy), abs(targetX-(unit.X+sign(dx))) + abs(targetY-(unit.Y-sign(dy)))},
		// Preferencia 6: diagonal perpendicular 2
		{unit.X - sign(dx), unit.Y + sign(dy), abs(targetX-(unit.X-sign(dx))) + abs(targetY-(unit.Y+sign(dy)))},
	}

	// Ordenar por distancia al objetivo (heurística)
	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].distance < candidates[i].distance {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}

	for _, c := range candidates {
		if state.isTileAllowedForUnit(unit, c.x, c.y) {
			return c.x, c.y, true
		}
	}

	// No hay movimiento posible
	return unit.X, unit.Y, false
}

// selectGoalTile ajusta el destino para que sea un tile alcanzable: si el tile objetivo
// no se puede pisar (agua para terrestres, ocupado, etc.), busca la mejor casilla
// alrededor del objetivo dentro del rango de ataque.
func (pf *PathFinder) selectGoalTile(state *GameState, unit *UnitState, targetX, targetY int) (int, int, bool) {
	// Si el tile objetivo es pisable, úsalo directamente
	if state.canUnitTypeEnter(unit.UnitType, unit.ID, targetX, targetY) {
		return targetX, targetY, true
	}

	// Buscar casillas alcanzables dentro del rango de ataque (o 1 si es 0)
	radius := unit.AttackRange
	if radius < 1 {
		radius = 1
	}

	bestFound := false
	bestX, bestY := 0, 0
	bestDistFromUnit := 1_000_000
	bestDistToTarget := 1_000_000

	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			manToTarget := abs(dx) + abs(dy)
			if manToTarget > radius {
				continue
			}

			nx := targetX + dx
			ny := targetY + dy
			if !state.canUnitTypeEnter(unit.UnitType, unit.ID, nx, ny) {
				continue
			}

			distFromUnit := abs(unit.X-nx) + abs(unit.Y-ny)
			if !bestFound || distFromUnit < bestDistFromUnit || (distFromUnit == bestDistFromUnit && manToTarget < bestDistToTarget) {
				bestFound = true
				bestX, bestY = nx, ny
				bestDistFromUnit = distFromUnit
				bestDistToTarget = manToTarget
			}
		}
	}

	if bestFound {
		return bestX, bestY, true
	}

	return targetX, targetY, false
}

// ClearCache limpia el cache de paths
func (pf *PathFinder) ClearCache() {
	pf.cache.Clear()
}

// InvalidatePath invalida un path específico cuando hay cambios en el mapa
func (pf *PathFinder) InvalidatePath(startX, startY, endX, endY int) {
	key := pf.cache.GetKey(startX, startY, endX, endY)
	delete(pf.cache.paths, key)
}
