import './MapViewer.css'

const TERRAIN_COLORS = {
  0: '#4a7c3c', // Grass
  1: '#6b8e23', // Path
  2: '#1e5aa0', // Water
}

export default function MapViewer({ gameMap, units }) {
  if (!gameMap || !gameMap.tiles) {
    return <div className="map-viewer">Loading map...</div>
  }

  const tileSize = 3 // pixels per tile
  const mapWidth = gameMap.width * tileSize
  const mapHeight = gameMap.height * tileSize

  return (
    <div className="map-viewer">
      <h3>🗺️ Game Map</h3>
      <svg 
        width={mapWidth} 
        height={mapHeight} 
        className="map-svg"
        viewBox={`0 0 ${gameMap.width} ${gameMap.height}`}
      >
        {/* Tiles */}
        {gameMap.tiles.map((row, y) =>
          row.map((tile, x) => (
            <rect
              key={`${x}-${y}`}
              x={x}
              y={y}
              width={1}
              height={1}
              fill={TERRAIN_COLORS[tile.terrainId] || '#666'}
              opacity={tile.walkable ? 1 : 0.6}
              stroke="none"
            />
          ))
        )}

        {/* Units */}
        {units && Object.values(units).map(unit => (
          <circle
            key={unit.id}
            cx={unit.x}
            cy={unit.y}
            r={0.4}
            fill={unit.playerId === 1 ? '#00ff88' : '#ff00ff'}
            stroke="white"
            strokeWidth="0.05"
          />
        ))}
      </svg>
      
      <div className="map-legend">
        <div className="legend-item">
          <div className="legend-color grass"></div>
          <span>Walkable</span>
        </div>
        <div className="legend-item">
          <div className="legend-color water"></div>
          <span>Water (Invalid)</span>
        </div>
        <div className="legend-item">
          <div className="legend-color unit1"></div>
          <span>Your Units</span>
        </div>
        <div className="legend-item">
          <div className="legend-color unit2"></div>
          <span>Enemy Units</span>
        </div>
      </div>
    </div>
  )
}
