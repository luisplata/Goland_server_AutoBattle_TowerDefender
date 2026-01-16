import './UnitLegend.css'

const UNIT_EMOJIS = {
  main_base: '👑',
  tower: '🏰',
  wall: '🧱',
  land_generator: '🏞️',
  naval_generator: '🌊',
  warrior: '⚔️',
  land_soldier: '🗡️',
  naval_ship: '⛵',
  default: '❓',
}

// Colores consistentes con el mapa
const COLOR_CONFIG = {
  PLAYER_1_HUE: 200,        // Azul (aliados)
  PLAYER_2_HUE: 0,          // Rojo (enemigos)
  SATURATION: 100,
  LIGHTNESS_MIN: 15,
  LIGHTNESS_MAX: 85,
}

const UNIT_TYPE_INTENSITIES = {
  main_base: 0.0,
  tower: 0.14,
  wall: 0.28,
  land_generator: 0.42,
  naval_generator: 0.5,
  warrior: 0.64,
  land_soldier: 0.78,
  naval_ship: 1.0,
}

const getUnitColor = (playerId, unitType) => {
  const hue = playerId === 1 ? COLOR_CONFIG.PLAYER_1_HUE : COLOR_CONFIG.PLAYER_2_HUE
  const saturation = COLOR_CONFIG.SATURATION
  const intensity = UNIT_TYPE_INTENSITIES[unitType] ?? 0.5
  const lightness = COLOR_CONFIG.LIGHTNESS_MIN + (intensity * (COLOR_CONFIG.LIGHTNESS_MAX - COLOR_CONFIG.LIGHTNESS_MIN))
  return `hsl(${hue}, ${saturation}%, ${lightness}%)`
}

const UNIT_DESCRIPTIONS = {
  main_base: { emoji: '👑', name: 'Main Base', desc: 'Your base. If destroyed, you lose.' },
  tower: { emoji: '🏰', name: 'Tower', desc: 'Defense structure. Range 25, Damage 25.' },
  wall: { emoji: '🧱', name: 'Wall', desc: 'Blocker. No target, just blocks pathways.' },
  land_generator: { emoji: '🏞️', name: 'Land Generator', desc: 'Spawns land soldiers.' },
  naval_generator: { emoji: '🌊', name: 'Naval Generator', desc: 'Spawns naval ships.' },
  warrior: { emoji: '⚔️', name: 'Warrior', desc: 'Basic land unit. Range 2, Damage 10.' },
  land_soldier: { emoji: '🗡️', name: 'Land Soldier', desc: 'Land unit. Range 2, Damage 15.' },
  naval_ship: { emoji: '⛵', name: 'Naval Ship', desc: 'Water unit. Range 15, Damage 20.' },
}

export default function UnitLegend() {
  return (
    <div className="unit-legend">
      <h3>📖 Unit Legend</h3>
      
      {/* Color indicators for teams */}
      <div style={{ 
        marginBottom: '1rem', 
        padding: '0.75rem', 
        background: 'rgba(0,0,0,0.3)', 
        borderRadius: '6px',
        border: '1px solid #444'
      }}>
        <div style={{ fontWeight: 'bold', marginBottom: '0.5rem', fontSize: '0.9rem' }}>🎨 Team Colors:</div>
        <div style={{ display: 'flex', gap: '1rem', justifyContent: 'space-around' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            <div style={{ 
              width: '20px', 
              height: '20px', 
              borderRadius: '50%', 
              backgroundColor: getUnitColor(1, 'warrior'),
              border: '2px solid white',
              boxShadow: '0 0 8px rgba(0,150,255,0.5)'
            }}></div>
            <span style={{ fontSize: '0.9rem' }}>🔵 Allies (Blue)</span>
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            <div style={{ 
              width: '20px', 
              height: '20px', 
              borderRadius: '50%', 
              backgroundColor: getUnitColor(2, 'warrior'),
              border: '2px solid white',
              boxShadow: '0 0 8px rgba(255,50,50,0.5)'
            }}></div>
            <span style={{ fontSize: '0.9rem' }}>🔴 Enemies (Red)</span>
          </div>
        </div>
      </div>

      <div className="legend-grid">
        {Object.entries(UNIT_DESCRIPTIONS).map(([key, info]) => (
          <div key={key} className="legend-item">
            <div style={{ display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
              <span className="legend-emoji">{info.emoji}</span>
              {/* Mostrar círculos de color para cada tipo */}
              <div style={{ display: 'flex', gap: '3px' }}>
                <div style={{ 
                  width: '12px', 
                  height: '12px', 
                  borderRadius: '50%', 
                  backgroundColor: getUnitColor(1, key),
                  border: '1px solid white'
                }} title="Ally color"></div>
                <div style={{ 
                  width: '12px', 
                  height: '12px', 
                  borderRadius: '50%', 
                  backgroundColor: getUnitColor(2, key),
                  border: '1px solid white'
                }} title="Enemy color"></div>
              </div>
            </div>
            <div className="legend-info">
              <div className="legend-name">{info.name}</div>
              <div className="legend-desc">{info.desc}</div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
