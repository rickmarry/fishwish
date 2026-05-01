import axios from "axios";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || "/api",
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("access_token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("access_token");
      window.location.href = "/auth";
    }
    return Promise.reject(error);
  }
);

export const spotsAPI = {
  list: (params) => api.get("/spots", { params }),
  get: (id) => api.get(`/spots/${id}`),
  details: (id) => api.get(`/spots/${id}/details`),
  nearby: (params) => api.get("/spots/nearby", { params }),
};

export const searchAPI = {
  search: (params) => api.get("/search", { params }),
  species: () => api.get("/search/species"),
  suggestions: (q) => api.get("/search/suggestions", { params: { q } }),
};

export const weatherAPI = {
  get: (params) => api.get("/weather", { params }),
  forecast: (params) => api.get("/weather/forecast", { params }),
  tides: (params) => api.get("/tides", { params }),
};

export const socialAPI = {
  reviews: (spotId) => api.get(`/spots/${spotId}/reviews`),
  createReview: (spotId, data) => api.post(`/spots/${spotId}/reviews`, data),
  logCatch: (data) => api.post("/catches", data),
  userCatches: (userId) => api.get(`/users/${userId}/catches`),
};

export const authAPI = {
  register: (data) => api.post("/auth/register", data),
  login: (data) => api.post("/auth/login", data),
  refresh: (data) => api.post("/auth/refresh", data),
};

export default api;
