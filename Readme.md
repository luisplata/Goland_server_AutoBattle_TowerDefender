# Autobattle Server — Guía de Integración de Cliente

Este servidor expone endpoints HTTP para crear/unirse a partidas y un WebSocket para actualizaciones en tiempo real. El juego opera por fases: `base_selection → turn_start → preparation → battle → turn_end`.

## Flujo Básico
1) Crear juego

```bash
curl -X POST http://localhost:8080/game/create
```

Respuesta incluye `gameId` y `snapshot`.

2) Unirse (obtener `playerId`)

```bash
curl -X POST "http://localhost:8080/game/join?gameId=1"
```

3) Conectar WebSocket

```bash
wscat -c "ws://localhost:8080/ws?gameId=1&playerId=1"
```

- Enviar `playerId` en la URL permite tracking de conexión (timeouts por desconexión).
- El servidor envía cada tick un mensaje `snapshot` y eventos `phase_changed` y `hand_updated` cuando corresponda.

4) Consultar estado puntual (re-sync)

```bash
curl "http://localhost:8080/game/state?gameId=1"
```

## Comandos (HTTP)
Endpoint: `POST /command/send`

Payload base:

```json
{ "gameId": 1, "playerId": 1, "type": "spawn_unit", "data": { /* según tipo */ } }
```

Tipos soportados:
- `place_base` (solo en `base_selection`): `{ data: { x, y } }`
- `spawn_unit` (en `preparation`, requiere carta en mano): `{ data: { unitType, x, y } }`
- `move_unit` (en `preparation`, fija destino): `{ data: { unitId, x, y } }`
- `ready` (en `preparation`, marca listo): `{ data: null }`
- `confirm_end` (cuando `snapshot.gameEnd.pending` es true): `{ data: null }`
- `end_turn` (legacy, tratado como `ready`)

Ejemplos:

Colocar base (fase `base_selection`):

```bash
curl -X POST http://localhost:8080/command/send \
  -H "Content-Type: application/json" \
  -d '{"gameId":1,"playerId":1,"type":"place_base","data":{"x":50,"y":50}}'
```

Jugar carta (torre):

```bash
curl -X POST http://localhost:8080/command/send \
  -H "Content-Type: application/json" \
  -d '{"gameId":1,"playerId":1,"type":"spawn_unit","data":{"unitType":"tower","x":52,"y":50}}'
```

Mover unidad (fija destino; el movimiento ocurre por ticks):

```bash
curl -X POST http://localhost:8080/command/send \
  -H "Content-Type: application/json" \
  -d '{"gameId":1,"playerId":1,"type":"move_unit","data":{"unitId":7,"x":60,"y":50}}'
```

Listo para batalla:

```bash
curl -X POST http://localhost:8080/command/send \
  -H "Content-Type: application/json" \
  -d '{"gameId":1,"playerId":1,"type":"ready"}'
```

Confirmar fin (cuando `gameEnd.pending`):

```bash
curl -X POST http://localhost:8080/command/send \
  -H "Content-Type: application/json" \
  -d '{"gameId":1,"playerId":1,"type":"confirm_end"}'
```

## WebSocket — Mensajes
- `snapshot`: estado completo ({ tick, units, players, map, currentPhase, turnNumber, humanPlayerId, aiPlayerId, humanPlayerReady, aiPlayerReady, config, currentPlayerTurn, gameEnd? })
- `phase_changed`: { type, tick, previousPhase, currentPhase, turnNumber, humanPlayerId, aiPlayerId }
- `hand_updated`: { type, playerId, hand, deckCount }

Nota: actualmente el servidor emite `snapshot` cada tick (no “wrapper” de update/kind).

## Reglas Importantes
- `place_base` solo en `base_selection`.
- `spawn_unit`, `move_unit`, `ready` se permiten en `preparation`.
- La IA se marca lista automáticamente después de `config.aiReadyDelay`.
- Spawns deben estar en área controlada (rango `buildRange` de tus estructuras/base), con terreno válido y sin ocupar tiles.
- Navales solo en agua; terrestres/estructuras solo en tiles walkable.

## Movimiento y Combate
- `move_unit` fija un destino; el movimiento ocurre por ticks usando pathfinding.
- Las unidades con objetivo en rango de ataque se detienen para atacar.

## Desconexiones y Fin de Juego
- Si un cliente WS identificado por `playerId` se desconecta por más de `config.disconnectTimeoutSeconds`, el juego termina en su contra.
- Cuando se destruye una base, `snapshot.gameEnd.pending = true`. El humano debe enviar `confirm_end` para cerrar la partida.

## Estadísticas de Unidades
`GET /unit-stats` → mapa de `unitType -> UnitStats` para poblar UI (hp, dps, rango, etc.).

## Herramientas
- Swagger UI: http://localhost:8080/docs (sirve `openapi.yml`).
- wscat: `npm i -g wscat` y `wscat -c "ws://localhost:8080/ws?gameId=1&playerId=1"`.

## Docker

Construir:

```bash
docker build -t autobattle-server .
```

Ejecutar:

```bash
docker run -p 8080:8080 --name autobattle autobattle-server
```

Compose:

```bash
docker compose up -d
docker compose up -d --build  # después de cambios
docker compose logs -f
```

Limpiar y volver a ejecutar

```bash
docker-compose down
docker system prune -f
docker-compose up --build -d
```
