package metrics

import gethmetrics "github.com/ethereum/go-ethereum/metrics"

// PerStringGaugeMap is a map that returns metrics per string
// errors per host or events per type are common usages
type PerStringGaugeMap struct {
	gaugeMap map[string]gethmetrics.Gauge
	instName string
	reg      gethmetrics.Registry
}

// NewPerStringGaugeMap creates a new instrument
func NewPerStringGaugeMap(registry gethmetrics.Registry, instrumentName string) *PerStringGaugeMap {
	return &PerStringGaugeMap{
		instName: instrumentName,
		gaugeMap: map[string]gethmetrics.Gauge{},
		reg:      registry,
	}
}

// Inc increments the instrument at a given host with a specific value
// creates one entry it one does not exist
// gethmetrics is thread-safe so no concurrency is handled at the map level
func (g *PerStringGaugeMap) Inc(str string, val int64) {
	if _, ok := g.gaugeMap[str]; !ok {
		g.gaugeMap[str] = gethmetrics.NewRegisteredGauge(g.instName, g.reg)
	}
	g.gaugeMap[str].Inc(val)
}

// totals returns the sum of all entries
func (g *PerStringGaugeMap) totals() int64 {
	total := int64(0)
	for _, gauge := range g.gaugeMap {
		total += gauge.Value()
	}
	return total
}
