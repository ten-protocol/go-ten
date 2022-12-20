package p2p

import gethmetrics "github.com/ethereum/go-ethereum/metrics"

// perStringGaugeMap is a map that returns metrics per string
// errors per host or events per type are common usages
type perStringGaugeMap struct {
	gaugeMap map[string]gethmetrics.Gauge
	instName string
	reg      gethmetrics.Registry
}

// newPerStringGaugeMap creates a new instrument
func newPerStringGaugeMap(registry gethmetrics.Registry, instrumentName string) *perStringGaugeMap {
	return &perStringGaugeMap{
		instName: instrumentName,
		gaugeMap: map[string]gethmetrics.Gauge{},
		reg:      registry,
	}
}

// inc increments the instrument at a given host with a specific value
// creates one entry it one does not exist
// gethmetrics is thread-safe so no concurrency is handled at the map level
func (g *perStringGaugeMap) inc(str string, val int64) {
	if _, ok := g.gaugeMap[str]; !ok {
		g.gaugeMap[str] = gethmetrics.NewRegisteredGauge(g.instName, g.reg)
	}
	g.gaugeMap[str].Inc(val)
}

// totals returns the sum of all entries
func (g *perStringGaugeMap) totals() int64 {
	total := int64(0)
	for _, gauge := range g.gaugeMap {
		total += gauge.Value()
	}
	return total
}
