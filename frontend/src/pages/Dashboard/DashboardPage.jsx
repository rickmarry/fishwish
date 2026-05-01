import { Link } from "react-router-dom";

function DashboardPage() {
  return (
    <div className="max-w-4xl mx-auto p-6">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Dashboard</h1>

      <div className="grid md:grid-cols-3 gap-6">
        <div className="card p-6">
          <h3 className="text-sm font-medium text-gray-500 mb-1">Saved Spots</h3>
          <p className="text-3xl font-bold text-ocean-600">0</p>
          <Link to="/search" className="text-ocean-600 text-sm mt-2 inline-block">
            Find spots &rarr;
          </Link>
        </div>

        <div className="card p-6">
          <h3 className="text-sm font-medium text-gray-500 mb-1">Catch Logs</h3>
          <p className="text-3xl font-bold text-forest-600">0</p>
          <button className="text-forest-600 text-sm mt-2 inline-block">
            Log a catch &rarr;
          </button>
        </div>

        <div className="card p-6">
          <h3 className="text-sm font-medium text-gray-500 mb-1">Reviews</h3>
          <p className="text-3xl font-bold text-yellow-600">0</p>
        </div>
      </div>

      <div className="mt-8 card p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Recent Activity</h2>
        <p className="text-gray-400 text-sm">
          Sign in to see your activity. <Link to="/auth" className="text-ocean-600">Sign in</Link>
        </p>
      </div>
    </div>
  );
}

export default DashboardPage;
