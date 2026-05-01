import { Outlet, Link, useLocation } from "react-router-dom";

const navLinks = [
  { path: "/", label: "Map" },
  { path: "/search", label: "Search" },
  { path: "/dashboard", label: "Dashboard" },
];

function Layout() {
  const location = useLocation();

  return (
    <div className="min-h-screen flex flex-col">
      <header className="bg-white border-b border-gray-200 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <Link to="/" className="flex items-center gap-2">
              <span className="text-2xl">&#128031;</span>
              <span className="text-xl font-bold text-ocean-700">FishFinder</span>
            </Link>

            <nav className="hidden md:flex items-center gap-1">
              {navLinks.map((link) => (
                <Link
                  key={link.path}
                  to={link.path}
                  className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                    location.pathname === link.path
                      ? "bg-ocean-50 text-ocean-700"
                      : "text-gray-600 hover:bg-gray-100"
                  }`}
                >
                  {link.label}
                </Link>
              ))}
            </nav>

            <Link to="/auth" className="btn-primary text-sm">
              Sign In
            </Link>
          </div>
        </div>
      </header>

      <main className="flex-1">
        <Outlet />
      </main>

      <footer className="bg-gray-900 text-gray-400 py-8">
        <div className="max-w-7xl mx-auto px-4 text-center text-sm">
          <p>FishFinder &mdash; Find your next great fishing spot.</p>
        </div>
      </footer>
    </div>
  );
}

export default Layout;
