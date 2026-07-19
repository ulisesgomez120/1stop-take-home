import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    // The backend's CORS allowlist is fixed to http://localhost:5173. If this
    // port is taken, fail loudly instead of silently drifting to 5174 and
    // breaking every API call with opaque CORS errors.
    strictPort: true,
  },
})
