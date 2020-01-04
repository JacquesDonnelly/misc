package hello

import (
	"fmt"
	"net/http"
	"os"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	ext "github.com/opentracing/opentracing-go/ext"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
	zipkin "github.com/uber/jaeger-client-go/zipkin"
)

type greeter struct {
	version string
	greetee string
}

var tracer opentracing.Tracer
var propagator zipkin.Propagator

func init() {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			CollectorEndpoint:   os.Getenv("JAEGER_ENDPOINT"),
		},
	}
	tracer, _, _ = cfg.New(
		"hello",
		config.Logger(jaeger.StdLogger),
	)
	propagator = zipkin.NewZipkinB3HTTPHeaderPropagator()
}

func (g *greeter) rootHandler(w http.ResponseWriter, r *http.Request) {
	// read in the B3 open tracing headers
	spanCtx, err := propagator.Extract(opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		fmt.Printf("Span context error: %s\n", err)
	}

	// write a new span for this handler
	serverSpan := tracer.StartSpan("root_hander", ext.RPCServerOption(spanCtx))
	defer serverSpan.Finish()

	fmt.Fprintf(w, "Hello, %s from %s!\n", g.greetee, r.RemoteAddr)
	fmt.Fprintf(w, "Welcome to %s\n", r.URL.Path)
	fmt.Fprintf(w, "Version %s\n", g.version)
}

func serveAdmin() {
	adminRouter := http.NewServeMux()

	adminRouter.HandleFunc("/healthz", healthzHandler)
	adminRouter.HandleFunc("/ready", readyHandler)
	adminRouter.HandleFunc("/set-unhealthy", setUnhealthyHandler)
	adminRouter.HandleFunc("/set-healthy", setHealthyHandler)
	adminRouter.HandleFunc("/set-unready", setUnreadyHandler)
	adminRouter.HandleFunc("/set-ready", setReadyHandler)

	adminServer := http.Server{
		Addr:    ":8081",
		Handler: adminRouter,
	}
	fmt.Println("Serving health endpoints on :8081")
	adminServer.ListenAndServe()
}

// Run starts the HTTP server
func Run(version string, name string) {
	g := &greeter{
		version: version,
		greetee: name,
	}

	go serveAdmin()
	go delayReady()

	helloRouter := http.NewServeMux()
	helloRouter.HandleFunc("/", g.rootHandler)
	helloServer := http.Server{
		Addr:    ":8080",
		Handler: helloRouter,
	}
	fmt.Println("Serving on :8080")
	helloServer.ListenAndServe()
}
