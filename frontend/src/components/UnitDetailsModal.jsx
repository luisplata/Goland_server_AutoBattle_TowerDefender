import './UnitDetailsModal.css'

export default function UnitDetailsModal({ unit, playerId, onClose }) {
  if (!unit) return null

  const isAlly = unit.playerId === playerId
  const teamLabel = isAlly ? '🔵 Aliado' : '🔴 Enemigo'

  const getUnitTypeLabel = (unitType) => {
    const labels = {
      main_base: 'Base Principal',
      tower: 'Torre',
      wall: 'Muro',
      land_generator: 'Generador de Tierra',
      naval_generator: 'Generador Naval',
      warrior: 'Guerrero',
      land_soldier: 'Soldado de Tierra',
      naval_ship: 'Barco',
    }
    return labels[unitType] || unitType
  }

  return (
    <div className="unit-details-modal-overlay" onClick={onClose}>
      <div className="unit-details-modal" onClick={(e) => e.stopPropagation()}>
        
        <div className="unit-details-header">
          <div>
            <h2>{getUnitTypeLabel(unit.unitType)}</h2>
            <span className="unit-details-team">{teamLabel}</span>
          </div>
        </div>

        <div className="unit-details-content">
          <div className="unit-details-grid">
            <div className="unit-details-item">
              <span className="label">ID:</span>
              <span className="value">{unit.id}</span>
            </div>
            <div className="unit-details-item">
              <span className="label">Posición:</span>
              <span className="value">({unit.x}, {unit.y})</span>
            </div>
            <div className="unit-details-item">
              <span className="label">HP:</span>
              <span className={`value ${unit.hp <= unit.maxHp * 0.25 ? 'critical' : unit.hp <= unit.maxHp * 0.5 ? 'warning' : 'healthy'}`}>
                {unit.hp}/{unit.maxHp}
              </span>
            </div>
            {unit.attackDamage > 0 && (
              <>
                <div className="unit-details-item">
                  <span className="label">Daño:</span>
                  <span className="value">{unit.attackDamage}</span>
                </div>
                <div className="unit-details-item">
                  <span className="label">Rango:</span>
                  <span className="value">{unit.attackRange}</span>
                </div>
                <div className="unit-details-item">
                  <span className="label">DPS:</span>
                  <span className="value">{unit.attackDps?.toFixed(2) || 'N/A'}</span>
                </div>
              </>
            )}
            {unit.detectionRange > 0 && (
              <div className="unit-details-item">
                <span className="label">Detección:</span>
                <span className="value">{unit.detectionRange}</span>
              </div>
            )}
            {unit.isGenerator && (
              <>
                <div className="unit-details-item">
                  <span className="label">Genera:</span>
                  <span className="value">{getUnitTypeLabel(unit.generatedUnitType)}</span>
                </div>
                <div className="unit-details-item">
                  <span className="label">Generadas:</span>
                  <span className="value">{unit.unitsGenerated}/{unit.maxUnitsGenerated}</span>
                </div>
              </>
            )}
            {unit.buildRange > 0 && (
              <div className="unit-details-item">
                <span className="label">Construcción:</span>
                <span className="value">{unit.buildRange}</span>
              </div>
            )}
          </div>

          <div className="hp-bar">
            <div 
              className={`hp-fill ${unit.hp <= unit.maxHp * 0.25 ? 'critical' : unit.hp <= unit.maxHp * 0.5 ? 'warning' : 'healthy'}`}
              style={{ width: `${(unit.hp / unit.maxHp) * 100}%` }}
            ></div>
          </div>

          {(unit.targetId > 0 || unit.status || unit.isBlocker || unit.isTargetable) && (
            <div className="unit-details-flags">
              {unit.targetId > 0 && <span className="unit-flag">🎯 Obj: #{unit.targetId}</span>}
              {unit.status && unit.status !== 'idle' && <span className="unit-flag">📍 {unit.status}</span>}
              {unit.isBlocker && <span className="unit-flag">🚫 Bloqueante</span>}
              {unit.isTargetable && <span className="unit-flag">🎯 Blanco</span>}
            </div>
          )}
        </div>

        <button className="unit-details-action-btn" onClick={onClose}>Cerrar</button>
      </div>
    </div>
  )
}
