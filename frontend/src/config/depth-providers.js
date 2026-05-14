const GEBCO_ATTRIBUTION =
  "Imagery reproduced from the GEBCO_2024 Grid, GEBCO Compilation Group (2024)";

const GEBCO_TILE_URL =
  "https://wms.gebco.net/mapserv?bbox={bbox-epsg-3857}&service=WMS" +
  "&request=GetMap&version=1.1.1&layers=GEBCO_LATEST_2&styles=" +
  "&format=image/png&transparent=true&srs=EPSG:3857&width=256&height=256";

const providers = {
  gebco: {
    id: "gebco",
    name: "GEBCO",
    type: "raster",
    attribution: GEBCO_ATTRIBUTION,
    getSourceConfig() {
      return {
        type: "raster",
        tiles: [GEBCO_TILE_URL],
        tileSize: 256,
        attribution: GEBCO_ATTRIBUTION,
      };
    },
    getLayerIds() {
      return ["depth-raster"];
    },
    addLayers(style) {
      style.layers.push({
        id: "depth-raster",
        type: "raster",
        source: "depth",
        layout: { visibility: "none" },
      });
    },
    supportsPointQuery: true,
  },

  vectorcharts: {
    id: "vectorcharts",
    name: "VectorCharts",
    type: "vector",
    attribution: "\u00a9 VectorCharts.com",
    getSourceConfig(key) {
      return {
        type: "vector",
        tiles: [
          `https://api.vectorcharts.com/v1/tiles/{z}/{x}/{y}.pbf?key=${key}`,
        ],
        maxzoom: 14,
      };
    },
    getLayerIds() {
      return ["depth-areas", "depth-contours"];
    },
    addLayers(style) {
      style.layers.push(
        {
          id: "depth-areas",
          type: "fill",
          source: "depth",
          "source-layer": "deptharea",
          paint: {
            "fill-color": [
              "interpolate",
              ["linear"],
              ["get", "depth"],
              -8000, "#0a1628",
              -200, "#1a3a5c",
              -50, "#2a5a8c",
              -10, "#4a8aba",
              -2, "#7abada",
              0, "#c8e8f0",
            ],
            "fill-opacity": 0.4,
          },
          layout: { visibility: "none" },
        },
        {
          id: "depth-contours",
          type: "line",
          source: "depth",
          "source-layer": "depthcontour",
          paint: {
            "line-color": "#4a6a8a",
            "line-opacity": 0.6,
            "line-width": 1,
          },
          layout: { visibility: "none" },
        }
      );
    },
    supportsPointQuery: false,
  },
};

const PROVIDER_NAME = import.meta.env.VITE_DEPTH_PROVIDER || "gebco";
const VECTORCHARTS_KEY = import.meta.env.VITE_VECTORCHARTS_KEY;

function resolveProvider() {
  if (PROVIDER_NAME === "vectorcharts") {
    if (
      VECTORCHARTS_KEY && VECTORCHARTS_KEY !== "your_vectorcharts_api_key_here"
    ) {
      return { provider: providers.vectorcharts, key: VECTORCHARTS_KEY };
    }
    console.warn(
      "VectorCharts selected but no valid VITE_VECTORCHARTS_KEY — falling back to GEBCO"
    );
  }
  return { provider: providers.gebco, key: null };
}

export const { provider: activeProvider, key: providerKey } = resolveProvider();
