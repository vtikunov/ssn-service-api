package metrics

import (
	"github.com/ozonmp/ssn-service-api/internal/config"
	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var notFoundErrorsTotal prometheus.Counter
var cudCountTotal *prometheus.CounterVec

// InitMetrics - инициализирует метрики.
func InitMetrics(cfg *config.Config) {
	notFoundErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Subsystem: cfg.Metrics.Subsystem,
		Name:      "not_found_errors_total",
		Help:      "Total not found errors",
	})

	cudCountTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Subsystem: cfg.Metrics.Subsystem,
		Name:      "cud_total",
		Help:      "Total CUD",
	}, []string{"type"})
}

// AddNotFoundErrorsTotal - увеличивает значение счетчика ошибок отсутствия объекта.
func AddNotFoundErrorsTotal(c uint) {
	if notFoundErrorsTotal == nil {
		return
	}

	notFoundErrorsTotal.Add(float64(c))
}

// AddCudCountTotal - увеличивает значение счетчика событий.
func AddCudCountTotal(c uint, eventType subscription.EventType) {
	if cudCountTotal == nil {
		return
	}

	cudCountTotal.WithLabelValues(string(eventType)).Add(float64(c))
}
