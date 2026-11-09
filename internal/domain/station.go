package domain

// WatershedZone define a região hidrográfica da estação na bacia do Itajaí.
type WatershedZone string

const (
	ZoneUpperValley  WatershedZone = "ALTO_VALE"   // Rio do Sul, Taió, Ituporanga
	ZoneMiddleValley WatershedZone = "MEDIO_VALE"  // Ibirama, Apiúna, Indaial
	ZoneLowerValley  WatershedZone = "BAIXO_VALE"  // Blumenau, Gaspar, Ilhota, Itajaí
)

// FloodThresholds define as cotas críticas do rio em metros para cada estação.
type FloodThresholds struct {
	AttentionMeters float64 `json:"attention_meters"` // Cota de atenção
	AlertMeters     float64 `json:"alert_meters"`     // Cota de alerta
	EmergencyMeters float64 `json:"emergency_meters"` // Cota de emergência / transbordamento
}

// Station representa os metadados de uma estação de monitoramento.
type Station struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	City       string          `json:"city"`
	River      string          `json:"river"`
	Zone       WatershedZone   `json:"zone"`
	Thresholds FloodThresholds `json:"thresholds"`
	IsActive   bool            `json:"is_active"`
}

// StationsCatalog retorna o catálogo inicial de estações monitoradas na bacia.
func StationsCatalog() map[string]Station {
	return map[string]Station{
		"BLU-CENTRO": {
			ID:    "BLU-CENTRO",
			Name:  "Ponte Adolfo Konder / Centro",
			City:  "Blumenau",
			River: "Rio Itajaí-Açu",
			Zone:  ZoneLowerValley,
			Thresholds: FloodThresholds{
				AttentionMeters: 6.00,
				AlertMeters:     7.00,
				EmergencyMeters: 8.00,
			},
			IsActive: true,
		},
		"RDS-TITOBUSS": {
			ID:    "RDS-TITOBUSS",
			Name:  "Ponte Dom Tito Buss",
			City:  "Rio do Sul",
			River: "Rio Itajaí do Sul",
			Zone:  ZoneUpperValley,
			Thresholds: FloodThresholds{
				AttentionMeters: 5.50,
				AlertMeters:     6.50,
				EmergencyMeters: 7.50,
			},
			IsActive: true,
		},
	}
}
