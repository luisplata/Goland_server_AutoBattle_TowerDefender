# 📦 Referencia de Tipos de Unidades

## Todos los tipos disponibles

### 🏰 Estructuras
| Tipo | Emoji | Descripción | Color (Azul/Rojo) |
|------|-------|-------------|-------------------|
| `main_base` | 👑 | Base principal | Muy oscuro (15%) |
| `tower` | 🏰 | Torre de defensa | Oscuro (24%) |
| `wall` | 🧱 | Muralla bloqueadora | Medio-oscuro (34%) |
| `land_generator` | 🏞️ | Generador de unidades terrestres | Medio (44%) |
| `naval_generator` | 🌊 | Generador de unidades navales | Medio (50%) |

### ⚔️ Unidades
| Tipo | Emoji | Descripción | Color (Azul/Rojo) |
|------|-------|-------------|-------------------|
| `warrior` | ⚔️ | Guerrero (legacy) | Claro (64%) |
| `land_soldier` | 🗡️ | Soldado terrestre | Claro (64%) |
| `naval_ship` | ⛵ | Barco de guerra | Muy claro (85%) |

---

## Dónde están definidos

### Emojis
- **MapViewer.jsx**: No tiene (solo muestra círculos)
- **GameBoard.jsx**: líneas 12-22

```javascript
const UNIT_EMOJIS = {
  main_base: '👑',
  tower: '🏰',
  // ... etc
  default: '❓',
}
```

### Colores
- **MapViewer.jsx**: líneas 10-30
- **GameBoard.jsx**: líneas 25-45

```javascript
const UNIT_TYPE_INTENSITIES = {
  main_base: 0.0,     // 15% lightness
  tower: 0.14,        // 24% lightness
  // ... etc
}
```

---

## Cómo agregar un nuevo tipo

1. **En el backend** (`game/unit_types.go`):
   ```go
   const TypeNewUnit = "new_unit"
   ```

2. **En el frontend** (GameBoard.jsx):
   ```javascript
   // Agregar emoji
   const UNIT_EMOJIS = {
     // ...
     new_unit: '🆕',  // Tu emoji aquí
   }
   
   // Agregar intensidad de color
   const UNIT_TYPE_INTENSITIES = {
     // ...
     new_unit: 0.5,   // Entre 0 (muy oscuro) y 1 (muy claro)
   }
   ```

3. Listo! Los colores se aplicarán automáticamente en ambos componentes.

---

## Paleta de colores por intensidad

| Intensidad | Brillo | Ejemplo (Azul) |
|-----------|--------|----------------|
| 0.0 | 15% | 🔵 Azul muy oscuro |
| 0.14 | 24% | 🔵 Azul oscuro |
| 0.28 | 34% | 🔵 Azul medio-oscuro |
| 0.42 | 44% | 🔵 Azul medio |
| 0.5 | 50% | 🔵 Azul neutro |
| 0.64 | 64% | 🔵 Azul claro |
| 0.78 | 75% | 🔵 Azul más claro |
| 1.0 | 85% | 🔵 Azul muy claro |
