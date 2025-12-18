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

# 4. Ver estado
`curl "http://localhost:8080/game/state?gameId=1"`

# wscat (instalar: npm install -g wscat)
`wscat -c "ws://localhost:8080/ws?gameId=1"`