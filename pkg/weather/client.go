package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Conditions struct {
	TemperatureF    float64 `json:"temperature_f"`
	WindSpeedMph    float64 `json:"wind_speed_mph"`
	WindDirection   string  `json:"wind_direction"`
	Humidity        float64 `json:"humidity"`
	BarometricPress float64 `json:"barometric_pressure_in"`
	WaterTempF      float64 `json:"water_temp_f"`
	TideHeightFt    float64 `json:"tide_height_ft"`
	TideState       string  `json:"tide_state"`
	CloudCover      float64 `json:"cloud_cover_pct"`
	PrecipChance    float64 `json:"precip_chance_pct"`
	UVIndex         int     `json:"uv_index"`
	UpdatedAt       string  `json:"updated_at"`
}

type Client struct {
	apiKey string
	client *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetConditions(lat, lon float64) (*Conditions, error) {
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f"+
			"&current=temperature_2m,relative_humidity_2m,precipitation_probability,cloud_cover,wind_speed_10m,wind_direction_10m,surface_pressure",
		lat, lon,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status: %d", resp.StatusCode)
	}

	var data struct {
		Current struct {
			Temp            float64 `json:"temperature_2m"`
			Humidity        float64 `json:"relative_humidity_2m"`
			PrecipChance    float64 `json:"precipitation_probability"`
			CloudCover      float64 `json:"cloud_cover"`
			WindSpeed       float64 `json:"wind_speed_10m"`
			WindDirection   int     `json:"wind_direction_10m"`
			SurfacePressure float64 `json:"surface_pressure"`
		} `json:"current"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &Conditions{
		TemperatureF:    celsiusToFahrenheit(data.Current.Temp),
		WindSpeedMph:    kmhToMph(data.Current.WindSpeed),
		WindDirection:   degToCardinal(data.Current.WindDirection),
		Humidity:        data.Current.Humidity,
		BarometricPress: hpaToInHg(data.Current.SurfacePressure),
		WaterTempF:      0,
		TideHeightFt:    0,
		TideState:       "unknown",
		CloudCover:      data.Current.CloudCover,
		PrecipChance:    data.Current.PrecipChance,
		UVIndex:         0,
		UpdatedAt:       now.Format(time.RFC3339),
	}, nil
}

func celsiusToFahrenheit(c float64) float64 {
	return c*9.0/5.0 + 32.0
}

func kmhToMph(k float64) float64 {
	return k * 0.621371
}

func hpaToInHg(h float64) float64 {
	return h * 0.02953
}

func degToCardinal(deg int) string {
	directions := []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}
	return directions[(deg+11)%360/22]
}
