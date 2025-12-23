# 1. Crear juego
`curl -X POST http://localhost:8080/game/create`

# 2. Unirse (obtener playerId)
`curl "http://localhost:8080/game/join?gameId=1"`

# 3. Enviar comando
```
curl -X POST http://localhost:8080/command/send \
  -H "Content-Type: application/json" \
  -d '{"gameId":1,"playerId":1,"type":"dummy"}'
```

## Spawnear unidad (warrior)
```
curl -X POST http://localhost:8080/command/send \
  -H "Content-Type: application/json" \
  -d '{"gameId":1,"playerId":1,"type":"spawn_unit","data":{"unitType":"warrior","x":50,"y":50}}'
```

## Mover unidad
Enviar movimiento (example):
```
curl -X POST http://localhost:8080/command/send \
  -H "Content-Type: application/json" \
  -d '{"gameId":1,"playerId":1,"type":"move_unit","data":{"unitId":1,"x":51,"y":50}}'
```

Notas de movimiento por tiempo:
- Las unidades móviles avanzan 1 casilla cada cierto intervalo de ticks.
- Ejemplo inicial: `warrior` se mueve 1 casilla por ~1s (5 ticks si el tick es 200ms).
- `tower` no se mueve.
- El comando `move_unit` fija destino; el avance se realiza automáticamente en cada tick.

### Troubleshooting movimiento
- Ver logs del servidor: `docker compose logs -f`
- Confirmar que la ruta es transitable: el mapa marca agua como no-walkable; si cruza agua la unidad se detiene.
- Probar movimientos horizontales cortos (p.ej. `x: 55 -> 60` con `y` fijo) para evitar agua.
- Observar deltas con `wscat` (ver abajo) y buscar eventos `moved` cada ~1s.

# 4. Ver estado
`curl "http://localhost:8080/game/state?gameId=1"`

# wscat (instalar: npm install -g wscat)
`wscat -c "ws://localhost:8080/ws?gameId=1"`

## Docker

Construir la imagen:
```
docker build -t autobattle-server .
```

Ejecutar el contenedor (expone 8080):
```
docker run -p 8080:8080 --name autobattle autobattle-server
```

## Docker Compose

Levantar el servicio:
```
docker compose up -d
```

Recrear después de cambios en el código:
```
docker compose up -d --build
```
