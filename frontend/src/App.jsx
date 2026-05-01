import { Routes, Route } from "react-router-dom";
import Layout from "./components/Layout/Layout";
import MapPage from "./pages/Map/MapPage";
import SpotDetailPage from "./pages/SpotDetail/SpotDetailPage";
import SearchPage from "./pages/Search/SearchPage";
import DashboardPage from "./pages/Dashboard/DashboardPage";
import AuthPage from "./pages/Auth/AuthPage";

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<MapPage />} />
        <Route path="search" element={<SearchPage />} />
        <Route path="spots/:id" element={<SpotDetailPage />} />
        <Route path="dashboard" element={<DashboardPage />} />
        <Route path="auth" element={<AuthPage />} />
      </Route>
    </Routes>
  );
}

export default App;
