import { useState, useCallback } from "react";
import { weatherAPI } from "../state/api";

export function useWeather() {
  const [conditions, setConditions] = useState(null);
  const [fishingScore, setFishingScore] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchWeather = useCallback(async (lat, lon) => {
    setLoading(true);
    setError(null);
    try {
      const { data } = await weatherAPI.get({ lat, lon });
      setConditions(data.conditions);
      setFishingScore(data.fishing_score);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }, []);

  return { conditions, fishingScore, loading, error, fetchWeather };
}
