import './GameStatusPanel.css'
import { useState } from 'react'

export default function GameStatusPanel({ state, playerId, gameId }) {
  const [isExpanded, setIsExpanded] = useState(true)

  if (!state) {
    return null
  }

  const getPhaseInfo = (phase) => {
    const phaseData = {
      'base_selection': { emoji: '🏰', color: '#9C27B0', label: 'Base Selection', duration: 0 },
      'turn_start': { emoji: '🎬', color: '#4CAF50', label: 'Turn Start', duration: state?.config?.turnStartDuration || 15 },
      'preparation': { emoji: '🎴', color: '#2196F3', label: 'Preparation', duration: state?.config?.preparationDuration || 150 },
      'battle': { emoji: '⚔️', color: '#FF5722', label: 'Battle!', duration: state?.config?.battleDuration || 300 },
      'turn_end': { emoji: '🏁', color: '#9E9E9E', label: 'Turn End', duration: state?.config?.turnEndDuration || 15 },
    }
    return phaseData[phase] || { emoji: '❓', color: '#666', label: 'Unknown', duration: 0 }
  }

  const phaseInfo = getPhaseInfo(state.currentPhase)
  const currentPlayer = state.currentPlayerTurn || 0
  const isHumanTurn = currentPlayer === state.humanPlayerId
  const isSimultaneous = currentPlayer === 0
  const isYourTurn = isSimultaneous || isHumanTurn

  // Calcular tiempo restante en el frontend
  const ticksPerSecond = 5 // 200ms por tick = 5 ticks/segundo
  const ticksSincePhaseStart = (state.tick || 0) - (state.phaseStartTick || 0)
  const elapsedSeconds = ticksSincePhaseStart / ticksPerSecond
  const phaseDurationSeconds = phaseInfo.duration / ticksPerSecond
  const remainingSeconds = Math.max(0, Math.ceil(phaseDurationSeconds - elapsedSeconds))
  
  // Calcular progreso de la fase
  const progress = phaseInfo.duration > 0 ? Math.min(100, (ticksSincePhaseStart / phaseInfo.duration) * 100) : 0

  // Formatear tiempo MM:SS
  const formatTime = (seconds) => {
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${mins}:${secs.toString().padStart(2, '0')}`
  }

  // Contar unidades
  const allUnits = state.units ? Object.values(state.units) : []
  const myUnits = allUnits.filter(u => u.playerId === playerId && u.hp > 0)
  const enemyUnits = allUnits.filter(u => u.playerId !== playerId && u.hp > 0)
  const myBaseHP = allUnits.find(u => u.playerId === playerId && u.unitType.toLowerCase() === 'base')?.hp || 0
  const enemyBaseHP = allUnits.find(u => u.playerId !== playerId && u.unitType.toLowerCase() === 'base')?.hp || 0

  return (
    <div className={`game-status-panel ${isYourTurn ? 'your-turn' : ''} ${isExpanded ? 'expanded' : 'collapsed'}`}>
      <div className="status-header">
        <div className="status-title">🎮 Game State</div>
        <button 
          className="status-toggle-btn"
          onClick={() => setIsExpanded(!isExpanded)}
          title={isExpanded ? 'Collapse' : 'Expand'}
        >
          {isExpanded ? '▼' : '▶'}
        </button>
      </div>

      {isExpanded && (
        <div className="status-content">
          <div className="status-row">
            <span className="status-label">🎮 Game:</span>
            <span className="status-value">{gameId || '-'}</span>
          </div>

          <div className="status-row">
            <span className="status-label">🎯 Turn:</span>
            <span className="status-value">{state.turnNumber || 0}</span>
          </div>

          {/* Phase con barra de progreso */}
          <div className="phase-section">
            <div className="status-row phase-row">
              <span className="status-label">Phase:</span>
              <span className={`phase-badge ${state.currentPhase}`} style={{ backgroundColor: phaseInfo.color }}>
                <span className="phase-emoji">{phaseInfo.emoji}</span> {phaseInfo.label}
              </span>
            </div>
            
            {/* Timer y barra de progreso */}
            {phaseInfo.duration > 0 && (
              <div className="phase-timer-section">
                <div className="phase-timer">
                  ⏱️ {formatTime(remainingSeconds)}
                </div>
                <div className="phase-progress-bar">
                  <div 
                    className="phase-progress-fill" 
                    style={{ 
                      width: `${Math.min(100, progress)}%`,
                      backgroundColor: phaseInfo.color
                    }}
                  />
                </div>
              </div>
            )}
          </div>

          {/* Turn indicator prominente */}
          <div className={`turn-indicator ${isSimultaneous ? 'simultaneous' : (isHumanTurn ? 'human' : 'ai')}`}>
            {isSimultaneous ? '⏱️ SIMULTANEOUS' : (isHumanTurn ? '👤 YOUR TURN' : '🤖 AI TURN')}
          </div>

          <div className="status-divider"></div>

          {/* Unidades */}
          <div className="status-row">
            <span className="status-label">👤 Your Units:</span>
            <span className="status-value unit-count">{myUnits.length}</span>
          </div>

          <div className="status-row">
            <span className="status-label">🤖 Enemy Units:</span>
            <span className="status-value unit-count">{enemyUnits.length}</span>
          </div>

          {/* Base HP */}
          <div className="status-row">
            <span className="status-label">🏰 Bases:</span>
            <span className="base-hp">
              <span className="my-base">♥ {myBaseHP}</span>
              <span className="separator">vs</span>
              <span className="enemy-base">♥ {enemyBaseHP}</span>
            </span>
          </div>
        </div>
      )}
    </div>
  )
}
