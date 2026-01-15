import './GameBoard.css'

const UNIT_EMOJIS = {
  warrior: '⚔️',
  tower: '🏰',
  wall: '🧱',
  naval_generator: '🌊',
}

export default function GameBoard({ state }) {
  if (!state) {
    return <div className="game-board">Waiting for game state...</div>
  }

  const renderGameInfo = () => {
    return (
      <div className="game-info">
        <div className="info-row">
          <span>📍 Tick: {state.tick || 0}</span>
          <span>🔄 Phase: {state.currentPhase || 'unknown'}</span>
          <span>🎯 Turn: {state.turnNumber || 0}</span>
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

  const renderUnits = () => {
    if (!state.units || Object.keys(state.units).length === 0) {
      return <div className="no-units">No units on the map</div>
    }

    return (
      <div className="units-section">
        <h3>⚡ Units</h3>
        <div className="units-grid">
          {Object.entries(state.units).map(([id, unit]) => (
            <div key={id} className={`unit-card player-${unit.playerId}`}>
              <div className="unit-header">
                <span className="unit-emoji">{UNIT_EMOJIS[unit.unitType] || '?'}</span>
                <span className="unit-name">{unit.unitType}</span>
              </div>
              <div className="unit-details">
                <span>ID: {unit.id}</span>
                <span>📍 ({unit.x}, {unit.y})</span>
                <span>❤️ HP: {unit.hp || 'N/A'}</span>
                <span>Owner: Player {unit.playerId}</span>
              </div>
            </div>
          ))}
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
