import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    host: true,
    proxy: {
      "/api/users": "http://user-service:8081",
      "/api/spots": "http://spot-service:8082",
      "/api/search": "http://search-service:8083",
      "/api/weather": "http://weather-service:8084",
      "/api/social": "http://social-service:8085",
    },
  },
});
