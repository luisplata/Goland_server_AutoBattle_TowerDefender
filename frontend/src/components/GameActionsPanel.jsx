import './GameActionsPanel.css'

export default function GameActionsPanel({ state, playerId, onCommand, selectedTile, selectedCard, onSelectCard, gameMap }) {
  if (!state || !playerId) {
    return null
  }

  const currentPhase = state?.currentPhase
  const isBaseSelectionPhase = currentPhase === 'base_selection'
  const isPreparationPhase = currentPhase === 'preparation'
  const hasPlacedBase = playerId === state?.humanPlayerId ? (state?.humanBaseId > 0) : (state?.aiBaseId > 0)

  // En base_selection solo muestra si es el jugador y no ha colocado base
  const shouldShowBaseSelection = isBaseSelectionPhase && !hasPlacedBase
  
  // En preparation muestra si está en esa fase
  const shouldShowPreparation = isPreparationPhase

  // Solo mostrar para el jugador humano
  if (playerId !== state?.humanPlayerId) {
    return null
  }

  // Si no hay ninguna acción disponible, no mostrar
  if (!shouldShowBaseSelection && !shouldShowPreparation) {
    return null
  }

  const myPlayer = state?.players?.[playerId]
  const myHand = myPlayer?.hand || []

  const CARD_EMOJIS = {
    warrior: '⚔️',
    tower: '🏰',
    wall: '🧱',
    naval_generator: '🌊',
    land_generator: '🏞️',
  }

  const getCardEmoji = (cardType) => CARD_EMOJIS[cardType] || '❓'

  const handlePlaceBase = () => {
    if (!selectedTile) {
      alert('Click on the map to select a position for your base')
      return
    }

    if (!selectedTile.walkable) {
      alert('Cannot place base on water')
      return
    }

    onCommand({
      type: 'place_base',
      data: { x: selectedTile.x, y: selectedTile.y }
    })
  }

  const handleSpawnUnit = () => {
    if (!selectedCard) {
      alert('Select a card first')
      return
    }

    if (!selectedTile) {
      alert('Click on the map to select a position')
      return
    }

    onCommand({
      type: 'spawn_unit',
      data: {
        unitType: selectedCard,
        x: selectedTile.x,
        y: selectedTile.y
      }
    })

    onSelectCard(null)
  }

  const handleEndTurn = () => {
    onCommand({
      type: 'end_turn'
    })
  }

  // Fase de selección de base
  if (shouldShowBaseSelection) {
    return (
      <div className="game-actions-panel">
        <div className="actions-header">
          <div className="actions-title">🏰 Place Base</div>
        </div>
        <div className="actions-content">
          <p className="action-hint">Click on the map to select a position</p>
          <button 
            className="action-button action-button-primary"
            onClick={handlePlaceBase}
            disabled={!selectedTile}
          >
            Place Base Here
          </button>
        </div>
      </div>
    )
  }

  // Fase de preparación - solo botones de acción
  if (shouldShowPreparation) {
    return (
      <div className="game-actions-panel">
        <div className="actions-header">
          <div className="actions-title">⚡ Actions</div>
        </div>
        <div className="actions-content">
          {/* Spawn button */}
          <button 
            className="action-button action-button-primary"
            onClick={handleSpawnUnit}
            disabled={!selectedCard || !selectedTile}
          >
            Spawn Unit
          </button>

          {/* End turn button */}
          <button 
            className="action-button action-button-secondary"
            onClick={handleEndTurn}
          >
            End Turn
          </button>
        </div>
      </div>
    )
  }

  return null
}
