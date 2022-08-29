package initialize

import (
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/log"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func Tracer() {
	if !global.Conf.Tracer.Enable {
		log.WithContext(ctx).Info("tracer is not enabled")
		return
	}
	driverOpts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(global.Conf.Tracer.Endpoint),
		otlptracehttp.WithHeaders(global.Conf.Tracer.Headers),
	}
	if global.Conf.Tracer.Insecure {
		driverOpts = append(driverOpts, otlptracehttp.WithInsecure())
	}
	// create exporters
	driver := otlptracehttp.NewClient(driverOpts...)

	exporter, err := otlptrace.New(ctx, driver)

	if err != nil {
		panic(errors.Wrap(err, "initialize tracer failed"))
	}

	// A custom ID Generator to generate traceIDs that conform to
	// AWS X-Ray traceID format
	idg := xray.NewIDGenerator()

	// create options
	opts := []trace.TracerProviderOption{
		trace.WithSyncer(exporter),
		trace.WithIDGenerator(idg),
		trace.WithSampler(trace.AlwaysSample()),
	}

	srcRes := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(global.ProName),
	)

	opts = append(opts, trace.WithResource(srcRes))

	// create tracer provider
	tp := trace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})
	global.Tracer = tp
	log.WithContext(ctx).Info("initialize tracer success")
	return
}
