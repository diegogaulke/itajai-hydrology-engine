# Itajaí Hydrology Engine

> Telemetry ingestion and predictive hydrological monitoring for the Rio Itajaí-Açu watershed (Santa Catarina, Brazil).

## Overview

The Itajaí River Valley is historically prone to flash floods and severe flooding events during El Niño cycles. Water level surges in **Blumenau** are heavily influenced by meteorological events and dam management upstream in the Upper Valley (**Rio do Sul**, Taió, and Ituporanga) and Middle Valley (**Ibirama**, Apiúna, and Indaial).

`itajai-hydrology-engine` is a resilient backend service written in **Go (Golang)** designed to continuously ingest, normalize, and compute trend vectors from dispersed telemetry endpoints provided by municipal civil defense agencies and national water authorities (**ANA — *Agência Nacional de Águas e Saneamento Básico***).

## Architecture Highlights

* **Fault-Tolerant Scraping / Ingestion:** Decoupled `Harvester` interfaces with isolated HTTP timeouts and circuit-breaking strategies. Public agency portals often experience downtime during major storms; an error on one station adapter never halts the ingestion of the rest of the watershed.
* **Unified Domain Model:** Regardless of upstream formatting differences (metric strings, Brazilian date encodings, variable sensor units), data points are standardized into canonical UTC time-series records.
* **Concurrent Polling:** Goroutine-based fan-out polling architecture that executes across all upstream observation points simultaneously.

## Station Coverage

| Station ID | Location | Watershed Segment | Sensor Target | Primary Upstream Influence |
|---|---|---|---|---|
| `BLU-CENTRO` | Blumenau | Lower/Middle Valley | Rio Itajaí-Açu | River level gauge & local rainfall |
| `RDS-TITOBUSS`| Rio do Sul | Upper Valley | Ponte Dom Tito Buss | Peak crest ~16-24h prior to Blumenau impact |

## Getting Started

### Prerequisites

* Go 1.22+

### Run Locally

```bash
# Clone the repository
git clone [https://github.com/seu-usuario/itajai-hydrology-engine.git](https://github.com/seu-usuario/itajai-hydrology-engine.git)
cd itajai-hydrology-engine

# Run the ingestion collector
go run cmd/collector/main.go
