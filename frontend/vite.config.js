import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  // Serve the build under /playgame so asset URLs resolve when hosted by the Go server
  base: '/playgame/',
})
