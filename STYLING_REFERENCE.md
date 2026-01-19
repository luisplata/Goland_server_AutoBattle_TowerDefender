# 🎨 Frontend Styling Reference Guide

## Quick CSS Variables Reference

```css
/* Copy to FloatingPanels.css :root section */

/* Colors */
--panel-bg-dark: rgba(15, 15, 20, 0.92);      /* Dark panel background */
--panel-bg-light: rgba(30, 30, 40, 0.95);     /* Lighter panel background */
--panel-border: #00ff88;                       /* Neon green border */
--panel-accent: #4CAF50;                       /* Material green accent */
--panel-border-secondary: rgba(0, 255, 136, 0.3);

/* Text */
--panel-text-primary: #ffffff;                 /* White text */
--panel-text-secondary: #a0a0a0;               /* Gray text */

/* Shadows */
--panel-shadow: 0 4px 16px rgba(0, 0, 0, 0.4),
                0 0 12px rgba(0, 255, 136, 0.15);
--panel-shadow-hover: 0 8px 24px rgba(0, 0, 0, 0.5),
                      0 0 20px rgba(0, 255, 136, 0.25);
```

---

## CSS Class Reference

### Panel Container
```css
.floating-panel {
  position: relative;
  width: 100%;
  max-width: 340px;
  background: var(--panel-bg-dark);
  border: 2px solid var(--panel-border);
  border-radius: 12px;
  backdrop-filter: blur(10px);
  box-shadow: var(--panel-shadow);
  transition: all 0.3s ease;
}

/* Add to panels for hover effect */
.floating-panel:hover {
  box-shadow: var(--panel-shadow-hover);
  border-color: var(--panel-accent);
}
```

### Header Styling
```css
.panel-header {
  background: linear-gradient(135deg, rgba(0, 255, 136, 0.12), rgba(0, 255, 136, 0.05));
  border-bottom: 1px solid rgba(0, 255, 136, 0.4);
  padding: 0.875rem 1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.panel-title {
  font-weight: 700;
  color: var(--panel-border);
  font-size: 1rem;
  text-shadow: 0 0 8px rgba(0, 255, 136, 0.2);
}
```

### Toggle Button
```css
.panel-toggle-btn {
  width: 32px;
  height: 32px;
  background: rgba(0, 255, 136, 0.15);
  border: 1.5px solid var(--panel-border);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.panel-toggle-btn:hover {
  background: rgba(0, 255, 136, 0.25);
  transform: translateY(-1px);
}

.panel-toggle-btn:active {
  transform: scale(0.92);
}
```

### Content Area
```css
.panel-content {
  padding: 1rem;
  background: var(--panel-bg-light);
  max-height: 60vh;
  overflow-y: auto;
}

/* Custom scrollbar */
.panel-content::-webkit-scrollbar {
  width: 6px;
}

.panel-content::-webkit-scrollbar-track {
  background: rgba(0, 255, 136, 0.05);
  border-radius: 10px;
}

.panel-content::-webkit-scrollbar-thumb {
  background: rgba(0, 255, 136, 0.3);
  border-radius: 10px;
}

.panel-content::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 255, 136, 0.5);
}
```

### Buttons
```css
/* Primary Button (Green) */
.panel-button-primary {
  background: linear-gradient(135deg, var(--panel-accent), #66BB6A);
  color: #ffffff;
  padding: 0.75rem 1rem;
  border: none;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 4px 12px rgba(76, 175, 80, 0.3);
}

.panel-button-primary:hover {
  box-shadow: 0 6px 16px rgba(76, 175, 80, 0.5);
  transform: translateY(-2px);
}

.panel-button-primary:active {
  transform: translateY(0);
}

/* Secondary Button (Blue) */
.panel-button-secondary {
  background: rgba(100, 150, 200, 0.2);
  color: #90CAF9;
  border: 1px solid rgba(144, 202, 249, 0.5);
  padding: 0.75rem 1rem;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.panel-button-secondary:hover {
  background: rgba(100, 150, 200, 0.3);
  transform: translateY(-2px);
}

.panel-button-secondary:active {
  transform: translateY(0);
}
```

---

## Implementation Examples

### Creating a New Panel

```jsx
/* MyPanel.jsx */
import './MyPanel.css';

function MyPanel() {
  const [isExpanded, setIsExpanded] = useState(true);
  
  return (
    <div className="floating-panel">
      <div className="panel-header">
        <h3 className="panel-title">My Panel</h3>
        <button 
          className="panel-toggle-btn"
          onClick={() => setIsExpanded(!isExpanded)}
        >
          {isExpanded ? '▼' : '▶'}
        </button>
      </div>
      
      {isExpanded && (
        <div className="panel-content">
          <div className="panel-row">
            <span className="panel-label">Label</span>
            <span className="panel-value">Value</span>
          </div>
          <button className="panel-button panel-button-primary">
            Action
          </button>
        </div>
      )}
    </div>
  );
}
```

```css
/* MyPanel.css */
@import './FloatingPanels.css';

.my-panel {
  position: fixed;
  top: 1rem;
  right: 1rem;
  width: 100%;
  max-width: 340px;
  z-index: 100;
}
```

### Creating a Custom Color Theme

```css
/* Override in FloatingPanels.css :root */
:root {
  /* Blue theme example */
  --panel-border: #00a8ff;
  --panel-accent: #0084d1;
  --panel-shadow: 0 4px 16px rgba(0, 0, 0, 0.4), 
                  0 0 12px rgba(0, 168, 255, 0.15);
}

/* All panels automatically update */
```

---

## Spacing Standards

```css
/* Consistent spacing throughout UI */

/* Padding */
--space-xs: 0.25rem    /* Extra small */
--space-sm: 0.5rem     /* Small */
--space-md: 0.75rem    /* Medium (standard gap) */
--space-lg: 1rem       /* Large (content padding) */
--space-xl: 1.5rem     /* Extra large */

/* Border Radius */
--radius-sm: 4px       /* Button radius */
--radius-md: 6px       /* Standard radius */
--radius-lg: 12px      /* Panel radius */

/* Shadow Elevation */
--shadow-none: none
--shadow-sm: 0 2px 4px rgba(0, 0, 0, 0.1)
--shadow-md: var(--panel-shadow)
--shadow-lg: var(--panel-shadow-hover)

/* Transitions */
--transition-quick: 0.2s cubic-bezier(0.4, 0, 0.2, 1)
--transition-normal: 0.3s ease
--transition-slow: 0.5s ease-in-out
```

---

## Color Palette Breakdown

### Neon Green Theme
```
Primary:    #00ff88 (Neon bright green)
Secondary:  #4CAF50 (Material green)
Accent:     #00ff88 (Same as primary)
Glow:       rgba(0, 255, 136, 0.x) (Varying opacity)
```

### Neutral Tones
```
Dark BG:    rgba(15, 15, 20, 0.92)
Light BG:   rgba(30, 30, 40, 0.95)
Primary Text: #ffffff
Secondary Text: #a0a0a0
Muted:      #666666
```

### Type-Specific Colors
```
Spawn:      #4CAF50 (Green)
Damage:     #FF5722 (Orange)
Death:      #FF1744 (Red)
Phase:      #2196F3 (Blue)
Turn:       #FFD700 (Gold)
```

---

## Animation Reference

### Standard Transitions
```css
/* Hover effect */
transition: all 0.3s ease;

/* Button press effect */
transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);

/* Custom animations */
@keyframes slideUp {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}
```

### Transform Effects
```css
/* Hover scale */
transform: scale(1.05);

/* Hover lift */
transform: translateY(-1px);  /* Small lift */
transform: translateY(-2px);  /* Large lift */

/* Press effect */
transform: scale(0.92);       /* 8% scale down */
transform: scale(0.95);       /* 5% scale down */
```

---

## Responsive Breakpoints

```css
/* Mobile (default) */
@media (max-width: 480px) {
  /* Adjust for small phones */
  --panel-max-width: 280px;
}

/* Tablet */
@media (max-width: 768px) {
  /* Adjust for tablets */
  --panel-max-width: 320px;
}

/* Desktop */
@media (min-width: 769px) {
  /* Keep 340px max-width */
  --panel-max-width: 340px;
}
```

---

## Common Issues & Solutions

### Issue: Colors not updating
**Solution:** Make sure to use CSS variables
```css
/* ✅ Correct */
background: var(--panel-bg-dark);

/* ❌ Wrong */
background: rgba(15, 15, 20, 0.92);
```

### Issue: Hover effect not working
**Solution:** Ensure transition is defined
```css
/* ✅ Correct */
transition: all 0.3s ease;
transform: scale(1.05);

/* ❌ Wrong - no transition */
transform: scale(1.05);
```

### Issue: Scrollbar not styled
**Solution:** Use webkit prefixes for Chrome/Safari
```css
/* ✅ Correct */
::-webkit-scrollbar { width: 6px; }
::-webkit-scrollbar-track { background: ...; }
::-webkit-scrollbar-thumb { background: ...; }

/* ❌ Wrong - standard scrollbar property */
scrollbar-width: thin;
```

### Issue: Blur effect not working
**Solution:** Ensure backdrop-filter is present
```css
/* ✅ Correct */
backdrop-filter: blur(10px);

/* ❌ Wrong - missing backdrop-filter */
filter: blur(10px);
```

---

## Developer Workflow

### 1. **Update Global Variables**
Edit `FloatingPanels.css` `:root` section for theme changes

### 2. **Update Panel-Specific Styling**
Edit individual panel CSS files for custom overrides

### 3. **Apply Base Classes**
Use `.floating-panel`, `.panel-header`, `.panel-toggle-btn` etc.

### 4. **Test in Dev Server**
```bash
npm run dev
```

### 5. **Build for Production**
```bash
npm run build
```

---

## Performance Tips

✅ Use CSS variables instead of repeated values
✅ Minimize use of box-shadow (performance cost)
✅ Use transform for animations (GPU accelerated)
✅ Avoid filter effects on large elements
✅ Use will-change sparingly for animated elements

```css
/* Example: Optimized animation */
.panel-button {
  will-change: transform;
  transition: transform 0.2s ease;
}

.panel-button:hover {
  transform: translateY(-2px);  /* GPU accelerated */
}
```

---

## Support & Documentation

- **Full Documentation:** See `STYLING_IMPROVEMENTS.md`
- **Before/After Comparison:** See `STYLING_BEFORE_AFTER.md`
- **Completion Report:** See `STYLING_COMPLETE.md`
- **Source Code:** `frontend/src/components/FloatingPanels.css`

---

**Last Updated:** 2024
**Status:** Production Ready ✅
