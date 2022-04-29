package cache_service

import (
	"github.com/piupuer/go-helper/pkg/tracing"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer(tracing.Cache)
