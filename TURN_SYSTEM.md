# AutoBattle Server - Sistema por Turnos

## 🎮 Sistema de Turnos

El juego ahora funciona en modo **1 vs IA por turnos**:
- Cuando un jugador se une al juego, automáticamente se crea un jugador IA
- El jugador humano siempre comienza primero
- Cada jugador puede ejecutar acciones solo durante su turno
- El turno se alterna entre el jugador humano y la IA

## 1. Crear juego
`curl -X POST http://localhost:8080/game/create`

Respuesta:
```json
{
  "gameId": 1,
  "snapshot": {
    "tick": 0,
    "players": {},
    "units": {},
    "currentPlayerTurn": 0,
    "aiPlayerId": 0,
    "humanPlayerId": 0,
    "turnBasedMode": true
  }
}
```

## 2. Unirse (obtener playerId)
`curl "http://localhost:8080/game/join?gameId=1"`

Respuesta:
```json
{
  "id": 1,
  "isAi": false
}
```

**Nota:** Al unirse el primer jugador (humano), se crea automáticamente el jugador IA con `id: 2`.

## 3. Enviar comando

Los comandos solo se procesan si es el turno del jugador que los envía.

### Spawnear unidad (warrior)
```bash
curl -X POST http://localhost:8080/command/send \
  -H "Content-Type: application/json" \
  -d '{"gameId":1,"playerId":1,"type":"spawn_unit","data":{"unitType":"warrior","x":50,"y":50}}'
```

### Mover unidad
```bash
curl -X POST http://localhost:8080/command/send \
  -H "Content-Type: application/json" \
  -d '{"gameId":1,"playerId":1,"type":"move_unit","data":{"unitId":1,"x":51,"y":50}}'
```

### 🎯 Terminar turno
```bash
curl -X POST http://localhost:8080/command/send \
  -H "Content-Type: application/json" \
  -d '{"gameId":1,"playerId":1,"type":"end_turn"}'
```

**Importante:** 
- Solo puedes realizar acciones durante tu turno
- Debes terminar tu turno manualmente con el comando `end_turn`
- Actualmente la IA termina su turno automáticamente (lógica pendiente)

## 4. Ver estado del juego
`curl "http://localhost:8080/game/state?gameId=1"`

Respuesta incluye información del turno:
```json
{
  "tick": 10,
  "currentPlayerTurn": 1,
  "humanPlayerId": 1,
  "aiPlayerId": 2,
  "turnBasedMode": true,
  "players": {
    "1": {"id": 1, "isAi": false},
    "2": {"id": 2, "isAi": true}
  },
  "units": { ... }
}
```

## 5. WebSocket para actualizaciones en tiempo real

```bash
# wscat (instalar: npm install -g wscat)
wscat -c "ws://localhost:8080/ws?gameId=1"
```

El servidor envía actualizaciones cada tick con información del turno actual.

## 🧪 Testing con test.html

Abre `test.html` en tu navegador y usa la consola (F12):

```javascript
// Ver comandos disponibles
help()

// Flujo típico de juego
await createGame()
await joinGame()
connectWS()

// Realizar acciones en tu turno
await spawnUnit('warrior', 5, 5)
await moveUnit(1, 10, 5)

// Terminar tu turno
await endTurn()
```

### Notas de movimiento por tiempo:
- Las unidades móviles avanzan 1 casilla cada cierto intervalo de ticks.
- Ejemplo inicial: `warrior` se mueve 1 casilla por ~1s (5 ticks si el tick es 200ms).
- `tower` no se mueve.
- El comando `move_unit` fija destino; el avance se realiza automáticamente en cada tick.

### Troubleshooting movimiento
- Ver logs del servidor: `docker compose logs -f`
- Confirmar que la ruta es transitable: el mapa marca agua como no-walkable; si cruza agua la unidad se detiene.
- Probar movimientos horizontales cortos (p.ej. `x: 55 -> 60` con `y` fijo) para evitar agua.
- Observar deltas con `wscat` (ver abajo) y buscar eventos `moved` cada ~1s.

## 🐳 Docker

Construir la imagen:
```bash
docker build -t autobattle-server .
```

Ejecutar el contenedor (expone 8080):
```bash
docker run -p 8080:8080 --name autobattle autobattle-server
```

## Docker Compose

Levantar el servicio:
```bash
docker compose up -d
```

Recrear después de cambios en el código:
```bash
docker compose up -d --build
```

Ver logs:
```bash
docker compose logs -f
```

## 📝 Tipos de Comandos

| Comando | Descripción | Requiere Turno |
|---------|-------------|----------------|
| `spawn_unit` | Crear una unidad | ✅ Sí |
| `move_unit` | Mover una unidad | ✅ Sí |
| `end_turn` | Terminar turno | ✅ Sí |

## 🤖 Estado de la IA

Actualmente la IA:
- Se crea automáticamente cuando el jugador humano se une
- Recibe su turno automáticamente
- Termina su turno inmediatamente (sin lógica implementada)

La implementación de la lógica de IA se agregará posteriormente.

## 🔧 Próximas Implementaciones

- [ ] Lógica de IA para tomar decisiones durante su turno
- [ ] Límite de tiempo por turno
- [ ] Sistema de recursos por turno
- [ ] Validación de acciones permitidas por turno
