import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  optimizeDeps: {
    include: ["highlight.js"],
  },
  server: {
    port: 5173,
    proxy: {
      // Proxy API requests to Go backend (includes SSE /api/events)
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
      // Proxy static assets to Go backend
      "/static": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});
