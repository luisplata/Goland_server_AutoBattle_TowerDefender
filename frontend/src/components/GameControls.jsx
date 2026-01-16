import { useState, useEffect } from 'react'
import './GameControls.css'

const CARD_EMOJIS = {
  warrior: '⚔️',
  tower: '🏰',
  wall: '🧱',
  naval_generator: '🌊',
  land_generator: '🏞️',
}

export default function GameControls({ state, playerId, onCommand, selectedTile, gameMap, onClearSelection, selectedUnitId, selectedCard, onSelectCard }) {
  const [spawnX, setSpawnX] = useState('50')
  const [spawnY, setSpawnY] = useState('50')
  const [unitStats, setUnitStats] = useState({})

  // Cargar estadísticas de unidades al montar el componente
  useEffect(() => {
    fetch('http://localhost:8080/unit-stats')
      .then(res => res.json())
      .then(data => setUnitStats(data))
      .catch(err => console.error('Error loading unit stats:', err))
  }, [])

  const currentPhase = state?.currentPhase
  const isBaseSelectionPhase = currentPhase === 'base_selection'
  const hasPlacedBase = playerId === state?.humanPlayerId ? (state?.humanBaseId > 0) : (state?.aiBaseId > 0)
  const canPlayCards = currentPhase === 'preparation'

  // Limpiar carta seleccionada cuando sale de la fase de preparation
  useEffect(() => {
    if (!canPlayCards && selectedCard) {
      onSelectCard(null)
    }
  }, [canPlayCards, selectedCard, onSelectCard])

  // Obtener la mano del jugador actual
  const myPlayer = state?.players?.[playerId]
  const myHand = myPlayer?.hand || []

  // Obtener unidades propias
  const myUnits = state?.units ? 
    Object.values(state.units).filter(u => u.playerId === playerId) : 
    []

  // Permitir acciones simultáneas en preparation: ambos pueden actuar
  const isMyTurn = currentPhase === 'preparation' || state?.currentPlayerTurn === playerId
  const showEndTurn = currentPhase === 'preparation'

  // Área controlada (para validar estructuras en el cliente)
  const controlledArea = (() => {
    if (!playerId || !state?.units || !selectedCard) return new Set()
    const set = new Set()
    Object.values(state.units).forEach(u => {
      if (u.playerId !== playerId || !u.buildRange || u.buildRange <= 0 || u.hp <= 0) return
      for (let dy = -u.buildRange; dy <= u.buildRange; dy++) {
        for (let dx = -u.buildRange; dx <= u.buildRange; dx++) {
          const manhattan = Math.abs(dx) + Math.abs(dy)
          if (manhattan <= u.buildRange) {
            const tx = u.x + dx
            const ty = u.y + dy
            if (tx >= 0 && tx < (gameMap?.width ?? 0) && ty >= 0 && ty < (gameMap?.height ?? 0)) {
              set.add(`${tx},${ty}`)
            }
          }
        }
      }
    })
    return set
  })()

  const isStructureCard = selectedCard && ['tower', 'wall', 'land_generator', 'naval_generator'].includes(selectedCard)
  const isTileControlled = selectedTile ? controlledArea.has(`${selectedTile.x},${selectedTile.y}`) : false
  const isWalkable = selectedTile ? selectedTile.walkable : true
  const isValidPlacement = isWalkable && (!isStructureCard || isTileControlled)

  // Sync spawn coordinates with selected tile from map
  useEffect(() => {
    if (selectedTile) {
      setSpawnX(String(selectedTile.x))
      setSpawnY(String(selectedTile.y))
    }
  }, [selectedTile])

  const handleSpawnFromCard = () => {
    if (!selectedCard || !spawnX || !spawnY) {
      alert('Selecciona una carta y coordenadas')
      return
    }

    const x = parseInt(spawnX)
    const y = parseInt(spawnY)
    const walkable = selectedTile ? selectedTile.walkable : (gameMap?.tiles?.[y]?.[x]?.walkable ?? true)
    if (!walkable) {
      alert('La posición seleccionada no es válida (agua)')
      return
    }

    onCommand({
      type: 'spawn_unit',
      data: {
        unitType: selectedCard,
        x,
        y
      }
    })
    
    onSelectCard(null)
    setSpawnX('50')
    setSpawnY('50')
    
    // Clear the tile selection after spawning
    if (onClearSelection) {
      onClearSelection()
    }
  }

  const handleEndTurn = () => {
    onCommand({
      type: 'end_turn'
    })
  }

  const handlePlaceBase = () => {
    if (!selectedTile) {
      alert('Click on the map to select a position for your base')
      return
    }

    const x = parseInt(spawnX)
    const y = parseInt(spawnY)
    
    if (!selectedTile.walkable) {
      alert('Cannot place base on water')
      return
    }

    onCommand({
      type: 'place_base',
      data: { x, y }
    })
  }

  // Si estamos en fase de selección de base
  if (isBaseSelectionPhase) {
    return (
      <div className="game-controls">
        <h2>🏰 Place Your Base</h2>
        
        {hasPlacedBase ? (
          <div className="base-placed-info">
            <p>✅ Base placed! Waiting for opponent...</p>
          </div>
        ) : (
          <>
            <div className="base-selection-info">
              <p>📍 Click on the map to select a position for your main base</p>
              <p className="help-text">Your base will generate warriors automatically</p>
            </div>

            {selectedTile && (
              <div className="selected-position">
                <p>Selected: ({selectedTile.x}, {selectedTile.y})</p>
                <p className={selectedTile.walkable ? 'valid' : 'invalid'}>
                  {selectedTile.walkable ? '✅ Valid position' : '❌ Invalid (water)'}
                </p>
              </div>
            )}

            <button 
              onClick={handlePlaceBase}
              className="btn-action btn-place-base"
              disabled={!selectedTile || !selectedTile.walkable}
            >
              🏰 Place Base Here
            </button>
          </>
        )}
      </div>
    )
  }

  return (
    <div className="game-controls">
      <h2>🎮 Controls</h2>

      <div className="turn-status">
        {isMyTurn ? (
          <div className="my-turn">
            <span>✅ Your Turn</span>
            {showEndTurn && (
              <button onClick={handleEndTurn} className="btn-end-turn">
                End Turn
              </button>
            )}
          </div>
        ) : (
          <div className="opponent-turn">
            <span>⏳ Opponent Turn</span>
          </div>
        )}
      </div>

      {/* HAND SECTION */}
      <div className="controls-section hand-section">
        <h3>🎴 Hand ({myHand.length})</h3>
        {!canPlayCards && (
          <div className="help-text" style={{ color: '#ff9800' }}>Cards disabled outside preparation phase</div>
        )}
        <div className="hand-grid">
          {myHand.length === 0 ? (
            <div className="empty-hand">No cards in hand</div>
          ) : (
            myHand.map((card, index) => {
              const stats = unitStats[card]
              return (
                <div
                  key={index}
                  className={`card ${selectedCard === card ? 'selected' : ''}`}
                  onClick={() => {
                    if (!canPlayCards || !isMyTurn) return
                    onSelectCard(selectedCard === card ? null : card)
                  }}
                  style={{ pointerEvents: canPlayCards && isMyTurn ? 'auto' : 'none', opacity: canPlayCards && isMyTurn ? 1 : 0.5 }}
                  title={stats ? `HP: ${stats.hp} | DMG: ${stats.attackDamage} | Range: ${stats.attackRange}` : ''}
                >
                  <div className="card-emoji">{CARD_EMOJIS[card] || '?'}</div>
                  <div className="card-name">{card}</div>
                  {stats && (
                    <div className="card-stats">
                      <div className="stat">❤️ {stats.hp}</div>
                      {stats.attackDamage > 0 && <div className="stat">⚔️ {stats.attackDamage}</div>}
                    </div>
                  )}
                </div>
              )
            })
          )}
        </div>
      </div>

      {/* SPAWN FROM CARD SECTION */}
      {selectedCard && canPlayCards && (
        <div className="controls-section spawn-section">
          <h3>📍 Spawn {selectedCard}</h3>
          <div className="form-row">
            <div className="form-group">
              <label>X:</label>
              <input 
                type="number" 
                value={spawnX}
                onChange={(e) => setSpawnX(e.target.value)}
                disabled={!canPlayCards || !isMyTurn}
                min="0"
                max="100"
              />
            </div>
            <div className="form-group">
              <label>Y:</label>
              <input 
                type="number" 
                value={spawnY}
                onChange={(e) => setSpawnY(e.target.value)}
                disabled={!canPlayCards || !isMyTurn}
                min="0"
                max="100"
              />
            </div>
          </div>

          <div style={{ display: 'flex', gap: '0.8rem', alignItems: 'center' }}>
            <button 
              onClick={handleSpawnFromCard}
              className="btn-action btn-spawn"
              disabled={!canPlayCards || !isMyTurn || !isValidPlacement}
            >
              {CARD_EMOJIS[selectedCard]} Spawn {selectedCard}
            </button>
            {selectedTile && (
              <span style={{ opacity: 0.8 }}>
                {isValidPlacement ? '✅ Posición válida' : isStructureCard && !isTileControlled ? '❌ Fuera de área controlada' : '❌ Posición inválida (agua)'}
              </span>
            )}
          </div>
        </div>
      )}

      {selectedCard && !canPlayCards && (
        <div className="controls-section spawn-section" style={{ opacity: 0.6 }}>
          <h3>📍 Spawn {selectedCard}</h3>
          <div className="help-text" style={{ color: '#ff9800' }}>Cards can only be played in preparation phase</div>
        </div>
      )}

      <div className="info-box">
        <h3>ℹ️ Card Info</h3>
        <ul>
          <li>⚔️ warrior: Mobile unit</li>
          <li>🏰 tower: Static defense</li>
          <li>🧱 wall: Blocker</li>
          <li>🌊 naval_generator: Special unit</li>
          <li>🏞️ land_generator: Land unit</li>
        </ul>
      </div>
    </div>
  )
}
