package game

import (
	"math"
)

const (
	MapWidth  = 100
	MapHeight = 100
)

// Tipos de terreno
const (
	TerrainGrass = 0
	TerrainPath  = 1
	TerrainWater = 2
)

type Tile struct {
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Walkable  bool `json:"walkable"`
	TerrainID int  `json:"terrainId"`
}

type GameMap struct {
	Width  int      `json:"width"`
	Height int      `json:"height"`
	Tiles  [][]Tile `json:"tiles"`
}

func NewGameMap() *GameMap {
	gameMap := &GameMap{
		Width:  MapWidth,
		Height: MapHeight,
		Tiles:  make([][]Tile, MapHeight),
	}

	// Generar mapa con Perlin noise
	for y := 0; y < MapHeight; y++ {
		gameMap.Tiles[y] = make([]Tile, MapWidth)
		for x := 0; x < MapWidth; x++ {
			// Generar valor de noise
			noiseValue := perlinNoise(float64(x), float64(y))

			// Mapear noise a tipo de terreno
			terrainID := TerrainGrass
			walkable := true

			if noiseValue < 0.3 {
				// Agua
				terrainID = TerrainWater
				walkable = false
			} else if noiseValue < 0.4 {
				// Camino
				terrainID = TerrainPath
				walkable = true
			} else {
				// Pasto
				terrainID = TerrainGrass
				walkable = true
			}

			gameMap.Tiles[y][x] = Tile{
				X:         x,
				Y:         y,
				Walkable:  walkable,
				TerrainID: terrainID,
			}
		}
	}

	return gameMap
}

// Implementación simple de Perlin noise 2D
func perlinNoise(x, y float64) float64 {
	// Escalar coordenadas
	scale := 0.05
	x *= scale
	y *= scale

	// Obtener coordenadas de la celda
	x0 := math.Floor(x)
	x1 := x0 + 1
	y0 := math.Floor(y)
	y1 := y0 + 1

	// Posición relativa dentro de la celda
	sx := x - x0
	sy := y - y0

	// Gradientes en las esquinas
	n0 := dotGridGradient(int(x0), int(y0), x, y)
	n1 := dotGridGradient(int(x1), int(y0), x, y)
	ix0 := interpolate(n0, n1, sx)

	n0 = dotGridGradient(int(x0), int(y1), x, y)
	n1 = dotGridGradient(int(x1), int(y1), x, y)
	ix1 := interpolate(n0, n1, sx)

	value := interpolate(ix0, ix1, sy)

	// Normalizar a [0, 1]
	return (value + 1) / 2
}

func dotGridGradient(ix, iy int, x, y float64) float64 {
	// Vector gradiente pseudoaleatorio
	angle := pseudoRandom(ix, iy) * 2 * math.Pi
	gx := math.Cos(angle)
	gy := math.Sin(angle)

	// Vector distancia
	dx := x - float64(ix)
	dy := y - float64(iy)

	return dx*gx + dy*gy
}

func pseudoRandom(x, y int) float64 {
	// Hash simple para generar valor pseudoaleatorio
	n := x*374761393 + y*668265263
	n = (n ^ (n >> 13)) * 1274126177
	return float64(n&0x7fffffff) / 2147483648.0
}

func interpolate(a, b, t float64) float64 {
	// Interpolación suave (smoothstep)
	t = t * t * (3 - 2*t)
	return a + t*(b-a)
}

func (m *GameMap) IsWalkable(x, y int) bool {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return false
	}
	return m.Tiles[y][x].Walkable
}
