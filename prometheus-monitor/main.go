package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

var (
	lemonsKey = attribute.Key("ex.com/lemons")
)

func initMeterServer() {
	config := prometheus.Config{
		DefaultHistogramBoundaries: []float64{0.1, 0.2, 0.5, 1, 2},
	}
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}

	global.SetMeterProvider(exporter.MeterProvider())

	r := gin.Default()
	r.GET("/metrics", gin.WrapH(exporter))
	go func() {
		_ = r.Run(":" + os.Getenv("PORT"))
	}()

	fmt.Println("Prometheus server running on :8080")
}

func main() {
	initMeterServer()

	meter := global.Meter("")

	observerLock := new(sync.RWMutex)
	var value float64
	commonAttrs := []attribute.KeyValue{attribute.String("A", "1"), attribute.String("A", "2")}

	gaugeObserver, err := meter.AsyncFloat64().Gauge("http_request_sample")
	if err != nil {
		log.Panicf("failed to initialize instrument: %v", err)
	}
	_ = meter.RegisterCallback([]instrument.Asynchronous{gaugeObserver}, func(ctx context.Context) {
		// use to capture the current value
		observerLock.RLock()
		value := value
		observerLock.RUnlock()
		gaugeObserver.Observe(ctx, value, commonAttrs...)
	})

	histogram, err := meter.SyncFloat64().Histogram("http_request_duration_seconds")
	if err != nil {
		log.Panicf("failed to initialize instrument: %v", err)
	}

	ctx := context.Background()

	for {
		observerLock.Lock()
		value = rand.Float64() * 3
		histogram.Record(ctx, value, commonAttrs...)
		observerLock.Unlock()
		time.Sleep(1 * time.Second)
	}
}
