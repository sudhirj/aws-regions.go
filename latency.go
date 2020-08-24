package awsregions

import (
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"
)

type measurement struct {
	region  string
	latency time.Duration
}

type LatencyChecker struct {
	regions   []string
	latencies map[string]time.Duration
	sync.RWMutex
}

func NewLatencyChecker(regions ...string) LatencyChecker {
	return LatencyChecker{regions: regions, latencies: map[string]time.Duration{}}
}

func (lc *LatencyChecker) Start() {
	lc.latencies = make(map[string]time.Duration)
	receiver := make(chan measurement)

	for _, region := range lc.regions {
		go keepMeasuring(region, receiver)
	}

	for {
		select {
		case m := <-receiver:
			lc.Lock()
			lc.latencies[m.region] = (m.latency + lc.latencies[m.region]) / 2
			sort.Slice(lc.regions, func(i, j int) bool {
				return lc.latencies[lc.regions[i]] < lc.latencies[lc.regions[j]]
			})
			lc.Unlock()
		}
	}
}

func (lc *LatencyChecker) FastestRegion() string {
	return lc.SortedRegions()[0]
}

func (lc *LatencyChecker) SortedRegions() []string {
	dup := make([]string, len(lc.regions))
	lc.RLock()
	copy(dup, lc.regions)
	lc.RUnlock()

	return dup
}

func (lc *LatencyChecker) Latencies() map[string]time.Duration {
	dup := map[string]time.Duration{}
	for k, v := range lc.latencies {
		dup[k] = v
	}
	return dup
}

func keepMeasuring(region string, receiver chan measurement) {
	receiver <- measure(region)

	for {
		receiver <- measure(region)
		time.Sleep(time.Duration(rand.Int63n(30)) * time.Second)
	}
}

func measure(region string) measurement {
	start := time.Now()
	resp, _ := http.Get("https://dynamodb." + region + ".amazonaws.com")
	defer resp.Body.Close()

	end := time.Now()
	latency := end.Sub(start)

	return measurement{region: region, latency: latency}
}
