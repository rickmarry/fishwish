import { useState, useEffect } from "react";
import { socialAPI } from "../../state/api";

function ReviewSection({ spotId }) {
  const [reviews, setReviews] = useState([]);
  const [newReview, setNewReview] = useState({ rating: 5, content: "" });
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    loadReviews();
  }, [spotId]);

  const loadReviews = async () => {
    try {
      const { data } = await socialAPI.reviews(spotId);
      setReviews(data);
    } catch (e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      const { data } = await socialAPI.createReview(spotId, newReview);
      setReviews([data, ...reviews]);
      setNewReview({ rating: 5, content: "" });
    } catch (e) {
      console.error(e);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="card p-6">
      <h2 className="text-lg font-semibold text-gray-900 mb-4">Reviews</h2>

      <form onSubmit={handleSubmit} className="mb-6 p-4 bg-gray-50 rounded-lg">
        <div className="flex items-center gap-2 mb-2">
          <span className="text-sm text-gray-600">Rating:</span>
          {[1, 2, 3, 4, 5].map((n) => (
            <button
              key={n}
              type="button"
              onClick={() => setNewReview({ ...newReview, rating: n })}
              className={`text-xl ${n <= newReview.rating ? "text-yellow-400" : "text-gray-300"}`}
            >
              ★
            </button>
          ))}
        </div>
        <textarea
          value={newReview.content}
          onChange={(e) => setNewReview({ ...newReview, content: e.target.value })}
          placeholder="Share your experience at this spot..."
          className="input h-24 resize-none"
        />
        <button type="submit" disabled={submitting || !newReview.content} className="btn-primary mt-2 text-sm">
          {submitting ? "Submitting..." : "Post Review"}
        </button>
      </form>

      {loading && <p className="text-gray-400 text-sm">Loading reviews...</p>}

      {!loading && reviews.length === 0 && (
        <p className="text-gray-400 text-sm">No reviews yet. Be the first!</p>
      )}

      <div className="space-y-4">
        {reviews.map((review) => (
          <div key={review.id} className="border-t border-gray-100 pt-4">
            <div className="flex items-center justify-between mb-1">
              <span className="font-medium text-gray-900 text-sm">
                {review.username || "Anonymous"}
              </span>
              <span className="text-yellow-500 text-sm">
                {"★".repeat(review.rating)}
              </span>
            </div>
            <p className="text-gray-600 text-sm">{review.content}</p>
          </div>
        ))}
      </div>
    </div>
  );
}

export default ReviewSection;
