package domain

import (
	"fmt"
	"time"
)

// RiverReading representa uma leitura fluviométrica normalizada.
type RiverReading struct {
	StationID             string    `json:"station_id"`
	City                  string    `json:"city"`
	RiverName             string    `json:"river_name"`
	WaterLevelMeters      float64   `json:"water_level_meters"`
	VariationLastHourM    float64   `json:"variation_last_hour_meters"`
	Rainfall1hMM          float64   `json:"rainfall_1h_mm"`
	Rainfall24hMM         float64   `json:"rainfall_24h_mm"`
	ObservedAt            time.Time `json:"observed_at"`
	Source                string    `json:"source"`
}

func (r RiverReading) String() string {
	return fmt.Sprintf("[%s | %s] Nível: %.2fm (Δ1h: %+.2fm) | Chuva 1h: %.1fmm | Medido em: %s",
		r.City, r.StationID, r.WaterLevelMeters, r.VariationLastHourM, r.Rainfall1hMM, r.ObservedAt.Format("15:04 02/01/2006"))
}
