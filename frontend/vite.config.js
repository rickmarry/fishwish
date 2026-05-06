import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    host: true,
    proxy: {
      "/api/users": "http://localhost:8081",
      "/api/spots": "http://localhost:8082",
      "/api/search": "http://localhost:8083",
      "/api/weather": "http://localhost:8084",
      "/api/social": "http://localhost:8085",
    },
  },
});
