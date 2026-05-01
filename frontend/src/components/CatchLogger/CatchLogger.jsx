import { useState } from "react";
import { socialAPI } from "../../state/api";

function CatchLogger({ spotId, onSuccess }) {
  const [form, setForm] = useState({
    species: "",
    weight_lbs: "",
    length_in: "",
    bait_used: "",
  });
  const [submitting, setSubmitting] = useState(false);
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      await socialAPI.logCatch({ spot_id: spotId, ...form });
      setSuccess(true);
      setForm({ species: "", weight_lbs: "", length_in: "", bait_used: "" });
      onSuccess?.();
      setTimeout(() => setSuccess(false), 3000);
    } catch (e) {
      console.error(e);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="card p-6">
      <h2 className="text-lg font-semibold text-gray-900 mb-4">Log a Catch</h2>

      {success && (
        <div className="bg-green-50 text-green-700 text-sm p-3 rounded-lg mb-4">
          Catch logged successfully!
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Species
          </label>
          <input
            type="text"
            value={form.species}
            onChange={(e) => setForm({ ...form, species: e.target.value })}
            className="input"
            placeholder="Largemouth Bass"
            required
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Weight (lbs)
            </label>
            <input
              type="number"
              step="0.1"
              value={form.weight_lbs}
              onChange={(e) => setForm({ ...form, weight_lbs: e.target.value })}
              className="input"
              placeholder="4.5"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Length (in)
            </label>
            <input
              type="number"
              step="0.1"
              value={form.length_in}
              onChange={(e) => setForm({ ...form, length_in: e.target.value })}
              className="input"
              placeholder="18.5"
            />
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Bait/Lure Used
          </label>
          <input
            type="text"
            value={form.bait_used}
            onChange={(e) => setForm({ ...form, bait_used: e.target.value })}
            className="input"
            placeholder="Plastic worm, green pumpkin"
          />
        </div>

        <button type="submit" disabled={submitting || !form.species} className="btn-primary w-full">
          {submitting ? "Logging..." : "Log Catch"}
        </button>
      </form>
    </div>
  );
}

export default CatchLogger;
