import { useEffect, useState } from 'react'
import './GameEventLog.css'

export default function GameEventLog({ state, playerId }) {
  const [events, setEvents] = useState([])
  const [lastTick, setLastTick] = useState(0)
  const [isExpanded, setIsExpanded] = useState(false)

  useEffect(() => {
    if (!state || !playerId) return

    // Detectar cambios en el juego para agregar eventos
    if (state.tick !== lastTick) {
      const newEvents = []
      const allUnits = state.units ? Object.values(state.units) : []

      // Detectar nuevas unidades (spawn)
      const currentUnitIds = allUnits.map(u => u.id)

      // Detectar cambios de fase
      if (lastTick > 0) {
        // Aquí agregaríamos lógica más compleja para detectar eventos
        // Por ahora, agregamos eventos de demo
      }

      setLastTick(state.tick)
    }
  }, [state, playerId, lastTick])

  // Agregar evento manualmente (puede ser llamado desde otros componentes)
  const addEvent = (type, message) => {
    const newEvent = {
      id: Date.now(),
      type,
      message,
      timestamp: new Date().getTime()
    }
    setEvents(prev => [newEvent, ...prev].slice(0, 10)) // Máximo 10 eventos
  }

  // Exponer método globalmente para que otros componentes puedan agregar eventos
  useEffect(() => {
    window.addGameEvent = addEvent
    return () => delete window.addGameEvent
  }, [])

  return (
    <div className={`game-event-log ${isExpanded ? 'expanded' : 'collapsed'}`}>
      {/* Botón toggle */}
      <button 
        className="event-log-toggle-button"
        onClick={() => setIsExpanded(!isExpanded)}
        title={isExpanded ? 'Hide events' : 'Show events'}
      >
        📋
      </button>

      {/* Panel expandido */}
      {isExpanded && (
        <div className="event-log-panel">
          <div className="event-log-header">📋 Events</div>
          <div className="event-log-content">
            {events.length === 0 ? (
              <div className="no-events">No events yet</div>
            ) : (
              events.map(event => (
                <div key={event.id} className={`event-item event-${event.type}`}>
                  <span className="event-time">
                    {new Date(event.timestamp).toLocaleTimeString().split(':').slice(1).join(':')}
                  </span>
                  <span className="event-message">{event.message}</span>
                </div>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  )
}
