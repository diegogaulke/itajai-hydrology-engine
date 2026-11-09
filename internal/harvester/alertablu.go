package harvester

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/seu-usuario/itajai-hydrology-engine/internal/domain"
)

type AlertaBluHarvester struct {
	client  *http.Client
	jsonURL string
}

func NewAlertaBluHarvester(client *http.Client) *AlertaBluHarvester {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	return &AlertaBluHarvester{
		client:  client,
		jsonURL: "https://defesacivil.blumenau.sc.gov.br/static/data/nivel_oficial.json",
	}
}

func (a *AlertaBluHarvester) Name() string {
	return "AlertaBLU-Blumenau"
}

type alertaBluNivelDTO struct {
	Nivel       float64 `json:"nivel"`
	HoraLeitura string  `json:"horaLeitura"`
}

type alertaBluResponseDTO struct {
	Fonte  string              `json:"fonte"`
	Niveis []alertaBluNivelDTO `json:"niveis"`
}

func (a *AlertaBluHarvester) FetchReadings(ctx context.Context) ([]domain.RiverReading, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.jsonURL, nil)
	if err != nil {
		return nil, fmt.Errorf("alertablu: falha ao criar request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")
	req.Header.Set("Referer", "https://defesacivil.blumenau.sc.gov.br/d/nivel-do-rio")
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("alertablu: falha de conexao: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("alertablu: status HTTP inesperado: %d", resp.StatusCode)
	}

	var data alertaBluResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("alertablu: falha ao decodificar JSON: %w", err)
	}

	total := len(data.Niveis)
	if total == 0 {
		return nil, fmt.Errorf("alertablu: array de niveis veio vazio")
	}

	// O último elemento do array é a leitura mais recente
	latest := data.Niveis[total-1]

	var variationLastHour float64
	if total >= 2 {
		previous := data.Niveis[total-2]
		variationLastHour = latest.Nivel - previous.Nivel
	}

	observedAt, err := time.Parse(time.RFC3339, latest.HoraLeitura)
	if err != nil {
		observedAt = time.Now().UTC()
	}

	return []domain.RiverReading{
		{
			StationID:          "BLU-CENTRO",
			City:               "Blumenau",
			RiverName:          "Rio Itajaí-Açu",
			WaterLevelMeters:   latest.Nivel,
			VariationLastHourM: variationLastHour,
			ObservedAt:         observedAt.UTC(),
			Source:             a.Name(),
		},
	}, nil
}
