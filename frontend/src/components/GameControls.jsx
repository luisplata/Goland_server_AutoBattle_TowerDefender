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
  const [moveX, setMoveX] = useState('51')
  const [moveY, setMoveY] = useState('50')

  const currentPhase = state?.currentPhase
  const isBaseSelectionPhase = currentPhase === 'base_selection'
  const hasPlacedBase = playerId === state?.humanPlayerId ? (state?.humanBaseId > 0) : (state?.aiBaseId > 0)

  // Obtener la mano del jugador actual
  const myPlayer = state?.players?.[playerId]
  const myHand = myPlayer?.hand || []

  // Obtener unidades propias
  const myUnits = state?.units ? 
    Object.values(state.units).filter(u => u.playerId === playerId) : 
    []

  const isMyTurn = state?.currentPlayerTurn === playerId

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

  const handleMoveUnit = () => {
    if (!selectedUnitId || !moveX || !moveY) {
      alert('Selecciona una unidad y coordenadas')
      return
    }

    onCommand({
      type: 'move_unit',
      data: {
        unitId: parseInt(selectedUnitId),
        x: parseInt(moveX),
        y: parseInt(moveY)
      }
    })
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
            <button onClick={handleEndTurn} className="btn-end-turn">
              End Turn
            </button>
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
        <div className="hand-grid">
          {myHand.length === 0 ? (
            <div className="empty-hand">No cards in hand</div>
          ) : (
            myHand.map((card, index) => (
              <div
                key={index}
                className={`card ${selectedCard === card ? 'selected' : ''}`}
                onClick={() => onSelectCard(selectedCard === card ? null : card)}
                style={{ pointerEvents: isMyTurn ? 'auto' : 'none', opacity: isMyTurn ? 1 : 0.5 }}
              >
                <div className="card-emoji">{CARD_EMOJIS[card] || '?'}</div>
                <div className="card-name">{card}</div>
              </div>
            ))
          )}
        </div>
      </div>

      {/* SPAWN FROM CARD SECTION */}
      {selectedCard && (
        <div className="controls-section spawn-section">
          <h3>📍 Spawn {selectedCard}</h3>
          <div className="form-row">
            <div className="form-group">
              <label>X:</label>
              <input 
                type="number" 
                value={spawnX}
                onChange={(e) => setSpawnX(e.target.value)}
                disabled={!isMyTurn}
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
                disabled={!isMyTurn}
                min="0"
                max="100"
              />
            </div>
          </div>

          <div style={{ display: 'flex', gap: '0.8rem', alignItems: 'center' }}>
            <button 
              onClick={handleSpawnFromCard}
              className="btn-action btn-spawn"
              disabled={!isMyTurn || (selectedTile && !selectedTile.walkable)}
            >
              {CARD_EMOJIS[selectedCard]} Spawn {selectedCard}
            </button>
            {selectedTile && (
              <span style={{ opacity: 0.8 }}>
                {selectedTile.walkable ? '✅ Posición válida' : '❌ Posición inválida (agua)'}
              </span>
            )}
          </div>
        </div>
      )}

      {/* MOVE UNIT SECTION */}
      <div className="controls-section">
        <h3>⚡ Move Unit</h3>
        {myUnits.length === 0 ? (
          <div className="empty-units">No units to move</div>
        ) : (
          <>
            {selectedUnitId ? (
              <div className="selected-unit-info">
                <span>Selected: {myUnits.find(u => u.id === selectedUnitId)?.unitType || 'Unknown'} (ID: {selectedUnitId})</span>
                <span className="help-text">Click unit in "My Units" to change selection</span>
              </div>
            ) : (
              <div className="help-text">👆 Click a unit in "My Units" section above to select it</div>
            )}

            <div className="form-row">
              <div className="form-group">
                <label>Target X:</label>
                <input 
                  type="number" 
                  value={moveX}
                  onChange={(e) => setMoveX(e.target.value)}
                  disabled={!isMyTurn}
                  min="0"
                  max="100"
                />
              </div>
              <div className="form-group">
                <label>Target Y:</label>
                <input 
                  type="number" 
                  value={moveY}
                  onChange={(e) => setMoveY(e.target.value)}
                  disabled={!isMyTurn}
                  min="0"
                  max="100"
                />
              </div>
            </div>

            <button 
              onClick={handleMoveUnit}
              className="btn-action"
              disabled={!isMyTurn || !selectedUnitId}
            >
              Move Unit
            </button>
          </>
        )}
      </div>

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
