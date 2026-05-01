import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { searchAPI } from "../../state/api";
import SpotCard from "../../components/SpotCard/SpotCard";

function SearchPage() {
  const [query, setQuery] = useState("");
  const [filters, setFilters] = useState({ type: "", species: "" });
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [suggestions, setSuggestions] = useState([]);
  const [showSuggestions, setShowSuggestions] = useState(false);

  useEffect(() => {
    if (query.length >= 2) {
      const timer = setTimeout(async () => {
        try {
          const { data } = await searchAPI.suggestions(query);
          setSuggestions(data);
          setShowSuggestions(true);
        } catch {
          setSuggestions([]);
        }
      }, 300);
      return () => clearTimeout(timer);
    }
    setSuggestions([]);
    setShowSuggestions(false);
  }, [query]);

  const handleSearch = async (e) => {
    e.preventDefault();
    if (!query.trim()) return;

    setLoading(true);
    setShowSuggestions(false);
    try {
      const { data } = await searchAPI.search({
        q: query,
        type: filters.type,
        species: filters.species,
      });
      setResults(data);
    } catch (e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto p-6">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Search Fishing Spots</h1>

      <form onSubmit={handleSearch} className="mb-6 relative">
        <div className="flex gap-2">
          <div className="flex-1 relative">
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Search spots, species, locations..."
              className="input pr-10"
            />
            {showSuggestions && suggestions.length > 0 && (
              <div className="absolute top-full left-0 right-0 mt-1 bg-white border border-gray-200 rounded-lg shadow-lg z-10">
                {suggestions.map((s, i) => (
                  <button
                    key={i}
                    type="button"
                    className="w-full text-left px-4 py-2 hover:bg-gray-50 text-sm"
                    onClick={() => {
                      setQuery(s);
                      setShowSuggestions(false);
                    }}
                  >
                    {s}
                  </button>
                ))}
              </div>
            )}
          </div>
          <select
            value={filters.type}
            onChange={(e) => setFilters({ ...filters, type: e.target.value })}
            className="input w-32"
          >
            <option value="">All Types</option>
            <option value="lake">Lake</option>
            <option value="river">River</option>
            <option value="ocean">Ocean</option>
            <option value="pond">Pond</option>
            <option value="reservoir">Reservoir</option>
            <option value="estuary">Estuary</option>
          </select>
          <button type="submit" className="btn-primary">
            Search
          </button>
        </div>
      </form>

      {loading && (
        <div className="text-center py-12 text-gray-400">Searching...</div>
      )}

      {!loading && results.length === 0 && query && (
        <div className="text-center py-12 text-gray-400">
          No results found. Try a different search.
        </div>
      )}

      <div className="grid sm:grid-cols-2 gap-4">
        {results.map((spot) => (
          <Link key={spot.id} to={`/spots/${spot.id}`}>
            <SpotCard spot={spot} />
          </Link>
        ))}
      </div>
    </div>
  );
}

export default SearchPage;
