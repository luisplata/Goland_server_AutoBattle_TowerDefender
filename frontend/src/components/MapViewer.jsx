import { useState, useRef, useEffect, useCallback } from 'react'
import './MapViewer.css'

const TERRAIN_COLORS = {
  0: '#4a7c3c', // Grass
  1: '#6b8e23', // Path
  2: '#1e5aa0', // Water
}

// ========== COLORES DE UNIDADES - EDITA AQUI ==========
const COLOR_CONFIG = {
  PLAYER_1_HUE: 200,        // Azul
  PLAYER_2_HUE: 0,          // Rojo
  SATURATION: 100,          // Saturación (0-100)
  LIGHTNESS_MIN: 15,        // Brillo mínimo (muy oscuro - casi negro)
  LIGHTNESS_MAX: 85,        // Brillo máximo (muy claro - casi blanco)
}

// Intensidades relativas de tipos de unidad (0 = muy oscuro, 1 = muy claro)
const UNIT_TYPE_INTENSITIES = {
  main_base: 0.0,           // Muy oscuro (15%)
  tower: 0.14,              // Oscuro (24%)
  wall: 0.28,               // Medio-oscuro (34%)
  land_generator: 0.42,     // Medio (44%)
  naval_generator: 0.5,     // Medio (50%)
  warrior: 0.64,            // Claro (64%)
  land_soldier: 0.78,       // Más claro (75%)
  naval_ship: 1.0,          // Muy claro (85%)
}
// ====================================================

const getUnitColorIntensity = (unitType) => {
  return UNIT_TYPE_INTENSITIES[unitType] ?? 0.5
}

// Generar color HSL basado en bando y tipo
const getUnitColor = (playerId, unitType) => {
  const hue = playerId === 1 ? COLOR_CONFIG.PLAYER_1_HUE : COLOR_CONFIG.PLAYER_2_HUE
  const saturation = COLOR_CONFIG.SATURATION
  const intensity = getUnitColorIntensity(unitType)
  const lightness = COLOR_CONFIG.LIGHTNESS_MIN + (intensity * (COLOR_CONFIG.LIGHTNESS_MAX - COLOR_CONFIG.LIGHTNESS_MIN))
  return `hsl(${hue}, ${saturation}%, ${lightness}%)`
}

export default function MapViewer({ gameMap, units, selectedTile, onSelectTile, disableZoom = false, playerId, selectedCard }) {
  const [zoom, setZoom] = useState(3)
  const [pan, setPan] = useState({ x: 0, y: 0 })
  const [isPanning, setIsPanning] = useState(false)
  const [startPan, setStartPan] = useState({ x: 0, y: 0 })
  const mapContainerRef = useRef(null)
  
  // Usar useCallback para el handler del wheel
  const handleWheel = useCallback((e) => {
    e.preventDefault()
    
    if (disableZoom) return

    const container = mapContainerRef.current
    if (!container) return

    const rect = container.getBoundingClientRect()
    const delta = e.deltaY > 0 ? -0.2 : 0.2
    
    // Mouse position relative to container
    const mouseX = e.clientX - rect.left
    const mouseY = e.clientY - rect.top
    
    // Handle zoom at mouse position
    setZoom(oldZoom => {
      const newZoom = Math.max(0.5, Math.min(10, oldZoom + delta))
      
      if (oldZoom === newZoom) return oldZoom // No change in zoom
      
      // Calculate the new pan to keep the point under the mouse stationary
      // Formula: newPan = mousePos - (mousePos - oldPan) * (newZoom / oldZoom)
      setPan(prevPan => {
        const zoomRatio = newZoom / oldZoom
        return {
          x: mouseX - (mouseX - prevPan.x) * zoomRatio,
          y: mouseY - (mouseY - prevPan.y) * zoomRatio
        }
      })
      
      return newZoom
    })
  }, [disableZoom])

  // Zoom at center of viewport (for buttons)
  const handleZoomAtCenter = useCallback((delta) => {
    if (disableZoom) return

    const container = mapContainerRef.current
    if (!container) return

    const rect = container.getBoundingClientRect()
    const centerX = rect.width / 2
    const centerY = rect.height / 2
    
    setZoom(oldZoom => {
      const newZoom = Math.max(0.5, Math.min(10, oldZoom + delta))
      
      if (oldZoom === newZoom) return oldZoom
      
      setPan(prevPan => {
        const zoomRatio = newZoom / oldZoom
        return {
          x: centerX - (centerX - prevPan.x) * zoomRatio,
          y: centerY - (centerY - prevPan.y) * zoomRatio
        }
      })
      
      return newZoom
    })
  }, [disableZoom])
  
  // Agregar wheel listener con { passive: false } para poder usar preventDefault
  useEffect(() => {
    const container = mapContainerRef.current
    if (!container) return

    container.addEventListener('wheel', handleWheel, { passive: false })
    
    return () => {
      container.removeEventListener('wheel', handleWheel)
    }
  }, [handleWheel])
  
  if (!gameMap || !gameMap.tiles) {
    return <div className="map-viewer">Loading map...</div>
  }

  const tileSize = 3 // pixels per tile
  const mapWidth = gameMap.width * tileSize
  const mapHeight = gameMap.height * tileSize

  const handleMouseDown = (e) => {
    if (e.button === 1 || e.shiftKey) { // Middle mouse or Shift+Click for pan
      setIsPanning(true)
      setStartPan({ x: e.clientX - pan.x, y: e.clientY - pan.y })
      e.preventDefault()
    }
  }

  const handleMouseMove = (e) => {
    if (isPanning) {
      setPan({
        x: e.clientX - startPan.x,
        y: e.clientY - startPan.y
      })
    }
  }

  const handleMouseUp = () => {
    setIsPanning(false)
  }

  const resetView = () => {
    setZoom(3)
    setPan({ x: 0, y: 0 })
  }

  // Calculate controlled area when player is selecting a card
  const calculateControlledArea = () => {
    if (!playerId || !units || !selectedCard) return new Set()

    const controlled = new Set()
    const myUnits = Object.values(units).filter(u => u.playerId === playerId && u.hp > 0)
    
    for (const unit of myUnits) {
      if (!unit.buildRange || unit.buildRange <= 0) continue
      
      // Add all tiles within Manhattan distance of buildRange
      for (let dy = -unit.buildRange; dy <= unit.buildRange; dy++) {
        for (let dx = -unit.buildRange; dx <= unit.buildRange; dx++) {
          const manhattanDist = Math.abs(dx) + Math.abs(dy)
          if (manhattanDist <= unit.buildRange) {
            const tileX = unit.x + dx
            const tileY = unit.y + dy
            if (tileX >= 0 && tileX < gameMap.width && tileY >= 0 && tileY < gameMap.height) {
              controlled.add(`${tileX},${tileY}`)
            }
          }
        }
      }
    }
    
    return controlled
  }

  const controlledArea = calculateControlledArea()
  
  // Check if selected card is a structure (requires controlled area)
  const isStructureCard = selectedCard && ['tower', 'wall', 'land_generator', 'naval_generator'].includes(selectedCard)

  const handleTileClick = (e) => {
    if (isPanning || e.button !== 0) return // Only handle left click for tile selection
    const svg = e.currentTarget
    const pt = svg.createSVGPoint()
    pt.x = e.clientX
    pt.y = e.clientY
    const svgP = pt.matrixTransform(svg.getScreenCTM().inverse())
    
    const x = Math.floor(svgP.x)
    const y = Math.floor(svgP.y)
    
    if (x < 0 || x >= gameMap.width || y < 0 || y >= gameMap.height) return
    const tile = gameMap.tiles[y][x]
    onSelectTile && onSelectTile({ x, y, walkable: tile.walkable })
  }

  return (
    <div className="map-viewer">
      <h3>🗺️ Game Map</h3>
      
      <div className="map-controls">
        <button onClick={() => handleZoomAtCenter(0.2)} className="zoom-btn" title="Zoom In" disabled={disableZoom}>🔍 +</button>
        <button onClick={() => handleZoomAtCenter(-0.2)} className="zoom-btn" title="Zoom Out" disabled={disableZoom}>🔍 -</button>
        <button onClick={resetView} className="zoom-btn" title="Reset View" disabled={disableZoom}>🎯</button>
        <span className="zoom-level">{((zoom / 3) * 100).toFixed(0)}%</span>
      </div>

      <div 
        className="map-container"
        ref={mapContainerRef}
        onMouseDown={(e) => {
          handleMouseDown(e)
        }}
        onMouseMove={handleMouseMove}
        onMouseUp={handleMouseUp}
        onMouseLeave={handleMouseUp}
        onContextMenu={(e) => e.preventDefault()}
        style={{ cursor: isPanning ? 'grabbing' : 'crosshair' }}
      >
        <svg 
          width={mapWidth} 
          height={mapHeight} 
          className="map-svg"
          onMouseDown={handleTileClick}
          style={{ 
            transform: `translate(${pan.x}px, ${pan.y}px) scale(${zoom})`, 
            transformOrigin: '0 0' 
          }}
        >
        {/* Tiles */}
        {gameMap.tiles.map((row, y) =>
          row.map((tile, x) => {
            const isControlled = controlledArea.has(`${x},${y}`)
            const showControlled = selectedCard && controlledArea.size > 0
            // For structures, show strict control (red outside). For units, show info only (blue tint)
            const isValidForStructure = !isStructureCard || isControlled
            
            return (
              <rect
                key={`${x}-${y}`}
                x={x}
                y={y}
                width={1}
                height={1}
                fill={TERRAIN_COLORS[tile.terrainId] || '#666'}
                opacity={tile.walkable ? 1 : 0.6}
                stroke={showControlled ? (isValidForStructure ? (isControlled ? '#00ff88' : 'none') : '#ff4444') : 'none'}
                strokeWidth={showControlled ? 0.02 : 0}
                fillOpacity={showControlled && isStructureCard && !isControlled ? 0.3 : 1}
              />
            )
          })
        )}

        {/* Selection highlight */}
        {selectedTile && (
          <rect
            x={selectedTile.x}
            y={selectedTile.y}
            width={1}
            height={1}
            fill="none"
            stroke={selectedTile.walkable ? '#00ff88' : '#ff4444'}
            strokeWidth={0.1}
          />
        )}

        {/* Units */}
        {units && Object.values(units).map(unit => {
          const hpPercent = unit.maxHp ? (unit.hp / unit.maxHp) : 1
          const color = getUnitColor(unit.playerId, unit.unitType)
          const cx = unit.x + 0.5 // center of tile instead of vertex
          const cy = unit.y + 0.5
          
          return (
            <g key={unit.id}>
              <circle
                cx={cx}
                cy={cy}
                r={0.4}
                fill={color}
                stroke="white"
                strokeWidth="0.05"
              />
              {/* HP bar over unit */}
              <rect
                x={cx - 0.3}
                y={cy - 0.7}
                width={0.6}
                height={0.08}
                fill="#1a1a1a"
                stroke="none"
              />
              <rect
                x={cx - 0.3}
                y={cy - 0.7}
                width={0.6 * hpPercent}
                height={0.08}
                fill={hpPercent > 0.5 ? '#4CAF50' : hpPercent > 0.25 ? '#FF9800' : '#F44336'}
                stroke="none"
              />
              {/* Attack indicator for units with damage */}
              {unit.attackDamage > 0 && (
                <circle
                  cx={cx}
                  cy={cy}
                  r={unit.attackRange || 1}
                  fill="none"
                  stroke={color}
                  strokeWidth="0.02"
                  opacity="0.2"
                />
              )}
              {/* Build range indicator for structures */}
              {unit.buildRange > 0 && (
                <circle
                  cx={cx}
                  cy={cy}
                  r={unit.buildRange}
                  fill="none"
                  stroke="#00ff88"
                  strokeWidth="0.03"
                  opacity="0.15"
                  strokeDasharray="0.2,0.2"
                />
              )}
            </g>
          )
        })}
      </svg>      </div>      
      <div className="map-legend">
        <div className="legend-item">
          <div className="legend-color grass"></div>
          <span>Walkable</span>
        </div>
        <div className="legend-item">
          <div className="legend-color water"></div>
          <span>Water (Invalid)</span>
        </div>
        <div style={{ marginTop: '0.5rem', borderTop: '1px solid #666', paddingTop: '0.5rem' }}>
          <div style={{ fontSize: '0.85rem', marginBottom: '0.3rem', fontWeight: 'bold' }}>Unit Colors:</div>
          <div className="legend-item" style={{ fontSize: '0.8rem' }}>
            <div style={{ width: '12px', height: '12px', borderRadius: '50%', backgroundColor: getUnitColor(1, 'main_base'), border: '1px solid white' }}></div>
            <span>Your Base</span>
          </div>
          <div className="legend-item" style={{ fontSize: '0.8rem' }}>
            <div style={{ width: '12px', height: '12px', borderRadius: '50%', backgroundColor: getUnitColor(1, 'warrior'), border: '1px solid white' }}></div>
            <span>Your Warrior</span>
          </div>
          <div className="legend-item" style={{ fontSize: '0.8rem' }}>
            <div style={{ width: '12px', height: '12px', borderRadius: '50%', backgroundColor: getUnitColor(2, 'main_base'), border: '1px solid white' }}></div>
            <span>Enemy Base</span>
          </div>
          <div className="legend-item" style={{ fontSize: '0.8rem' }}>
            <div style={{ width: '12px', height: '12px', borderRadius: '50%', backgroundColor: getUnitColor(2, 'warrior'), border: '1px solid white' }}></div>
            <span>Enemy Warrior</span>
          </div>
        </div>
      </div>

      <div className="selection-info">
        {selectedTile ? (
          <span>
            Selected: ({selectedTile.x}, {selectedTile.y}) — {selectedTile.walkable ? 'Walkable ✅' : 'Water ❌'}
          </span>
        ) : selectedCard && controlledArea.size > 0 ? (
          isStructureCard ? (
            <span>🟢 Green border: Can build here | 🔴 Red border: Outside controlled area</span>
          ) : (
            <span>🟢 Green border: Your controlled area (units can spawn anywhere walkable)</span>
          )
        ) : (
          <span>Scroll: Zoom | Shift+Drag or Middle Click: Pan</span>
        )}
      </div>
    </div>
  )
}
