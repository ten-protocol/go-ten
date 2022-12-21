package p2p

import gethmetrics "github.com/ethereum/go-ethereum/metrics"

// perHostMetrics is a map that returns metrics per string
// errors per host or events per type are common usages
type perHostMetrics struct {
	gaugeMap map[string]gethmetrics.Gauge
	instName string
	reg      gethmetrics.Registry
}

// newperHostMetricMap creates a new p2p instrument for host based tracking
func newperHostMetricMap(registry gethmetrics.Registry, instrumentName string) *perHostMetrics {
	return &perHostMetrics{
		instName: instrumentName,
		gaugeMap: map[string]gethmetrics.Gauge{},
		reg:      registry,
	}
}

// inc increments the instrument at a given host with a specific value
// creates one entry it one does not exist
// gethmetrics is thread-safe so no concurrency is handled at the map level
func (g *perHostMetrics) inc(str string, val int64) {
	if _, ok := g.gaugeMap[str]; !ok {
		g.gaugeMap[str] = gethmetrics.NewRegisteredGauge(g.instName, g.reg)
	}
	g.gaugeMap[str].Inc(val)
}

// totals returns the sum of all entries
func (g *perHostMetrics) totals() int64 {
	total := int64(0)
	for _, gauge := range g.gaugeMap {
		total += gauge.Value()
	}
	return total
}
