import { useState } from 'react'
import './GameControls.css'

const UNIT_TYPES = ['warrior', 'tower', 'wall', 'naval_generator']

export default function GameControls({ state, playerId, onCommand }) {
  const [selectedUnitType, setSelectedUnitType] = useState('warrior')
  const [spawnX, setSpawnX] = useState('50')
  const [spawnY, setSpawnY] = useState('50')
  const [selectedUnitId, setSelectedUnitId] = useState('')
  const [moveX, setMoveX] = useState('51')
  const [moveY, setMoveY] = useState('50')

  const handleSpawnUnit = () => {
    if (!selectedUnitType || !spawnX || !spawnY) {
      alert('Please fill all spawn fields')
      return
    }

    onCommand({
      type: 'spawn_unit',
      data: {
        unitType: selectedUnitType,
        x: parseInt(spawnX),
        y: parseInt(spawnY)
      }
    })
  }

  const handleMoveUnit = () => {
    if (!selectedUnitId || !moveX || !moveY) {
      alert('Please select a unit and enter coordinates')
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

  const isMyTurn = state?.currentPlayerTurn === playerId

  const myUnits = state?.units ? 
    Object.values(state.units).filter(u => u.playerId === playerId) : 
    []

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

      <div className="controls-section">
        <h3>📍 Spawn Unit</h3>
        <div className="form-group">
          <label>Unit Type:</label>
          <select 
            value={selectedUnitType}
            onChange={(e) => setSelectedUnitType(e.target.value)}
            disabled={!isMyTurn}
          >
            {UNIT_TYPES.map(type => (
              <option key={type} value={type}>{type}</option>
            ))}
          </select>
        </div>

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
          onClick={handleSpawnUnit}
          className="btn-action"
          disabled={!isMyTurn}
        >
          Spawn Unit
        </button>
      </div>

      <div className="controls-section">
        <h3>⚡ Move Unit</h3>
        <div className="form-group">
          <label>Select Unit:</label>
          <select 
            value={selectedUnitId}
            onChange={(e) => setSelectedUnitId(e.target.value)}
            disabled={!isMyTurn || myUnits.length === 0}
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
          disabled={!isMyTurn || myUnits.length === 0}
        >
          Move Unit
        </button>
      </div>

      <div className="info-box">
        <h3>ℹ️ Info</h3>
        <ul>
          <li>⚔️ warrior: Mobile unit</li>
          <li>🏰 tower: Static defense</li>
          <li>🧱 wall: Blocker</li>
          <li>🌊 naval_generator: Special unit</li>
        </ul>
      </div>
    </div>
  )
}
