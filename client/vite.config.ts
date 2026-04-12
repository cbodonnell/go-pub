import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    proxy: {
      '/users': 'http://localhost:8080',
      '/activities': 'http://localhost:8080',
      '/objects': 'http://localhost:8080',
      '/uploads': 'http://localhost:8080',
      '/.well-known': 'http://localhost:8080',
      '/auth': 'http://localhost:8080',
      '/proxy': 'http://localhost:8080',
    },
  },
})
