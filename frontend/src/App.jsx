import { useState, useEffect } from 'react'
import GameBoard from './components/GameBoard'
import GameControls from './components/GameControls'
import './App.css'

function App() {
  const [gameId, setGameId] = useState(null)
  const [playerId, setPlayerId] = useState(null)
  const [gameState, setGameState] = useState(null)
  const [ws, setWs] = useState(null)
  const [connected, setConnected] = useState(false)

  const API_URL = 'http://localhost:8080'
  const WS_URL = 'ws://localhost:8080/ws'

  // Crear juego
  const createGame = async () => {
    try {
      const res = await fetch(`${API_URL}/game/create`, { method: 'POST' })
      const data = await res.json()
      setGameId(data.gameId)
      setGameState(data.snapshot)
    } catch (err) {
      console.error('Error creating game:', err)
    }
  }

  // Unirse al juego
  const joinGame = async () => {
    if (!gameId) return
    try {
      const res = await fetch(`${API_URL}/game/join?gameId=${gameId}`)
      const data = await res.json()
      setPlayerId(data.id)
      connectWebSocket(data.id)
    } catch (err) {
      console.error('Error joining game:', err)
    }
  }

  // Conectar a WebSocket
  const connectWebSocket = (pid) => {
    const wsUrl = `${WS_URL}?gameId=${gameId}&playerId=${pid}`
    const newWs = new WebSocket(wsUrl)

    newWs.onopen = () => {
      setConnected(true)
      console.log('WebSocket connected to:', wsUrl)
    }

    newWs.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        
        // El backend envía directamente el objeto con type='snapshot' o similar
        if (message.type === 'snapshot' || message.type === 'update') {
          setGameState(message)
        } else if (message.payload) {
          // Por si vienen con payload
          setGameState(message.payload)
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
    try {
      await fetch(`${API_URL}/command/send`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          gameId,
          playerId,
          ...command
        })
      })
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
        {!gameId ? (
          <div className="setup">
            <h2>Start New Game</h2>
            <button onClick={createGame} className="btn-primary">Create Game</button>
          </div>
        ) : !playerId ? (
          <div className="setup">
            <h2>Join Game {gameId}</h2>
            <button onClick={joinGame} className="btn-primary">Join as Human Player</button>
          </div>
        ) : (
          <>
            <GameBoard state={gameState} />
            <GameControls 
              state={gameState}
              playerId={playerId}
              onCommand={sendCommand}
            />
          </>
        )}
      </div>
    </div>
  )
}

export default App
