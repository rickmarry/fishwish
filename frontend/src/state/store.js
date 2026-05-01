import { create } from "zustand";
import { authAPI } from "../state/api";

export const useAuthStore = create((set) => ({
  user: null,
  token: localStorage.getItem("access_token"),
  isLoading: false,

  login: async (email, password) => {
    set({ isLoading: true });
    try {
      const { data } = await authAPI.login({ email, password });
      localStorage.setItem("access_token", data.access_token);
      set({ user: { email }, token: data.access_token, isLoading: false });
    } catch (e) {
      set({ isLoading: false });
      throw e;
    }
  },

  register: async (email, username, password) => {
    set({ isLoading: true });
    try {
      const { data } = await authAPI.register({ email, username, password });
      localStorage.setItem("access_token", data.access_token);
      set({ user: { email, username }, token: data.access_token, isLoading: false });
    } catch (e) {
      set({ isLoading: false });
      throw e;
    }
  },

  logout: () => {
    localStorage.removeItem("access_token");
    set({ user: null, token: null });
  },
}));

export const useSpotStore = create((set) => ({
  spots: [],
  selectedSpot: null,
  isLoading: false,

  fetchSpots: async (params) => {
    set({ isLoading: true });
    try {
      const { data } = await import("../state/api").then((m) => m.spotsAPI.list(params));
      set({ spots: data, isLoading: false });
    } catch {
      set({ isLoading: false });
    }
  },

  selectSpot: (spot) => set({ selectedSpot: spot }),
}));
