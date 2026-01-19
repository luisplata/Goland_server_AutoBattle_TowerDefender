import { useEffect, useState } from 'react'
import './GameCardDetailsPanel.css'

export default function GameCardDetailsPanel({ state, playerId, selectedCard }) {
  // Solo mostrar en la fase de preparation
  const isPreparationPhase = state?.currentPhase === 'preparation'
  if (!state || playerId !== state?.humanPlayerId || !isPreparationPhase) {
    return null
  }

  const [unitStats, setUnitStats] = useState({})

  const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:7070'

  // Cargar estadísticas de unidades
  useEffect(() => {
    fetch(`${API_URL}/unit-stats`)
      .then(res => res.json())
      .then(data => setUnitStats(data))
      .catch(err => console.error('Error loading unit stats:', err))
  }, [API_URL])

  if (!state || !selectedCard) {
    return null
  }

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
    naval_generator: 'Naval Generator',
    land_generator: 'Land Generator',
  }

  const CARD_DESCRIPTIONS = {
    warrior: 'Mobile unit that can move and attack enemies.',
    tower: 'Defensive structure that attacks from range.',
    wall: 'Defensive structure that blocks movement.',
    naval_generator: 'Generates naval units over time.',
    land_generator: 'Generates land units over time.',
  }

  const getCardEmoji = (cardType) => CARD_EMOJIS[cardType] || '❓'
  const getCardName = (cardType) => CARD_NAMES[cardType] || cardType
  const getCardDescription = (cardType) => CARD_DESCRIPTIONS[cardType] || ''

  const stats = unitStats[selectedCard] || {}
  const myHand = state?.players?.[playerId]?.hand || []
  const cardCount = myHand.filter(c => c === selectedCard).length

  return (
    <div className="game-card-details-panel">
      <div className="details-header">
        <div className="details-title">📋 Card Details</div>
      </div>

      <div className="details-content">
        {/* Card header with emoji and name */}
        <div className="card-display">
          <div className="card-display-emoji">{getCardEmoji(selectedCard)}</div>
          <div className="card-display-name">{getCardName(selectedCard)}</div>
        </div>

        {/* Description */}
        <div className="card-description">
          {getCardDescription(selectedCard)}
        </div>

        {/* Stats grid */}
        {Object.keys(stats).length > 0 && (
          <div className="stats-grid">
            {stats.hp !== undefined && (
              <div className="stat-item">
                <span className="stat-icon">❤️</span>
                <span className="stat-name">HP</span>
                <span className="stat-val">{stats.hp}</span>
              </div>
            )}
            {stats.attackDamage !== undefined && (
              <div className="stat-item">
                <span className="stat-icon">⚔️</span>
                <span className="stat-name">DMG</span>
                <span className="stat-val">{stats.attackDamage}</span>
              </div>
            )}
            {stats.attackRange !== undefined && (
              <div className="stat-item">
                <span className="stat-icon">📏</span>
                <span className="stat-name">Range</span>
                <span className="stat-val">{stats.attackRange}</span>
              </div>
            )}
            {stats.buildRange !== undefined && (
              <div className="stat-item">
                <span className="stat-icon">🏗️</span>
                <span className="stat-name">Build</span>
                <span className="stat-val">{stats.buildRange}</span>
              </div>
            )}
            {stats.armor !== undefined && (
              <div className="stat-item">
                <span className="stat-icon">🛡️</span>
                <span className="stat-name">Armor</span>
                <span className="stat-val">{stats.armor}</span>
              </div>
            )}
            {stats.cost !== undefined && (
              <div className="stat-item">
                <span className="stat-icon">💰</span>
                <span className="stat-name">Cost</span>
                <span className="stat-val">{stats.cost}</span>
              </div>
            )}
          </div>
        )}

        {/* Available count */}
        <div className="card-count-display">
          <span className="count-label">Available in Hand:</span>
          <span className="count-value">{cardCount}</span>
        </div>
      </div>
    </div>
  )
}
