import { useState, useRef, useEffect } from 'react'
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

export default function MapViewer({ gameMap, units, selectedTile, onSelectTile, disableZoom = false }) {
  const [zoom, setZoom] = useState(1)
  const [pan, setPan] = useState({ x: 0, y: 0 })
  const [isPanning, setIsPanning] = useState(false)
  const [startPan, setStartPan] = useState({ x: 0, y: 0 })
  const mapContainerRef = useRef(null)
  
  // Agregar wheel listener con { passive: false } para poder usar preventDefault
  useEffect(() => {
    const container = mapContainerRef.current
    if (!container) return

    const wheelHandler = (e) => {
      e.preventDefault()
      
      if (disableZoom) return

      const rect = container.getBoundingClientRect()
      const delta = e.deltaY > 0 ? -0.1 : 0.1
      handleZoomAtPoint(delta, e.clientX, e.clientY, rect)
    }

    container.addEventListener('wheel', wheelHandler, { passive: false })
    
    return () => {
      container.removeEventListener('wheel', wheelHandler)
    }
  }, [disableZoom, zoom])
  
  if (!gameMap || !gameMap.tiles) {
    return <div className="map-viewer">Loading map...</div>
  }

  const tileSize = 3 // pixels per tile
  const mapWidth = gameMap.width * tileSize
  const mapHeight = gameMap.height * tileSize

  const handleZoomAtPoint = (delta, mouseX, mouseY, containerRect) => {
    const oldZoom = zoom
    const newZoom = Math.max(0.5, Math.min(5, oldZoom + delta))
    
    if (oldZoom === newZoom) return // No change in zoom
    
    // Mouse position relative to container
    const mouseXInContainer = mouseX - containerRect.left
    const mouseYInContainer = mouseY - containerRect.top
    
    // Calculate the new pan to keep the point under the mouse stationary
    const scale = newZoom / oldZoom
    
    setZoom(newZoom)
    setPan(prevPan => ({
      x: mouseXInContainer - (mouseXInContainer - prevPan.x) * scale,
      y: mouseYInContainer - (mouseYInContainer - prevPan.y) * scale
    }))
  }

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
    setZoom(1)
    setPan({ x: 0, y: 0 })
  }

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
        <button onClick={(e) => {
          if (disableZoom) return
          const container = e.target.closest('.map-viewer').querySelector('.map-container')
          const rect = container.getBoundingClientRect()
          handleZoomAtPoint(0.2, rect.left + rect.width / 2, rect.top + rect.height / 2, rect)
        }} className="zoom-btn" title="Zoom In" disabled={disableZoom}>🔍 +</button>
        <button onClick={(e) => {
          if (disableZoom) return
          const container = e.target.closest('.map-viewer').querySelector('.map-container')
          const rect = container.getBoundingClientRect()
          handleZoomAtPoint(-0.2, rect.left + rect.width / 2, rect.top + rect.height / 2, rect)
        }} className="zoom-btn" title="Zoom Out" disabled={disableZoom}>🔍 -</button>
        <button onClick={resetView} className="zoom-btn" title="Reset View" disabled={disableZoom}>🎯</button>
        <span className="zoom-level">{(zoom * 100).toFixed(0)}%</span>
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
          viewBox={`${-pan.x / (zoom * tileSize)} ${-pan.y / (zoom * tileSize)} ${gameMap.width / zoom} ${gameMap.height / zoom}`}
          onMouseDown={handleTileClick}
          style={{ transform: `scale(${zoom})`, transformOrigin: '0 0' }}
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
        ) : (
          <span>Scroll: Zoom | Shift+Drag or Middle Click: Pan</span>
        )}
      </div>
    </div>
  )
}
