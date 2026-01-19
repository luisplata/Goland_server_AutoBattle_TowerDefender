import { useState, useRef, useEffect } from 'react'
import './RangeIndicator.css'

export default function RangeIndicator({ gameMap, units, selectedUnitId, playerId, pan, zoom, tileSize = 3 }) {
  const [mousePos, setMousePos] = useState({ x: 0, y: 0 })
  const canvasRef = useRef(null)

  const selectedUnit = selectedUnitId && units ? units.find(u => u.id === selectedUnitId) : null

  useEffect(() => {
    if (!canvasRef.current || !selectedUnit) return

    const canvas = canvasRef.current
    const ctx = canvas.getContext('2d')
    
    // Limpiar canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height)

    if (!selectedUnit) return

    // Calcular posición en pantalla
    const cx = (selectedUnit.x * tileSize + pan.x) * zoom + canvas.width / 2
    const cy = (selectedUnit.y * tileSize + pan.y) * zoom + canvas.height / 2

    // Dibujar círculo de detección
    if (selectedUnit.detectionRange > 0) {
      ctx.strokeStyle = 'rgba(0, 255, 136, 0.3)'
      ctx.lineWidth = 2
      ctx.setLineDash([5, 5])
      ctx.beginPath()
      ctx.arc(cx, cy, selectedUnit.detectionRange * tileSize * zoom, 0, Math.PI * 2)
      ctx.stroke()
      ctx.setLineDash([])
    }

    // Dibujar círculo de ataque
    if (selectedUnit.attackRange > 0) {
      ctx.strokeStyle = 'rgba(255, 87, 34, 0.5)'
      ctx.lineWidth = 3
      ctx.beginPath()
      ctx.arc(cx, cy, selectedUnit.attackRange * tileSize * zoom, 0, Math.PI * 2)
      ctx.stroke()
    }

    // Dibujar círculo de construcción
    if (selectedUnit.buildRange > 0) {
      ctx.strokeStyle = 'rgba(33, 150, 243, 0.4)'
      ctx.lineWidth = 2
      ctx.setLineDash([3, 3])
      ctx.beginPath()
      ctx.arc(cx, cy, selectedUnit.buildRange * tileSize * zoom, 0, Math.PI * 2)
      ctx.stroke()
      ctx.setLineDash([])
    }
  }, [selectedUnit, pan, zoom, tileSize])

  const handleMouseMove = (e) => {
    const rect = canvasRef.current.getBoundingClientRect()
    setMousePos({
      x: e.clientX - rect.left,
      y: e.clientY - rect.top
    })
  }

  if (!selectedUnit) return null

  return (
    <canvas
      ref={canvasRef}
      className="range-indicator-canvas"
      onMouseMove={handleMouseMove}
      style={{
        position: 'fixed',
        top: 0,
        left: 0,
        pointerEvents: 'none'
      }}
    />
  )
}
