package lq

import (
	"context"
	"time"

	"github.com/internetarchive/Zeno/internal/pkg/log"
)

func vacuumer() {
	logger := log.NewFieldedLogger(&log.Fields{
		"component": "lq.vacuumer",
	})

	// Create a context to manage goroutines
	ctx, cancel := context.WithCancel(globalLQ.ctx)
	defer cancel()

	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Debug("Context cancelled, exiting")
			globalLQ.wg.Done()
			return
		case <-ticker.C:
			logger.Info("Vacuuming")
			_, err := globalLQ.client.dbWrite.Exec("VACUUM;")
			if err != nil {
				logger.Error("Vacuuming failed", err)
			}
			logger.Info("Vacuuming complete")
		}
	}

}
