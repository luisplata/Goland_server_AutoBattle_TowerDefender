# Optimizaciones de Rendimiento del Servidor

## Problema Detectado
El servidor consumía **97.66% de CPU** debido a un bucle de polling excesivo.

## Optimizaciones Aplicadas

### 1. **Bucle Principal del Servidor** ([main.go](main.go))
- ✅ **Antes**: `time.Sleep(1 * time.Millisecond)` = ~1000 iteraciones/segundo
- ✅ **Después**: `time.Sleep(10 * time.Millisecond)` = ~100 iteraciones/segundo
- **Impacto**: Reducción del 90% en iteraciones innecesarias
- **Resultado**: ~100 updates/seg es más que suficiente para juegos en tiempo real

### 2. **Sleep Adaptativo Cuando No Hay Juegos**
- ✅ Si `len(games) == 0`, duerme 100ms en lugar de 10ms
- **Impacto**: CPU casi 0% cuando no hay juegos activos

### 3. **Optimización de UpdateTargets** ([game_simulation.go](game/game_simulation.go))
- ✅ **Antes**: Se ejecutaba cada tick
- ✅ **Después**: Se ejecuta cada 5 ticks (`tick % 5 == 0`)
- **Impacto**: Reducción del 80% en cálculos de targeting
- **Justificación**: Los enemigos no cambian de posición tan rápido como para necesitar actualización cada tick

### 4. **Límite de Pathfinding** ([pathfinding.go](game/pathfinding.go))
- ✅ **Antes**: `maxPathSearchSteps = MapWidth * MapHeight` (potencialmente 900+ pasos)
- ✅ **Después**: `maxPathSearchSteps = 200`
- **Impacto**: Cálculos de A* más rápidos, evita búsquedas infinitas
- **Trade-off**: Caminos muy largos pueden fallar, pero es aceptable en un juego de torre defense

### 5. **Cache de Pathfinding con Límite de Memoria**
- ✅ Agregado `maxSize: 1000` entradas en PathCache
- ✅ Auto-limpieza cuando se excede el límite
- **Impacto**: Evita crecimiento infinito de memoria

## Resultados Esperados

| Métrica | Antes | Después (Estimado) |
|---------|-------|-------------------|
| CPU en idle | ~20-30% | ~1-5% |
| CPU con 1 juego | ~97% | ~15-30% |
| Updates/segundo | ~1000 | ~100 |
| Targeting updates/seg | ~1000 | ~200 |

## Monitoreo Post-Optimización

Después de reiniciar el servidor, verifica con:
```bash
docker stats
```

Deberías ver una reducción significativa en el uso de CPU.

## Optimizaciones Futuras (Si Aún Hay Problemas)

1. **Event-Driven Architecture**: Usar canales de Go en lugar de polling
2. **Goroutine Pool**: Limitar goroutines concurrentes
3. **Profiling**: Usar `pprof` para identificar hotspots específicos
4. **Lazy Evaluation**: Calcular pathfinding solo cuando la unidad realmente se mueve
5. **Spatial Partitioning**: Usar grids para búsquedas de enemigos más rápidas

## Comandos Útiles

```bash
# Ver uso de CPU en tiempo real
docker stats

# Reiniciar el contenedor para aplicar cambios
docker-compose restart

# Ver logs del servidor
docker-compose logs -f autobattle

# Profiling de Go (si es necesario)
go tool pprof http://localhost:8080/debug/pprof/profile
```
