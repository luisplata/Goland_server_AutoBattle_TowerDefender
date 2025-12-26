# Autobattle Server API

## HTTP Endpoints

### POST /game/create
- Body (optional):
```json
{
  "config": {
    "turnStartDuration": 15,
    "preparationDuration": 150,
    "battleDuration": 25,
    "turnEndDuration": 10,
    "aiReadyDelay": 5
  }
}
```
- Response 200:
```json
{
  "gameId": 1,
  "snapshot": { /* see Snapshot schema below */ }
}
```

### POST /game/join?gameId={id}
- Response 200:
```json
{
  "id": 1,
  "isAi": false
}
```

### GET /game/state?gameId={id}
- Response 200:
```json
{
  "type": "snapshot",
  "tick": 42,
  "currentPhase": "preparation",
  "turnNumber": 2,
  "humanPlayerId": 1,
  "aiPlayerId": 2,
  "humanPlayerReady": false,
  "aiPlayerReady": false,
  "config": { /* PhaseConfig */ },
  "players": {
    "1": { "id": 1, "isAi": false, "hand": ["tower", "wall"], "deckCount": 4, "connected": true },
    "2": { "id": 2, "isAi": true,  "hand": ["naval_generator"], "deckCount": 5, "connected": true }
  },
  "units": {
    "1": { "id": 1, "playerId": 1, "unitType": "tower", "x": 5, "y": 5, "hp": 500 }
  }
}
```

### POST /command/send
- Body (spawn example):
```json
{
  "gameId": 1,
  "playerId": 1,
  "type": "spawn_unit",
  "data": { "unitType": "tower", "x": 5, "y": 5 }
}
```
- Body (move example):
```json
{
  "gameId": 1,
  "playerId": 1,
  "type": "move_unit",
  "data": { "unitId": 7, "x": 8, "y": 9 }
}
```
- Body (ready example):
```json
{
  "gameId": 1,
  "playerId": 1,
  "type": "ready",
  "data": null
}
```
- Response: 202 Accepted (cola), 404 si gameId no existe, 400 si payload es invalido.

## WebSocket /ws?gameId={id}&playerId={playerId}
Recibe mensajes JSON:

### Update (snapshot completo)
```json
{
  "type": "update",
  "kind": "snapshot",
  "payload": {
    "type": "snapshot",
    "tick": 0,
    "currentPhase": "turn_start",
    "turnNumber": 1,
    "humanPlayerId": 1,
    "aiPlayerId": 2,
    "humanPlayerReady": false,
    "aiPlayerReady": false,
    "config": { /* PhaseConfig */ },
    "units": {},
    "players": {}
  }
}
```

### Update (delta incremental)
```json
{
  "type": "update",
  "kind": "delta",
  "payload": {
    "type": "delta",
    "tick": 21,
    "changedUnits": {
      "3": { "id": 3, "playerId": 1, "unitType": "land_soldier", "x": 6, "y": 5, "hp": 100 }
    },
    "removedUnits": [],
    "currentPhase": "battle",
    "turnNumber": 1,
    "humanPlayerReady": true,
    "aiPlayerReady": true
  }
}
```

### phase_changed
```json
{
  "type": "phase_changed",
  "previousPhase": "preparation",
  "currentPhase": "battle",
  "turnNumber": 2
}
```

### hand_updated
```json
{
  "type": "hand_updated",
  "playerId": 1,
  "hand": ["tower", "wall"],
  "deckCount": 4
}
```
Se emite al robar (inicio de preparation) o consumir carta (spawn).

## Esquemas
- PhaseConfig: ints (ticks) `turnStartDuration`, `preparationDuration`, `battleDuration`, `turnEndDuration`, `aiReadyDelay`.
- Player: `id`, `isAi`, `hand` (array de strings), `deckCount`, `connected` (bool).
- Unit: `id`, `playerId`, `unitType`, `x`, `y`, `hp`.
- Snapshot: campos de fase + `units` (map) + `players` (map) + `config`.
- Delta: `changedUnits` (map), `removedUnits` (array), fase y turn info.
- Cartas validas (unitType): `tower`, `land_generator`, `naval_generator`, `wall`, `warrior` (legacy). Generadas: `land_soldier`, `naval_ship` (no jugables por carta).

## Notas de validacion
- Comandos solo en `preparation` (excepto `ready`).
- Spawn requiere carta en mano y posicion valida: dentro de mapa, sin otra unidad, terreno apto (naval solo agua; resto walkable).
- Generadores spawnean en tiles adyacentes libres; si no hay espacio, loguean aviso.
 - Si un `playerId` se desconecta del WebSocket y permanece desconectado por 10s, la partida termina con derrota para ese jugador (se registra en logs).
