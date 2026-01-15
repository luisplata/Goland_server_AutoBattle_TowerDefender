# 🎨 Configuración de Colores de Unidades

## Ubicaciones donde editar colores

### 1. **MapViewer.jsx** - Colores en el mapa
📍 Archivo: `frontend/src/components/MapViewer.jsx` (líneas 10-30)

```javascript
const COLOR_CONFIG = {
  PLAYER_1_HUE: 200,        // Azul (0-360°)
  PLAYER_2_HUE: 0,          // Rojo (0-360°)
  SATURATION: 100,          // Saturación (0-100)
  LIGHTNESS_MIN: 15,        // Brillo mínimo - muy oscuro (0-100)
  LIGHTNESS_MAX: 85,        // Brillo máximo - muy claro (0-100)
}

const UNIT_TYPE_INTENSITIES = {
  main_base: 0.0,           // Muy oscuro (15%)
  tower: 0.14,              // Oscuro (24%)
  wall: 0.28,               // Medio-oscuro (34%)
  land_generator: 0.42,     // Medio (44%)
  naval_generator: 0.5,     // Medio (50%)
  warrior: 0.64,            // Claro (64%)
  land_soldier: 0.78,       // Más claro (75%)
  naval_ship: 1.0,          // Muy claro (85%)
}
```

### 2. **GameBoard.jsx** - Colores en las tarjetas de unidades
📍 Archivo: `frontend/src/components/GameBoard.jsx` (líneas 12-32)

**Misma configuración que MapViewer.jsx** para consistencia.

---

## Cómo personalizar

### Cambiar colores por bando
```javascript
PLAYER_1_HUE: 200   // Cambiar a 120 para verde, 300 para púrpura, etc.
PLAYER_2_HUE: 0     // Cambiar a otro valor
```

### Cambiar rango de brillo (diferencia entre tipos)
```javascript
LIGHTNESS_MIN: 15   // Sube para menos contraste, baja para más oscuro
LIGHTNESS_MAX: 85   // Baja para menos contraste, sube para más claro
```

Diferencia actual: 70% de contraste (15% a 85%)
- Para más diferencia: `MIN: 5, MAX: 95` (90% de rango)
- Para menos diferencia: `MIN: 35, MAX: 65` (30% de rango)

### Cambiar orden de tipos (oscuro a claro)
Edita los valores en `UNIT_TYPE_INTENSITIES` (0 = muy oscuro, 1 = muy claro):
```javascript
main_base: 0.0      // Cambiar a otro valor entre 0 y 1
warrior: 0.64       // etc...
```

---

## Valores HSL explicados

- **Hue (0-360°)**: Color
  - 0° = Rojo
  - 120° = Verde
  - 200° = Azul
  - 300° = Púrpura

- **Saturation (0-100%)**: Intensidad del color
  - 0% = Gris (sin color)
  - 100% = Color puro

- **Lightness (0-100%)**: Brillo
  - 0% = Negro
  - 50% = Color normal
  - 100% = Blanco

---

## Ejemplo de modificación

Para hacer las bases más oscuras que los towers:
```javascript
UNIT_TYPE_INTENSITIES: {
  main_base: 0.0,   // 15% (muy oscuro)
  tower: 0.14,      // 24% (oscuro)
}
```

Para cambiar el color tuyo a verde:
```javascript
COLOR_CONFIG: {
  PLAYER_1_HUE: 120  // Verde en lugar de azul (200)
}
```

---

## Rango actual de colores

| Tipo | Intensidad | Brillo | Azul Jugador 1 | Rojo Enemigo |
|------|-----------|--------|----------------|------------|
| main_base | 0.0 | 15% | Azul muy oscuro | Rojo muy oscuro |
| tower | 0.14 | 24% | Azul oscuro | Rojo oscuro |
| wall | 0.28 | 34% | Azul medio-oscuro | Rojo medio-oscuro |
| land_generator | 0.42 | 44% | Azul medio | Rojo medio |
| naval_generator | 0.5 | 50% | Azul neutro | Rojo neutro |
| warrior | 0.64 | 64% | Azul claro | Rojo claro |
| land_soldier | 0.78 | 75% | Azul más claro | Rojo más claro |
| naval_ship | 1.0 | 85% | Azul muy claro | Rojo muy claro |
