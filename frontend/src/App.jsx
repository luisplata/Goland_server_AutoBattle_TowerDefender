import { useState, useEffect } from 'react'
import GameBoard from './components/GameBoard'
import GameControls from './components/GameControls'
import CanvasMapViewer from './components/CanvasMapViewer'
import UnitLegend from './components/UnitLegend'
import './App.css'

function App() {
  const [gameId, setGameId] = useState(null)
  const [gameIdInput, setGameIdInput] = useState('')
  const [playerId, setPlayerId] = useState(null)
  const [gameState, setGameState] = useState(null)
  const [ws, setWs] = useState(null)
  const [connected, setConnected] = useState(false)
  const [selectedTile, setSelectedTile] = useState(null)
  const [lastTurn, setLastTurn] = useState(null)
  const [selectedUnitId, setSelectedUnitId] = useState(null)
  const [selectedCard, setSelectedCard] = useState(null)
  const [gameOver, setGameOver] = useState(null) // { loserId, winnerId, reason }

  const API_URL = 'http://localhost:7070'
  const WS_URL = 'ws://localhost:7070/ws'

  // Limpiar selección cuando cambia el turno
  useEffect(() => {
    if (gameState?.currentPlayerTurn && lastTurn !== null && gameState.currentPlayerTurn !== lastTurn) {
      setSelectedTile(null)
    }
    setLastTurn(gameState?.currentPlayerTurn)
  }, [gameState?.currentPlayerTurn])

  // Detectar fin de juego desde el estado del servidor (pendiente de confirmación)
  useEffect(() => {
    if (!gameState) return
    const ge = gameState.gameEnd
    if (ge && ge.pending && !ge.confirmed) {
      setGameOver({ loserId: ge.loserId, winnerId: ge.winnerId, reason: ge.reason })
    } else if (ge && ge.confirmed) {
      setGameOver({ loserId: ge.loserId, winnerId: ge.winnerId, reason: ge.reason })
    } else {
      setGameOver(null)
    }
  }, [gameState])

  // Crear juego
  const createGame = async () => {
    try {
      const res = await fetch(`${API_URL}/game/create`, { method: 'POST' })
      const data = await res.json()
      setGameId(data.gameId)
      setGameState(data.snapshot)
      setGameIdInput(String(data.gameId))
      return data.gameId
    } catch (err) {
      console.error('Error creating game:', err)
      return null
    }
  }

  // Unirse al juego
  const joinGame = async () => {
    // Si hay un gameId en el input, intentar unirse a ese
    let targetGameId = gameIdInput?.trim() !== '' ? Number(gameIdInput) : null
    
    if (targetGameId) {
      console.log('Attempting to join existing game:', targetGameId)
      
      // Check if we have a saved playerId for this game
      const savedPlayerId = localStorage.getItem(`playerId_${targetGameId}`)
      
      if (savedPlayerId) {
        // Reconnect with existing player ID
        console.log('Reconnecting with saved player ID:', savedPlayerId)
        setGameId(targetGameId)
        setPlayerId(Number(savedPlayerId))
        await fetchGameState(targetGameId)
        connectWebSocket(Number(savedPlayerId), targetGameId)
        return
      }
      
      try {
        const res = await fetch(`${API_URL}/game/join?gameId=${targetGameId}`)
        
        if (res.ok) {
          const data = await res.json()
          console.log('Successfully joined game:', targetGameId, 'Player ID:', data.id)
          setGameId(targetGameId)
          setPlayerId(data.id)
          localStorage.setItem(`playerId_${targetGameId}`, String(data.id))
          await fetchGameState(targetGameId)
          connectWebSocket(data.id, targetGameId)
          return
        } else if (res.status === 404) {
          console.log(`Game ${targetGameId} not found, will create a new one instead`)
        } else {
          console.error('Join failed:', res.status)
          return
        }
      } catch (err) {
        console.error('Error joining game:', err)
        return
      }
    }
    
    // Si no hay gameId o no existe, crear uno nuevo
    console.log('Creating new game...')
    try {
      const res = await fetch(`${API_URL}/game/create`, { method: 'POST' })
      if (!res.ok) {
        console.error('Failed to create game:', res.status)
        return
      }
      
      const data = await res.json()
      const newGameId = data.gameId
      console.log('Game created with ID:', newGameId)
      
      setGameId(newGameId)
      setGameState(data.snapshot)
      setGameIdInput(String(newGameId))
      
      // Ahora unirse al juego recién creado
      const joinRes = await fetch(`${API_URL}/game/join?gameId=${newGameId}`)
      if (!joinRes.ok) {
        console.error('Failed to join newly created game:', joinRes.status)
        return
      }
      
      const playerData = await joinRes.json()
      console.log('Successfully joined new game:', newGameId, 'Player ID:', playerData.id)
      setPlayerId(playerData.id)
      localStorage.setItem(`playerId_${newGameId}`, String(playerData.id))
      await fetchGameState(newGameId)
      connectWebSocket(playerData.id, newGameId)
    } catch (err) {
      console.error('Error creating/joining game:', err)
    }
  }

  // Obtener estado del juego
  const fetchGameState = async (gid) => {
    try {
      const res = await fetch(`${API_URL}/game/state?gameId=${gid}`)
      const data = await res.json()
      setGameState(data)
    } catch (err) {
      console.error('Error fetching game state:', err)
    }
  }

  // Conectar a WebSocket
    const connectWebSocket = (pid, gid) => {
      const wsUrl = `${WS_URL}?gameId=${gid}&playerId=${pid}`
    const newWs = new WebSocket(wsUrl)

    newWs.onopen = () => {
      setConnected(true)
      console.log('WebSocket connected to:', wsUrl)
    }

    newWs.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        
        // Procesar snapshots completos
        if (message.type === 'snapshot' || message.type === 'update') {
          setGameState(message)
          // Reset selection if map changes dimensions
          setSelectedTile((sel) => {
            const map = message?.map
            if (!map || !sel) return sel
            if (sel.x >= map.width || sel.y >= map.height) return null
            return sel
          })
        } 
        // Procesar deltas (cambios incrementales)
        else if (message.type === 'delta') {
          setGameState((prevState) => {
            if (!prevState) return prevState
            
            const newState = { ...prevState }
            
            // Aplicar unidades nuevas (spawned)
            if (message.spawned && message.spawned.length > 0) {
              message.spawned.forEach(unit => {
                newState.units = { ...newState.units, [unit.id]: unit }
              })
            }
            
            // Aplicar movimientos
            if (message.moved && message.moved.length > 0) {
              message.moved.forEach(move => {
                if (newState.units[move.id]) {
                  newState.units[move.id] = {
                    ...newState.units[move.id],
                    x: move.x,
                    y: move.y,
                  }
                }
              })
            }
            
            // Aplicar cambios de estado (TargetID, HP, Status)
            if (message.updated && message.updated.length > 0) {
              message.updated.forEach(update => {
                if (newState.units[update.id]) {
                  const unitUpdate = { ...newState.units[update.id] }
                  if (update.targetId !== undefined) unitUpdate.targetId = update.targetId
                  if (update.hp !== undefined) unitUpdate.hp = update.hp
                  if (update.status !== undefined) unitUpdate.status = update.status
                  newState.units[update.id] = unitUpdate
                }
              })
            }
            
            // Aplicar unidades muertas
            if (message.dead && message.dead.length > 0) {
              const newUnits = { ...newState.units }
              message.dead.forEach(unitId => {
                delete newUnits[unitId]
              })
              newState.units = newUnits
            }
            
            // Actualizar otros campos del estado
            if (message.tick !== undefined) newState.tick = message.tick
            if (message.currentPhase !== undefined) newState.currentPhase = message.currentPhase
            if (message.turnNumber !== undefined) newState.turnNumber = message.turnNumber
            if (message.humanPlayerReady !== undefined) newState.humanPlayerReady = message.humanPlayerReady
            if (message.aiPlayerReady !== undefined) newState.aiPlayerReady = message.aiPlayerReady
            
            return newState
          })
        } 
        else if (message.payload) {
          // Por si vienen con payload
          setGameState(message.payload)
          const map = message.payload?.map
          setSelectedTile((sel) => {
            if (!map || !sel) return sel
            if (sel.x >= map.width || sel.y >= map.height) return null
            return sel
          })
        }
      } catch (err) {
        console.error('Error parsing WebSocket message:', err)
      }
    }

    newWs.onerror = (err) => {
      console.error('WebSocket error:', err)
      setConnected(false)
    }

    newWs.onclose = () => {
      console.log('WebSocket disconnected')
      setConnected(false)
    }

    setWs(newWs)
  }

  // Enviar comando
  const sendCommand = async (command) => {
    if (!gameId || !playerId) return
    // Permitir confirmación de fin de juego aunque esté activo el overlay
    if (gameOver && command?.type !== 'confirm_end') return
    try {
      const res = await fetch(`${API_URL}/command/send`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          gameId,
          playerId,
          ...command
        })
      })

      // Si se confirmó fin de juego y el servidor aceptó, limpiar UI y conexión
      if (res.ok && command?.type === 'confirm_end') {
        try { ws?.close() } catch {}
        setWs(null)
        setGameOver(null)
        // Limpiar identificación persistida del jugador para este juego
        try { localStorage.removeItem(`playerId_${gameId}`) } catch {}
        // Resetear estado de sesión
        setPlayerId(null)
        setGameId(null)
        setGameIdInput('')
        setSelectedTile(null)
        setSelectedUnitId(null)
        setSelectedCard(null)
        setGameState(null)
        setConnected(false)
      }
    } catch (err) {
      console.error('Error sending command:', err)
    }
  }

  return (
    <div className="app">
      <header>
        <h1>🎮 AutoBattle Client</h1>
        <div className="status">
          <span>Game ID: {gameId || 'None'}</span>
          <span>Player ID: {playerId || 'None'}</span>
          <span className={`connection ${connected ? 'connected' : 'disconnected'}`}>
            {connected ? '🟢 Connected' : '🔴 Disconnected'}
          </span>
        </div>
      </header>

      <div className="container">
        {!playerId ? (
          <div className="setup">
            <h2>Start or Join a Game</h2>
            <div style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
              <input
                type="number"
                placeholder="Enter Game ID (optional)"
                value={gameIdInput}
                onChange={(e) => setGameIdInput(e.target.value)}
                style={{ padding: '0.6rem', borderRadius: 6, border: '1px solid #00ff88', background: 'rgba(255,255,255,0.1)', color: '#fff' }}
              />
              <button onClick={joinGame} className="btn-primary">Join Game</button>
              <button onClick={createGame} className="btn-primary">Create Game</button>
            </div>
            {gameId && (
              <div style={{ marginTop: '1rem', opacity: 0.8 }}>
                Current Game ID: {gameId}
              </div>
            )}
          </div>
        ) : (
          <>
            {gameOver && (
              <div className="overlay">
                <div className="overlay-content">
                  <h2>{gameOver.loserId === playerId ? 'Derrota' : 'Victoria'}</h2>
                  <p>
                    {gameOver.reason === 'human_base_destroyed' && gameOver.loserId === playerId ? 'Tu base fue destruida.' : ''}
                    {gameOver.reason === 'ai_base_destroyed' && gameOver.winnerId === playerId ? 'Has destruido la base enemiga.' : ''}
                  </p>
                  <div style={{ display: 'flex', gap: '1rem', marginTop: '1rem' }}>
                    <button className="btn-primary" onClick={() => sendCommand({ type: 'confirm_end' })}>Aceptar resultado</button>
                  </div>
                </div>
              </div>
            )}
            <UnitLegend />
            <CanvasMapViewer 
              gameMap={gameState?.map} 
              units={gameState?.units} 
              selectedTile={selectedTile}
              onSelectTile={(tile) => setSelectedTile(tile)}
              disableZoom={false}
              playerId={playerId}
              selectedCard={selectedCard}
            />
            <GameBoard 
              state={gameState}
              playerId={playerId}
              selectedUnitId={selectedUnitId}
              onSelectUnit={setSelectedUnitId}
            />
            <GameControls 
              state={gameState}
              playerId={playerId}
              onCommand={sendCommand}
              selectedTile={selectedTile}
              gameMap={gameState?.map}
              onClearSelection={() => setSelectedTile(null)}
              selectedUnitId={selectedUnitId}
              selectedCard={selectedCard}
              onSelectCard={setSelectedCard}
            />
          </>
        )}
      </div>
    </div>
  )
}

export default App
