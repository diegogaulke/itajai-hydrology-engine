package harvester

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/seu-usuario/itajai-hydrology-engine/internal/domain"
)

type RioDoSulHarvester struct {
	client  *http.Client
	apiURL  string
}

func NewRioDoSulHarvester(client *http.Client) *RioDoSulHarvester {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	return &RioDoSulHarvester{
		client: client,
		apiURL: "https://public.asthon.com.br/public/stations/live?city_id=4214805&_v=2",
	}
}

func (r *RioDoSulHarvester) Name() string {
	return "DefesaCivil-RioDoSul"
}

type asthonStationDTO struct {
	StationID     string   `json:"station_id"`
	LastReadingAt *string  `json:"last_reading_at"`
	LevelM        *float64 `json:"level_m"`
	RainfallMM    *float64 `json:"rainfall_mm"`
	OverflowCotaM *float64 `json:"overflow_cota_m"`
	BandKey       string   `json:"band_key"`
	BandLabel     string   `json:"band_label"`
}

func (r *RioDoSulHarvester) FetchReadings(ctx context.Context) ([]domain.RiverReading, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("riodosul: erro ao criar request: %w", err)
	}

	// Cabeçalhos exigidos pelo backend da Asthon / Cloudflare
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "https://defesacivil.riodosul.sc.gov.br")
	req.Header.Set("Referer", "https://defesacivil.riodosul.sc.gov.br/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.6.2 Safari/605.1.15")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("riodosul: falha de conexao: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("riodosul: status HTTP inesperado: %d", resp.StatusCode)
	}

	var stations []asthonStationDTO
	if err := json.NewDecoder(resp.Body).Decode(&stations); err != nil {
		return nil, fmt.Errorf("riodosul: falha ao deserializar json da Asthon: %w", err)
	}

	var readings []domain.RiverReading

	for _, s := range stations {
		if s.LastReadingAt == nil || s.LevelM == nil {
			continue
		}

		// Filtra as estações do leito principal do rio (cota de cheia/transbordamento de 6.5m)
		// Isso descarta barragens (cota > 20m) e ribeirões secundários (cota < 4m)
		if s.OverflowCotaM == nil || *s.OverflowCotaM != 6.5 {
			continue
		}

		observedAt, err := parseAsthonDate(*s.LastReadingAt)
		if err != nil {
			observedAt = time.Now().UTC()
		}

		rain := 0.0
		if s.RainfallMM != nil {
			rain = *s.RainfallMM
		}

		readings = append(readings, domain.RiverReading{
			StationID:        s.StationID,
			City:             "Rio do Sul",
			RiverName:        "Rio Itajaí do Sul / Açu",
			WaterLevelMeters: *s.LevelM,
			Rainfall1hMM:     rain,
			ObservedAt:       observedAt,
			Source:           r.Name(),
		})
	}

	if len(readings) == 0 {
		return nil, fmt.Errorf("riodosul: nenhuma estacao de rio principal encontrada no payload")
	}

	return readings, nil
}

func parseAsthonDate(raw string) (time.Time, error) {
	layouts := []string{
		"2006-01-02 15:04:05-07",
		"2006-01-02 15:04:05.999999-07",
		time.RFC3339,
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, raw); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("data nao reconhecida: %s", raw)
}
