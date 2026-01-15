import { useState } from 'react'
import './GameBoard.css'

const UNIT_EMOJIS = {
  // Base
  main_base: '👑',
  
  // Estructuras
  tower: '🏰',
  wall: '🧱',
  land_generator: '🏞️',
  naval_generator: '🌊',
  
  // Unidades terrestres
  warrior: '⚔️',
  land_soldier: '🗡️',
  
  // Unidades navales
  naval_ship: '⛵',
  
  // Default
  default: '❓',
}

// ========== COLORES DE UNIDADES - EDITA AQUI ==========
const COLOR_CONFIG = {
  PLAYER_1_HUE: 200,        // Azul
  PLAYER_2_HUE: 0,          // Rojo
  SATURATION: 100,          // Saturación (0-100)
  LIGHTNESS_MIN: 15,        // Brillo mínimo (muy oscuro - casi negro)
  LIGHTNESS_MAX: 85,        // Brillo máximo (muy claro - casi blanco)
}

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

const getUnitColor = (playerId, unitType) => {
  const hue = playerId === 1 ? COLOR_CONFIG.PLAYER_1_HUE : COLOR_CONFIG.PLAYER_2_HUE
  const saturation = COLOR_CONFIG.SATURATION
  const intensity = getUnitColorIntensity(unitType)
  const lightness = COLOR_CONFIG.LIGHTNESS_MIN + (intensity * (COLOR_CONFIG.LIGHTNESS_MAX - COLOR_CONFIG.LIGHTNESS_MIN))
  return `hsl(${hue}, ${saturation}%, ${lightness}%)`
}

export default function GameBoard({ state, playerId, selectedUnitId, onSelectUnit }) {
  const [myUnitsCollapsed, setMyUnitsCollapsed] = useState(false)
  const [enemyUnitsCollapsed, setEnemyUnitsCollapsed] = useState(true)
  
  if (!state) {
    return <div className="game-board">Waiting for game state...</div>
  }
  
  // Separar unidades propias y enemigas
  const allUnits = state.units ? Object.values(state.units) : []
  const myUnits = allUnits.filter(u => u.playerId === playerId)
  const enemyUnits = allUnits.filter(u => u.playerId !== playerId)

  const getPhaseInfo = (phase) => {
    const phaseData = {
      'turn_start': { emoji: '🎬', color: '#4CAF50', label: 'Turn Start' },
      'preparation': { emoji: '🎴', color: '#2196F3', label: 'Preparation' },
      'battle': { emoji: '⚔️', color: '#FF5722', label: 'Battle!' },
      'turn_end': { emoji: '🏁', color: '#9E9E9E', label: 'Turn End' },
    }
    return phaseData[phase] || { emoji: '❓', color: '#666', label: 'Unknown' }
  }

  const renderGameInfo = () => {
    const phaseInfo = getPhaseInfo(state.currentPhase)
    const currentPlayer = state.currentPlayerTurn || 0
    const isHumanTurn = currentPlayer === state.humanPlayerId
    
    return (
      <div className="game-info">
        <div className="info-row">
          <span>📍 Tick: {state.tick || 0}</span>
          <span className="phase-indicator" style={{backgroundColor: phaseInfo.color}}>
            {phaseInfo.emoji} {phaseInfo.label}
          </span>
          <span>🎯 Turn: {state.turnNumber || 0}</span>
          <span className={isHumanTurn ? 'your-turn' : 'opponent-turn'}>
            {isHumanTurn ? '👤 Your Turn' : '🤖 AI Turn'}
          </span>
        </div>
      </div>
    )
  }

  const renderPlayers = () => {
    if (!state.players) {
      return <div className="players-section">No player data available</div>
    }

    const playerEntries = Object.entries(state.players)

    if (playerEntries.length === 0) {
      return <div className="players-section">No players have joined yet</div>
    }

    return (
      <div className="players-section">
        <h3>👥 Players</h3>
        {playerEntries.map(([id, player]) => (
          <div key={id} className={`player ${player.isAi ? 'ai' : 'human'}`}>
            <div className="player-name">
              {player.isAi ? '🤖' : '👤'} Player {id} {id == state.humanPlayerId ? '(Human)' : id == state.aiPlayerId ? '(AI)' : ''}
            </div>
            <div className="player-stats">
              {player.hand && <span>Hand: {Array.isArray(player.hand) ? player.hand.join(', ') : 'empty'}</span>}
              <span>Deck: {player.deckCount || 0}</span>
              <span className={player.connected ? 'connected' : 'disconnected'}>
                {player.connected ? '🟢 Online' : '🔴 Offline'}
              </span>
              <span>Ready: {player.ready ? '✅' : '❌'}</span>
            </div>
          </div>
        ))}
      </div>
    )
  }

  const getStatusEmoji = (status) => {
    const statusMap = {
      'idle': '⏸️',
      'moving': '🚶',
      'waiting': '⏳',
      'blocked': '🚫',
      'attacking': '⚡',
    }
    return statusMap[status] || '❓'
  }

  const getStatusColor = (status) => {
    const colorMap = {
      'idle': '#9E9E9E',
      'moving': '#2196F3',
      'waiting': '#FF9800',
      'blocked': '#F44336',
      'attacking': '#FF1744',
    }
    return colorMap[status] || '#666'
  }

  const renderUnitCard = (unit, isMine) => {
    const hpPercent = unit.maxHp ? (unit.hp / unit.maxHp) * 100 : 100
    const isAlive = unit.hp > 0
    const isSelected = selectedUnitId === unit.id
    const unitColor = getUnitColor(unit.playerId, unit.unitType)
    const statusColor = getStatusColor(unit.status)
    
    return (
      <div 
        key={unit.id} 
        className={`unit-card ${!isAlive ? 'dead' : ''} ${isSelected ? 'selected' : ''} ${isMine ? 'clickable' : ''} status-${unit.status}`}
        onClick={() => isMine && onSelectUnit && onSelectUnit(unit.id)}
        style={{
          borderLeftColor: unitColor,
          backgroundColor: `${unitColor}15`
        }}
      >
        <div className="unit-header">
          <span className="unit-emoji" style={{color: unitColor}}>{UNIT_EMOJIS[unit.unitType] || UNIT_EMOJIS.default}</span>
          <span className="unit-name">{unit.unitType}</span>
          {isSelected && <span className="selected-badge">✓</span>}
        </div>

        <div className="unit-status-bar" style={{backgroundColor: statusColor}}>
          <span className="status-emoji">{getStatusEmoji(unit.status)}</span>
          <span className="status-text">{unit.status}</span>
        </div>
        
        <div className="hp-bar-container">
          <div className="hp-bar" style={{width: `${hpPercent}%`, backgroundColor: hpPercent > 50 ? '#4CAF50' : hpPercent > 25 ? '#FF9800' : '#F44336'}}></div>
          <span className="hp-text">{unit.hp}/{unit.maxHp}</span>
        </div>
        
        <div className="unit-stats">
          <span>📍 ({unit.x}, {unit.y})</span>
          {unit.attackDamage > 0 && <span>⚔️ {unit.attackDamage}</span>}
          {unit.attackRange > 0 && <span>🎯 {unit.attackRange}</span>}
          {unit.detectionRange > 0 && <span>👁️ {unit.detectionRange}</span>}
        </div>
      </div>
    )
  }

  const renderUnits = () => {
    return (
      <div className="units-section">
        {/* My Units */}
        <div className="unit-group">
          <div 
            className="unit-group-header my-units-header"
            onClick={() => setMyUnitsCollapsed(!myUnitsCollapsed)}
          >
            <h3>👤 My Units ({myUnits.length})</h3>
            <span className="collapse-icon">{myUnitsCollapsed ? '▼' : '▲'}</span>
          </div>
          {!myUnitsCollapsed && (
            <div className="units-grid">
              {myUnits.length === 0 ? (
                <div className="no-units">No units yet. Spawn from cards!</div>
              ) : (
                myUnits.map(unit => renderUnitCard(unit, true))
              )}
            </div>
          )}
        </div>

        {/* Enemy Units */}
        <div className="unit-group">
          <div 
            className="unit-group-header enemy-units-header"
            onClick={() => setEnemyUnitsCollapsed(!enemyUnitsCollapsed)}
          >
            <h3>🤖 Enemy Units ({enemyUnits.length})</h3>
            <span className="collapse-icon">{enemyUnitsCollapsed ? '▼' : '▲'}</span>
          </div>
          {!enemyUnitsCollapsed && (
            <div className="units-grid">
              {enemyUnits.length === 0 ? (
                <div className="no-units">No enemy units</div>
              ) : (
                enemyUnits.map(unit => renderUnitCard(unit, false))
              )}
            </div>
          )}
        </div>
      </div>
    )
  }

  return (
    <div className="game-board">
      {renderGameInfo()}
      <div className="board-content">
        {renderPlayers()}
        {renderUnits()}
      </div>
    </div>
  )
}
