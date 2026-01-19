import { useState } from 'react'
import './CollapsibleGameHandPanel.css'

export default function CollapsibleGameHandPanel({ state, playerId, selectedCard, onSelectCard }) {
  const [isExpanded, setIsExpanded] = useState(false)

  if (!state || playerId !== state?.humanPlayerId) {
    return null
  }

  const currentPhase = state?.currentPhase
  const isPreparationPhase = currentPhase === 'preparation'
  const myPlayer = state?.players?.[playerId]
  const myHand = myPlayer?.hand || []

  const CARD_EMOJIS = {
    warrior: '⚔️',
    tower: '🏰',
    wall: '🧱',
    naval_generator: '🌊',
    land_generator: '🏞️',
  }

  const CARD_NAMES = {
    warrior: 'Warrior',
    tower: 'Tower',
    wall: 'Wall',
    naval_generator: 'Naval Gen',
    land_generator: 'Land Gen',
  }

  const getCardEmoji = (cardType) => CARD_EMOJIS[cardType] || '❓'
  const getCardName = (cardType) => CARD_NAMES[cardType] || cardType

  // Contar cartas por tipo
  const cardCount = {}
  myHand.forEach(card => {
    cardCount[card] = (cardCount[card] || 0) + 1
  })

  const uniqueCards = Object.keys(cardCount)

  // En preparation, mostrar normalmente
  if (isPreparationPhase) {
    return (
      <div className="collapsible-hand-panel">
        <div className="collapsed-hand-content">
          <div className="collapsed-hand-header">
            <span className="collapsed-hand-title">🎴 Hand ({myHand.length})</span>
          </div>

          <div className="collapsed-hand-cards">
            {uniqueCards.length === 0 ? (
              <div className="no-cards">No cards</div>
            ) : (
              uniqueCards.map(cardType => (
                <div
                  key={cardType}
                  className={`collapsed-hand-card ${selectedCard === cardType ? 'selected' : ''}`}
                  onClick={() => onSelectCard(selectedCard === cardType ? null : cardType)}
                >
                  <div className="collapsed-card-emoji">{getCardEmoji(cardType)}</div>
                  <div className="collapsed-card-count">{cardCount[cardType]}</div>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    )
  }

  // En otros estados, mostrar panel colapsable
  return (
    <div className={`collapsible-hand-panel ${isExpanded ? 'expanded' : 'collapsed'}`}>
      {/* Botón toggle */}
      <button 
        className="hand-toggle-button"
        onClick={() => setIsExpanded(!isExpanded)}
        title={isExpanded ? 'Hide hand' : 'Show hand'}
      >
        🎴
      </button>

      {/* Panel expandido */}
      {isExpanded && (
        <div className="collapsed-hand-content">
          <div className="collapsed-hand-header">
            <span className="collapsed-hand-title">Hand ({myHand.length})</span>
          </div>

          <div className="collapsed-hand-cards">
            {uniqueCards.length === 0 ? (
              <div className="no-cards">No cards</div>
            ) : (
              uniqueCards.map(cardType => (
                <div
                  key={cardType}
                  className={`collapsed-hand-card ${selectedCard === cardType ? 'selected' : ''}`}
                  onClick={() => onSelectCard(selectedCard === cardType ? null : cardType)}
                >
                  <div className="collapsed-card-emoji">{getCardEmoji(cardType)}</div>
                  <div className="collapsed-card-count">{cardCount[cardType]}</div>
                </div>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  )
}
