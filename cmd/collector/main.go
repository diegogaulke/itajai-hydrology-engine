package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/seu-usuario/itajai-hydrology-engine/internal/harvester"
)

func main() {
	log.Println("Starting Itajai Hydrology Engine Ingestion Service...")

	httpClient := &http.Client{
		Timeout: 12 * time.Second,
	}

	harvesters := []harvester.Harvester{
		harvester.NewAlertaBluHarvester(httpClient),
		harvester.NewRioDoSulHarvester(httpClient),
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	// Coleta inicial no boot
	runIngestionCycle(ctx, harvesters)

	for {
		select {
		case <-ctx.Done():
			log.Println("Gracefully shutting down hydrology collector...")
			return
		case <-ticker.C:
			runIngestionCycle(ctx, harvesters)
		}
	}
}

func runIngestionCycle(ctx context.Context, harvesters []harvester.Harvester) {
	log.Println("--- Starting ingestion cycle across river stations ---")
	var wg sync.WaitGroup

	for _, h := range harvesters {
		wg.Add(1)
		go func(harv harvester.Harvester) {
			defer wg.Done()

			subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
			defer cancel()

			readings, err := harv.FetchReadings(subCtx)
			if err != nil {
				log.Printf("[WARN] Error collecting from %s: %v", harv.Name(), err)
				return
			}

			for _, r := range readings {
				log.Printf("[SUCCESS] %s", r.String())
			}
		}(h)
	}

	wg.Wait()
	log.Println("--- Ingestion cycle completed ---")
}
