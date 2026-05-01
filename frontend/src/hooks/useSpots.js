import { useState, useCallback } from "react";
import { spotsAPI } from "../state/api";

export function useSpots() {
  const [spots, setSpots] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchSpots = useCallback(async (params = {}) => {
    setLoading(true);
    setError(null);
    try {
      const { data } = await spotsAPI.list(params);
      setSpots(data);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }, []);

  const fetchNearby = useCallback(async (lat, lon, radius = 25) => {
    setLoading(true);
    setError(null);
    try {
      const { data } = await spotsAPI.nearby({ lat, lon, radius });
      setSpots(data);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }, []);

  return { spots, loading, error, fetchSpots, fetchNearby };
}
