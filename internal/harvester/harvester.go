package harvester

import (
	"context"
	"github.com/seu-usuario/itajai-hydrology-engine/internal/domain"
)

// Harvester define o contrato para coletores de órgãos de monitoramento.
type Harvester interface {
	Name() string
	FetchReadings(ctx context.Context) ([]domain.RiverReading, error)
}
