package device

import (
	"context"
	"time"

	"github.com/hashicorp/nomad/plugins/device"
)

// doStats is the long running goroutine that streams device statistics
func (d *SkeletonDevicePlugin) doStats(ctx context.Context, stats chan<- *device.StatsResponse, interval time.Duration) {
	defer close(stats)

	// Create a timer that will fire immediately for the first detection
	ticker := time.NewTimer(0)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ticker.Reset(interval)
		}

		d.writeStatsToChannel(stats, time.Now())
	}
}

// writeStatsToChannel collects device stats, partitions devices into
// device groups, and sends the data over the provided channel.
func (d *SkeletonDevicePlugin) writeStatsToChannel(stats chan<- *device.StatsResponse, timestamp time.Time) {
	stats <- &device.StatsResponse{}
}
