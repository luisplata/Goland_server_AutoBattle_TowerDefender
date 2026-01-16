import { useState, useRef, useEffect, useCallback } from 'react'
import './MapViewer.css'

const TERRAIN_COLORS = {
  0: '#4a7c3c', // Grass
  1: '#6b8e23', // Path
  2: '#1e5aa0', // Water
}

// Emojis para cada tipo de unidad (igual que la leyenda)
const UNIT_EMOJIS = {
  main_base: '👑',
  tower: '🏰',
  wall: '🧱',
  land_generator: '🏞️',
  naval_generator: '🌊',
  warrior: '⚔️',
  land_soldier: '🗡️',
  naval_ship: '⛵',
  default: '❓',
}

// Colores de fondo para los bandos
const TEAM_COLORS = {
  1: '#2196F3', // Azul para aliados (Player 1)
  2: '#F44336', // Rojo para enemigos (Player 2)
}

const getUnitEmoji = (unitType) => UNIT_EMOJIS[unitType] || UNIT_EMOJIS.default

const getTeamColor = (playerId) => TEAM_COLORS[playerId] || '#666'

export default function CanvasMapViewer({ gameMap, units, selectedTile, onSelectTile, disableZoom = false, playerId, selectedCard }) {
  const [zoom, setZoom] = useState(3)
  const [pan, setPan] = useState({ x: 0, y: 0 })
  const [isPanning, setIsPanning] = useState(false)
  const [startPan, setStartPan] = useState({ x: 0, y: 0 })
  const [selectedUnitId, setSelectedUnitId] = useState(null)
  const containerRef = useRef(null)
  const canvasRef = useRef(null)

  const tileSize = 3 // Base tile size in world units; zoom starts at 3 for perceived 100%
  const MAX_ZOOM = 6 // Máximo zoom permitido
  const MIN_ZOOM = 0.5 // Mínimo zoom permitido

  // Calcula los límites de pan para mantener el mapa visible
  const getClampedPan = (panX, panY, zoomValue) => {
    if (!containerRef.current || !gameMap) return { x: panX, y: panY }
    
    const rect = containerRef.current.getBoundingClientRect()
    const viewportWidth = rect.width
    const viewportHeight = rect.height
    
    // Dimensiones del mapa en pixels
    const mapWidth = gameMap.width * tileSize * zoomValue
    const mapHeight = gameMap.height * tileSize * zoomValue
    
    // Limitar pan para no salir demasiado del mapa
    let clampedX = panX
    let clampedY = panY
    
    // Si el mapa es más grande que el viewport, permitir pan controlado
    if (mapWidth > viewportWidth) {
      clampedX = Math.max(-mapWidth + 50, Math.min(viewportWidth - 50, panX))
    } else {
      clampedX = Math.max(viewportWidth - mapWidth, Math.min(0, panX))
    }
    
    if (mapHeight > viewportHeight) {
      clampedY = Math.max(-mapHeight + 50, Math.min(viewportHeight - 50, panY))
    } else {
      clampedY = Math.max(viewportHeight - mapHeight, Math.min(0, panY))
    }
    
    return { x: clampedX, y: clampedY }
  }

  // Controlled area calculation (same logic as SVG component)
  const controlledArea = (() => {
    if (!playerId || !units || !selectedCard) return new Set()
    const controlled = new Set()
    const myUnits = Object.values(units).filter(u => u.playerId === playerId && u.hp > 0)

    for (const unit of myUnits) {
      if (!unit.buildRange || unit.buildRange <= 0) continue
      for (let dy = -unit.buildRange; dy <= unit.buildRange; dy++) {
        for (let dx = -unit.buildRange; dx <= unit.buildRange; dx++) {
          const manhattan = Math.abs(dx) + Math.abs(dy)
          if (manhattan <= unit.buildRange) {
            const tx = unit.x + dx
            const ty = unit.y + dy
            if (tx >= 0 && tx < gameMap.width && ty >= 0 && ty < gameMap.height) {
              controlled.add(`${tx},${ty}`)
            }
          }
        }
      }
    }
    return controlled
  })()

  const isStructureCard = selectedCard && ['tower', 'wall', 'land_generator', 'naval_generator'].includes(selectedCard)

  // Zoom with wheel at mouse position
  const handleWheel = useCallback((e) => {
    e.preventDefault()
    if (disableZoom) return
    const container = containerRef.current
    if (!container) return

    const rect = container.getBoundingClientRect()
    const delta = e.deltaY > 0 ? -0.2 : 0.2
    const mouseX = e.clientX - rect.left
    const mouseY = e.clientY - rect.top

    setZoom(oldZoom => {
      const newZoom = Math.max(MIN_ZOOM, Math.min(MAX_ZOOM, oldZoom + delta))
      if (oldZoom === newZoom) return oldZoom
      
      // Calcular la posición del mundo antes del zoom
      const worldXBefore = (mouseX - pan.x) / oldZoom
      const worldYBefore = (mouseY - pan.y) / oldZoom
      
      // Calcular el nuevo pan para mantener el punto del mundo bajo el mouse
      setPan(prevPan => {
        let newPan = {
          x: mouseX - worldXBefore * newZoom,
          y: mouseY - worldYBefore * newZoom,
        }
        newPan = getClampedPan(newPan.x, newPan.y, newZoom)
        return newPan
      })
      
      return newZoom
    })
  }, [disableZoom, pan])

  // Zoom via buttons at center
  const handleZoomAtCenter = useCallback((delta) => {
    if (disableZoom) return
    const container = containerRef.current
    if (!container) return
    const rect = container.getBoundingClientRect()
    const cx = rect.width / 2
    const cy = rect.height / 2
    
    setZoom(oldZoom => {
      const newZoom = Math.max(MIN_ZOOM, Math.min(MAX_ZOOM, oldZoom + delta))
      if (oldZoom === newZoom) return oldZoom
      
      // Calcular la posición del mundo antes del zoom
      const worldXBefore = (cx - pan.x) / oldZoom
      const worldYBefore = (cy - pan.y) / oldZoom
      
      // Calcular el nuevo pan para mantener el punto del mundo en el centro
      setPan(prevPan => {
        let newPan = {
          x: cx - worldXBefore * newZoom,
          y: cy - worldYBefore * newZoom,
        }
        newPan = getClampedPan(newPan.x, newPan.y, newZoom)
        return newPan
      })
      
      return newZoom
    })
  }, [disableZoom, pan])

  useEffect(() => {
    const container = containerRef.current
    if (!container) return
    container.addEventListener('wheel', handleWheel, { passive: false })
    return () => {
      container.removeEventListener('wheel', handleWheel)
    }
  }, [handleWheel])

  // Mouse handlers for pan and click
  const handleMouseDown = (e) => {
    if (e.button === 1 || e.shiftKey) {
      setIsPanning(true)
      setStartPan({ x: e.clientX - pan.x, y: e.clientY - pan.y })
      e.preventDefault()
    }
  }

  const handleMouseMove = (e) => {
    if (isPanning) {
      let newPan = { x: e.clientX - startPan.x, y: e.clientY - startPan.y }
      newPan = getClampedPan(newPan.x, newPan.y, zoom)
      setPan(newPan)
    }
  }

  const handleMouseUp = () => setIsPanning(false)

  const resetView = () => {
    setZoom(3)
    setPan({ x: 0, y: 0 })
  }

  const pickTileFromEvent = (e) => {
    const canvas = canvasRef.current
    const container = containerRef.current
    if (!canvas || !container) return null
    const rect = container.getBoundingClientRect()
    const worldX = (e.clientX - rect.left - pan.x) / zoom
    const worldY = (e.clientY - rect.top - pan.y) / zoom
    const x = Math.floor(worldX / tileSize)
    const y = Math.floor(worldY / tileSize)
    if (x < 0 || y < 0 || x >= gameMap.width || y >= gameMap.height) return null
    const tile = gameMap.tiles[y][x]
    return { x, y, walkable: tile.walkable }
  }

  const handleClick = (e) => {
    if (isPanning || e.button !== 0) return
    const tile = pickTileFromEvent(e)
    if (tile && onSelectTile) onSelectTile(tile)
    
    // Check if there's a unit at this position
    if (tile && units) {
      const unitAtTile = Object.values(units).find(u => u.x === tile.x && u.y === tile.y)
      if (unitAtTile) {
        setSelectedUnitId(unitAtTile.id)
      } else {
        setSelectedUnitId(null)
      }
    }
  }

  // Draw everything
  const draw = useCallback(() => {
    const canvas = canvasRef.current
    if (!canvas || !gameMap) return
    const ctx = canvas.getContext('2d')
    ctx.setTransform(1, 0, 0, 1, 0, 0)
    ctx.clearRect(0, 0, canvas.width, canvas.height)

    // World transform: zoom + pan
    ctx.setTransform(zoom, 0, 0, zoom, pan.x, pan.y)

    // Draw tiles
    for (let y = 0; y < gameMap.height; y++) {
      const row = gameMap.tiles[y]
      for (let x = 0; x < gameMap.width; x++) {
        const tile = row[x]
        ctx.fillStyle = TERRAIN_COLORS[tile.terrainId] || '#666'
        ctx.globalAlpha = tile.walkable ? 1 : 0.6
        ctx.fillRect(x * tileSize, y * tileSize, tileSize, tileSize)
      }
    }
    ctx.globalAlpha = 1

    // Controlled area overlay
    if (selectedCard && controlledArea.size > 0) {
      ctx.lineWidth = 0.04 * tileSize
      for (const key of controlledArea) {
        const [tx, ty] = key.split(',').map(Number)
        const isControlled = true
        const isValidForStructure = !isStructureCard || isControlled
        const strokeColor = isValidForStructure ? '#00ff88' : '#ff4444'
        const fillAlpha = isStructureCard && !isControlled ? 0.3 : 0
        if (fillAlpha > 0) {
          ctx.fillStyle = `rgba(0,255,136,${fillAlpha})`
          ctx.fillRect(tx * tileSize, ty * tileSize, tileSize, tileSize)
        }
        ctx.strokeStyle = strokeColor
        ctx.strokeRect(tx * tileSize, ty * tileSize, tileSize, tileSize)
      }
    }

    // Selection highlight
    if (selectedTile) {
      ctx.strokeStyle = selectedTile.walkable ? '#00ff88' : '#ff4444'
      ctx.lineWidth = 0.1 * tileSize
      ctx.strokeRect(selectedTile.x * tileSize, selectedTile.y * tileSize, tileSize, tileSize)
    }

    // Draw targeting lines for selected unit only
    if (units && selectedUnitId && units[selectedUnitId]) {
      const unit = units[selectedUnitId]
      if (unit.targetId && units[unit.targetId]) {
        const target = units[unit.targetId]
        const cx = (unit.x + 0.5) * tileSize
        const cy = (unit.y + 0.5) * tileSize
        const targetCx = (target.x + 0.5) * tileSize
        const targetCy = (target.y + 0.5) * tileSize
        
        ctx.strokeStyle = '#ff00ff'
        ctx.lineWidth = 0.2 * tileSize
        ctx.globalAlpha = 1.0
        ctx.beginPath()
        ctx.moveTo(cx, cy)
        ctx.lineTo(targetCx, targetCy)
        ctx.stroke()
        ctx.globalAlpha = 1
      }
    }

    // Units
    if (units) {
      Object.values(units).forEach(unit => {
        const teamColor = getTeamColor(unit.playerId)
        const emoji = getUnitEmoji(unit.unitType)
        const cx = (unit.x + 0.5) * tileSize
        const cy = (unit.y + 0.5) * tileSize
        const r = 0.5 * tileSize
        const hpPercent = unit.maxHp ? unit.hp / unit.maxHp : 1
        const isSelected = unit.id === selectedUnitId

        // Detection Range (área de cambio de target)
        if (unit.detectionRange > 0) {
          ctx.strokeStyle = isSelected ? '#ffff00' : teamColor
          ctx.lineWidth = isSelected ? 0.08 * tileSize : 0.02 * tileSize
          ctx.setLineDash([0.3 * tileSize, 0.15 * tileSize])
          ctx.globalAlpha = isSelected ? 0.6 : 0.1
          ctx.beginPath()
          ctx.arc(cx, cy, unit.detectionRange * tileSize, 0, Math.PI * 2)
          ctx.stroke()
          ctx.setLineDash([])
          ctx.globalAlpha = 1
        }

        // Attack Range (área de ataque)
        if (unit.attackDamage > 0 && unit.attackRange > 0) {
          ctx.strokeStyle = isSelected ? '#ff4444' : teamColor
          ctx.lineWidth = isSelected ? 0.08 * tileSize : 0.02 * tileSize
          ctx.globalAlpha = isSelected ? 0.5 : 0.15
          ctx.beginPath()
          ctx.arc(cx, cy, unit.attackRange * tileSize, 0, Math.PI * 2)
          ctx.stroke()
          ctx.globalAlpha = 1
        }

        // Build Range
        if (unit.buildRange > 0) {
          ctx.strokeStyle = '#00ff88'
          ctx.lineWidth = 0.03 * tileSize
          ctx.setLineDash([0.2 * tileSize, 0.2 * tileSize])
          ctx.globalAlpha = 0.15
          ctx.beginPath()
          ctx.arc(cx, cy, unit.buildRange * tileSize, 0, Math.PI * 2)
          ctx.stroke()
          ctx.setLineDash([])
          ctx.globalAlpha = 1
        }

        // Círculo de fondo con color del equipo
        ctx.fillStyle = teamColor
        ctx.globalAlpha = 0.9
        ctx.beginPath()
        ctx.arc(cx, cy, r, 0, Math.PI * 2)
        ctx.fill()
        ctx.globalAlpha = 1

        // Borde del círculo (más grueso si está seleccionada)
        ctx.strokeStyle = isSelected ? '#ffff00' : '#fff'
        ctx.lineWidth = isSelected ? 0.15 * tileSize : 0.08 * tileSize
        ctx.beginPath()
        ctx.arc(cx, cy, r, 0, Math.PI * 2)
        ctx.stroke()

        // Dibujar emoji en el centro
        ctx.save()
        ctx.textAlign = 'center'
        ctx.textBaseline = 'middle'
        ctx.font = `${r * 1.4}px Arial`
        ctx.fillStyle = '#fff'
        ctx.strokeStyle = '#000'
        ctx.lineWidth = 0.03 * tileSize
        ctx.strokeText(emoji, cx, cy)
        ctx.fillText(emoji, cx, cy)
        ctx.restore()

        // HP bar
        const barW = 0.8 * tileSize
        const barH = 0.1 * tileSize
        const barY = cy + r + 0.15 * tileSize
        ctx.fillStyle = '#1a1a1a'
        ctx.fillRect(cx - barW / 2, barY, barW, barH)
        ctx.fillStyle = hpPercent > 0.5 ? '#4CAF50' : hpPercent > 0.25 ? '#FF9800' : '#F44336'
        ctx.fillRect(cx - barW / 2, barY, barW * hpPercent, barH)
        
        // Borde de HP bar
        ctx.strokeStyle = '#fff'
        ctx.lineWidth = 0.02 * tileSize
        ctx.strokeRect(cx - barW / 2, barY, barW, barH)
      })
    }
  }, [gameMap, units, controlledArea, isStructureCard, selectedTile, selectedUnitId, pan, zoom, tileSize])

  // Resize canvas to container size
  useEffect(() => {
    const canvas = canvasRef.current
    const container = containerRef.current
    if (!canvas || !container) return

    const resize = () => {
      const { clientWidth, clientHeight } = container
      canvas.width = clientWidth
      canvas.height = clientHeight
      draw()
    }

    resize()
    const observer = new ResizeObserver(resize)
    observer.observe(container)
    return () => observer.disconnect()
  }, [draw])


  // Redraw on changes
  useEffect(() => {
    draw()
  }, [draw])

  if (!gameMap || !gameMap.tiles) {
    return <div className="map-viewer">Loading map...</div>
  }

  return (
    <div className="map-viewer">
      <h3>🗺️ Game Map (Canvas)</h3>

      <div className="map-controls">
        <button onClick={() => handleZoomAtCenter(0.2)} className="zoom-btn" title="Zoom In" disabled={disableZoom}>🔍 +</button>
        <button onClick={() => handleZoomAtCenter(-0.2)} className="zoom-btn" title="Zoom Out" disabled={disableZoom}>🔍 -</button>
        <button onClick={resetView} className="zoom-btn" title="Reset View" disabled={disableZoom}>🎯</button>
        <span className="zoom-level">{((zoom / 3) * 100).toFixed(0)}%</span>
      </div>

      <div
        className="map-container"
        ref={containerRef}
        onMouseDown={handleMouseDown}
        onMouseMove={handleMouseMove}
        onMouseUp={handleMouseUp}
        onMouseLeave={handleMouseUp}
        onClick={handleClick}
        onContextMenu={(e) => e.preventDefault()}
        style={{ cursor: isPanning ? 'grabbing' : 'crosshair' }}
      >
        <canvas ref={canvasRef} className="map-canvas" />
      </div>

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
          <div style={{ fontSize: '0.85rem', marginBottom: '0.3rem', fontWeight: 'bold' }}>Team Colors:</div>
          <div className="legend-item" style={{ fontSize: '0.8rem' }}>
            <div style={{ width: '16px', height: '16px', borderRadius: '50%', backgroundColor: getTeamColor(1), border: '2px solid white', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '10px' }}>⚔️</div>
            <span>Allies (Blue)</span>
          </div>
          <div className="legend-item" style={{ fontSize: '0.8rem' }}>
            <div style={{ width: '16px', height: '16px', borderRadius: '50%', backgroundColor: getTeamColor(2), border: '2px solid white', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '10px' }}>⚔️</div>
            <span>Enemies (Red)</span>
          </div>
        </div>
      </div>

      <div className="selection-info">
        {selectedTile ? (
          <span>
            Selected: ({selectedTile.x}, {selectedTile.y}) — {selectedTile.walkable ? 'Walkable ✅' : 'Water ❌'}
          </span>
        ) : selectedCard && controlledArea.size > 0 ? (
          <span>Controlled area shown in green for structures</span>
        ) : (
          <span>Click a tile to select it</span>
        )}
      </div>

      {selectedUnitId && units && units[selectedUnitId] && (
        <div style={{
          marginTop: '0.6rem',
          padding: '0.6rem 0.8rem',
          background: 'rgba(0,0,0,0.25)',
          border: '1px solid #00ff88',
          borderRadius: 6,
          display: 'inline-block'
        }}>
          {(() => {
            const u = units[selectedUnitId]
            return (
              <div style={{ display: 'grid', gridTemplateColumns: 'auto auto', gap: '0.35rem 1rem', alignItems: 'center' }}>
                <div style={{ gridColumn: '1 / -1', fontWeight: 'bold' }}>
                  Unidad: {u.unitType} {u.playerId === playerId ? '(Aliada)' : '(Enemiga)'}
                </div>
                <div>HP:</div>
                <div>{u.hp}/{u.maxHp}</div>
                <div>DMG:</div>
                <div>{u.attackDamage}</div>
                <div>Rango Ataque:</div>
                <div>{u.attackRange}</div>
                <div>Rango Detección:</div>
                <div>{u.detectionRange}</div>
                {u.targetId && units[u.targetId] && (
                  <>
                    <div>Target:</div>
                    <div style={{ color: '#00ffff' }}>{units[u.targetId].unitType} (ID: {u.targetId})</div>
                  </>
                )}
                {u.spawnedById && units[u.spawnedById] && (
                  <>
                    <div>Generada por:</div>
                    <div>{units[u.spawnedById].unitType}</div>
                  </>
                )}
                {u.generatedUnitType && (
                  <>
                    <div>Genera:</div>
                    <div>{u.generatedUnitType}</div>
                  </>
                )}
                {u.status && (
                  <>
                    <div>Estado:</div>
                    <div>{u.status}</div>
                  </>
                )}
              </div>
            )
          })()}
        </div>
      )}
    </div>
  )
}
