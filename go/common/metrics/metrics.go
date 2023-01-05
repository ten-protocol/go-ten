package metrics

import (
	"fmt"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/ethereum/go-ethereum/metrics/exp"

	gethlog "github.com/ethereum/go-ethereum/log"
	gethmetrics "github.com/ethereum/go-ethereum/metrics"
)

var (
	_collectProcessMetricsRefreshDuration = 3 * time.Second
	_threadCreateProfile                  = pprof.Lookup("threadcreate")
)

// Service provides the metrics for the host
// it registers the gethmetrics Registry
// and handles the metrics server
type Service struct {
	registry gethmetrics.Registry
	port     uint
	logger   gethlog.Logger
}

func New(enabled bool, port uint, logger gethlog.Logger) *Service {
	gethmetrics.Enabled = enabled
	return &Service{
		registry: gethmetrics.NewRegistry(),
		port:     port,
		logger:   logger,
	}
}

// Start starts the metrics server
func (m *Service) Start() {
	// metrics not enabled
	if !gethmetrics.Enabled {
		return
	}

	// start the process collection metric on it's own thread
	go m.CollectProcessMetrics()

	// starts the metric server
	address := fmt.Sprintf("%s:%d", "0.0.0.0", m.port)
	m.logger.Info("HTTP Metric server started at %s", address)
	// TODO re-write this http server so to have a stop method
	exp.Setup(address)
}

// Registry returns the registry for the metrics service
func (m *Service) Registry() gethmetrics.Registry {
	return m.registry
}

func (m *Service) Stop() {
	// TODO re-write this http server so to have a stop method
}

// CollectProcessMetrics collect process and system metrics
// this method is an adapted copy of eth/metrics/metrics.go@CollectProcessMetrics method
func (m *Service) CollectProcessMetrics() {
	refreshFreq := int64(_collectProcessMetricsRefreshDuration / time.Second)

	// Create the various data collectors
	cpuStats := make([]*gethmetrics.CPUStats, 2)
	memstats := make([]*runtime.MemStats, 2)
	diskstats := make([]*gethmetrics.DiskStats, 2)
	for i := 0; i < len(memstats); i++ {
		cpuStats[i] = new(gethmetrics.CPUStats)
		memstats[i] = new(runtime.MemStats)
		diskstats[i] = new(gethmetrics.DiskStats)
	}
	// Define the various metrics to collect
	var (
		cpuSysLoad    = gethmetrics.GetOrRegisterGauge("system/cpu/sysload", m.registry)
		cpuSysWait    = gethmetrics.GetOrRegisterGauge("system/cpu/syswait", m.registry)
		cpuProcLoad   = gethmetrics.GetOrRegisterGauge("system/cpu/procload", m.registry)
		cpuThreads    = gethmetrics.GetOrRegisterGauge("system/cpu/threads", m.registry)
		cpuGoroutines = gethmetrics.GetOrRegisterGauge("system/cpu/goroutines", m.registry)

		memPauses = gethmetrics.GetOrRegisterMeter("system/memory/pauses", m.registry)
		memAllocs = gethmetrics.GetOrRegisterMeter("system/memory/allocs", m.registry)
		memFrees  = gethmetrics.GetOrRegisterMeter("system/memory/frees", m.registry)
		memHeld   = gethmetrics.GetOrRegisterGauge("system/memory/held", m.registry)
		memUsed   = gethmetrics.GetOrRegisterGauge("system/memory/used", m.registry)

		diskReads             = gethmetrics.GetOrRegisterMeter("system/disk/readcount", m.registry)
		diskReadBytes         = gethmetrics.GetOrRegisterMeter("system/disk/readdata", m.registry)
		diskReadBytesCounter  = gethmetrics.GetOrRegisterCounter("system/disk/readbytes", m.registry)
		diskWrites            = gethmetrics.GetOrRegisterMeter("system/disk/writecount", m.registry)
		diskWriteBytes        = gethmetrics.GetOrRegisterMeter("system/disk/writedata", m.registry)
		diskWriteBytesCounter = gethmetrics.GetOrRegisterCounter("system/disk/writebytes", m.registry)
	)
	// Iterate loading the different stats and updating the meters
	for i := 1; ; i++ {
		location1 := i % 2
		location2 := (i - 1) % 2

		gethmetrics.ReadCPUStats(cpuStats[location1])
		cpuSysLoad.Update((cpuStats[location1].GlobalTime - cpuStats[location2].GlobalTime) / refreshFreq)
		cpuSysWait.Update((cpuStats[location1].GlobalWait - cpuStats[location2].GlobalWait) / refreshFreq)
		cpuProcLoad.Update((cpuStats[location1].LocalTime - cpuStats[location2].LocalTime) / refreshFreq)
		cpuThreads.Update(int64(_threadCreateProfile.Count()))
		cpuGoroutines.Update(int64(runtime.NumGoroutine()))

		runtime.ReadMemStats(memstats[location1])
		memPauses.Mark(int64(memstats[location1].PauseTotalNs - memstats[location2].PauseTotalNs))
		memAllocs.Mark(int64(memstats[location1].Mallocs - memstats[location2].Mallocs))
		memFrees.Mark(int64(memstats[location1].Frees - memstats[location2].Frees))
		memHeld.Update(int64(memstats[location1].HeapSys - memstats[location1].HeapReleased))
		memUsed.Update(int64(memstats[location1].Alloc))

		if gethmetrics.ReadDiskStats(diskstats[location1]) == nil {
			diskReads.Mark(diskstats[location1].ReadCount - diskstats[location2].ReadCount)
			diskReadBytes.Mark(diskstats[location1].ReadBytes - diskstats[location2].ReadBytes)
			diskWrites.Mark(diskstats[location1].WriteCount - diskstats[location2].WriteCount)
			diskWriteBytes.Mark(diskstats[location1].WriteBytes - diskstats[location2].WriteBytes)

			diskReadBytesCounter.Inc(diskstats[location1].ReadBytes - diskstats[location2].ReadBytes)
			diskWriteBytesCounter.Inc(diskstats[location1].WriteBytes - diskstats[location2].WriteBytes)
		}
		time.Sleep(_collectProcessMetricsRefreshDuration)
	}
}
