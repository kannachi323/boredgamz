import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
  ],
  server: {
    port: Number(process.env.VITE_PORT ?? 5201),
    strictPort: true,
    proxy: {
      // Same-origin in dev, exactly like production where Go serves the SPA.
      '/api': { target: process.env.VITE_API_PROXY ?? 'http://127.0.0.1:8201', ws: true, changeOrigin: false },
    },
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
})
