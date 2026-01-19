import './GameHandPanel.css'

export default function GameHandPanel({ state, playerId, selectedCard, onSelectCard }) {
  if (!state || playerId !== state?.humanPlayerId) {
    return null
  }

  // Solo mostrar en la fase de preparation
  const isPreparationPhase = state?.currentPhase === 'preparation'
  if (!isPreparationPhase) {
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

  return (
    <div className="game-hand-panel">
      <div className="hand-header">
        <div className="hand-title">🎴 Hand</div>
        <div className="hand-count">{myHand.length}</div>
      </div>

      <div className="hand-cards-container">
        {uniqueCards.length === 0 ? (
          <div className="no-cards">No cards</div>
        ) : (
          uniqueCards.map(cardType => (
            <div
              key={cardType}
              className={`hand-card ${selectedCard === cardType ? 'selected' : ''}`}
              onClick={() => onSelectCard(selectedCard === cardType ? null : cardType)}
              title={getCardName(cardType)}
            >
              <div className="card-emoji">{getCardEmoji(cardType)}</div>
              <div className="card-name">{getCardName(cardType)}</div>
              <div className="card-count">{cardCount[cardType]}</div>
            </div>
          ))
        )}
      </div>
    </div>
  )
}
