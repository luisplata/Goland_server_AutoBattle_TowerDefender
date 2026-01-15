import { useState } from 'react'
import './GameControls.css'

const CARD_EMOJIS = {
  warrior: '⚔️',
  tower: '🏰',
  wall: '🧱',
  naval_generator: '🌊',
  land_generator: '🏞️',
}

export default function GameControls({ state, playerId, onCommand }) {
  const [selectedCard, setSelectedCard] = useState(null)
  const [spawnX, setSpawnX] = useState('50')
  const [spawnY, setSpawnY] = useState('50')
  const [selectedUnitId, setSelectedUnitId] = useState('')
  const [moveX, setMoveX] = useState('51')
  const [moveY, setMoveY] = useState('50')

  // Obtener la mano del jugador actual
  const myPlayer = state?.players?.[playerId]
  const myHand = myPlayer?.hand || []

  // Obtener unidades propias
  const myUnits = state?.units ? 
    Object.values(state.units).filter(u => u.playerId === playerId) : 
    []

  const isMyTurn = state?.currentPlayerTurn === playerId

  const handleSpawnFromCard = () => {
    if (!selectedCard || !spawnX || !spawnY) {
      alert('Selecciona una carta y coordenadas')
      return
    }

    onCommand({
      type: 'spawn_unit',
      data: {
        unitType: selectedCard,
        x: parseInt(spawnX),
        y: parseInt(spawnY)
      }
    })
    
    setSelectedCard(null)
    setSpawnX('50')
    setSpawnY('50')
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
                onClick={() => setSelectedCard(selectedCard === card ? null : card)}
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

          <button 
            onClick={handleSpawnFromCard}
            className="btn-action btn-spawn"
            disabled={!isMyTurn}
          >
            {CARD_EMOJIS[selectedCard]} Spawn {selectedCard}
          </button>
        </div>
      )}

      {/* MOVE UNIT SECTION */}
      <div className="controls-section">
        <h3>⚡ Move Unit</h3>
        {myUnits.length === 0 ? (
          <div className="empty-units">No units to move</div>
        ) : (
          <>
            <div className="form-group">
              <label>Select Unit:</label>
              <select 
                value={selectedUnitId}
                onChange={(e) => setSelectedUnitId(e.target.value)}
                disabled={!isMyTurn}
              >
                <option value="">-- Choose a unit --</option>
                {myUnits.map(unit => (
                  <option key={unit.id} value={unit.id}>
                    {unit.unitType} (ID: {unit.id}) @ ({unit.x}, {unit.y})
                  </option>
                ))}
              </select>
            </div>

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
