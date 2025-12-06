import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': 'http://localhost:8080',
      '/resource': 'http://localhost:8080',
      '/download': 'http://localhost:8080',
    }
  },
  build: {
    outDir: '../dist',
    emptyOutDir: true,
  }
})
