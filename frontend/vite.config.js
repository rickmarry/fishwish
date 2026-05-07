import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3006,
    strictPort: true,
    host: true,
    proxy: {
      "/api/users": {
        target: "http://localhost:8086",
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
      "/api/spots": {
        target: "http://localhost:8082",
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
      "/api/search": "http://localhost:8083",
      "/api/weather": {
        target: "http://localhost:8084",
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
      "/api/social": {
        target: "http://localhost:8085",
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  },
});
