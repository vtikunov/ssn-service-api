package metrics

import (
	"github.com/ozonmp/ssn-service-api/internal/retranslator/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var eventsCountInPool prometheus.Gauge

// InitMetrics - инициализирует метрики.
func InitMetrics(cfg *config.Config) {
	eventsCountInPool = promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: cfg.Metrics.Subsystem,
		Name:      "events_count_in_pool",
		Help:      "Total events in retranslator",
	})
}

// AddEventsCountInPool - увеличивает значение счетчика событий в пуле.
func AddEventsCountInPool(c uint) {
	if eventsCountInPool == nil {
		return
	}

	eventsCountInPool.Add(float64(c))
}

// SubEventsCountInPool - уменьшает значение счетчика событий в пуле.
func SubEventsCountInPool(c uint) {
	if eventsCountInPool == nil {
		return
	}

	eventsCountInPool.Sub(float64(c))
}
