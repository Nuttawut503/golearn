package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func FindItem(ctx context.Context) {
	tr := otel.Tracer("component-query")
	_, span := tr.Start(ctx, "query")
	defer span.End()
	span.SetAttributes(attribute.Key("search.id").String("a1b2c3d4"))
	// do something
	time.Sleep(time.Millisecond * 500)
}

func Filter(ctx context.Context) {
	tr := otel.Tracer("component-filter")
	_, span := tr.Start(ctx, "filter")
	defer span.End()
	// do something
	time.Sleep(time.Millisecond * 100)
	// code example when there is an error
	// span.SetStatus(codes.Error, "fail to filter items")
	// span.RecordError(errors.New("item format is not correct"))
	time.Sleep(time.Millisecond * 200)
}

func SetupHandler(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		tr := otel.Tracer("component-http-request")
		ctx, span := tr.Start(context.Background(), "http-request")
		FindItem(ctx)
		Filter(ctx)
		defer span.End()
	})
}
