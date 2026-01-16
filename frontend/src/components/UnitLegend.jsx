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
      <div className="legend-grid">
        {Object.entries(UNIT_DESCRIPTIONS).map(([key, info]) => (
          <div key={key} className="legend-item">
            <span className="legend-emoji">{info.emoji}</span>
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
