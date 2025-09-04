package server

import (
	"github.com/datpp/go-kratos-based-template/packages/types"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/google/wire"
	"go.opentelemetry.io/otel/exporters/prometheus"
	ometrics "go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, NewMetrics)

// Metrics holds the metrics instances
type Metrics struct {
	Seconds  ometrics.Float64Histogram
	Requests ometrics.Int64Counter
}

// NewMetrics creates new metrics with the given service name
func NewMetrics(app types.AppInfo) (*Metrics, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	meter := provider.Meter(*app.Name)

	requests, err := metrics.DefaultRequestsCounter(meter, metrics.DefaultServerRequestsCounterName)
	if err != nil {
		return nil, err
	}

	seconds, err := metrics.DefaultSecondsHistogram(meter, metrics.DefaultServerSecondsHistogramName)
	if err != nil {
		return nil, err
	}

	return &Metrics{
		Seconds:  seconds,
		Requests: requests,
	}, nil
}
