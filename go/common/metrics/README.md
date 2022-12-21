This package contains code related to the node's host metrics system.

The metrics package heavily uses geths own metric package.
- golang Package [here](https://github.com/ethereum/go-ethereum/tree/master/metrics)
- Documentation [here](https://geth.ethereum.org/docs/interface/metrics)


## Add metrics

Instrumentation can be added and used as the following example shows:

```go

import gethmetrics "github.com/ethereum/go-ethereum/metrics"

type SomeStruct struct {
        ...
        someMeter gethmetrics.Meter
        someGauge gethmetrics.Gauge
}

// registry should be passed down from the parent package
func New(registry gethmetrics.Registry) *SomeStruct {
	return &SomeStruct {
            someMeter : gethmetrics.NewRegisteredMeter("some/metric/path", registry),
            someGauge : gethmetrics.NewRegisteredGauge("some/metric/path", registry), 
        }
}

func (s *SomeStruct) Do()  {
	//some other code
        s.someMeter.Update(Some_Value)
	s.someGauge.Inc(1)
}

```