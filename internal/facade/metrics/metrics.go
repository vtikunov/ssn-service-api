package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/ozonmp/ssn-service-api/internal/facade/config"
	"github.com/ozonmp/ssn-service-api/internal/facade/model/subscription"
)

var cudCountTotal *prometheus.CounterVec

// InitMetrics - инициализирует метрики.
func InitMetrics(cfg *config.Config) {
	cudCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Subsystem: cfg.Metrics.Subsystem,
		Name:      "cud_total",
		Help:      "Total CUD",
	}, []string{"type"})
}

// AddCudCountTotal - увеличивает значение счетчика событий.
func AddCudCountTotal(c uint, eventType subscription.EventType) {
	if cudCountTotal == nil {
		return
	}

	cudCountTotal.WithLabelValues(string(eventType)).Add(float64(c))
}
