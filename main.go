package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"

	"github.com/prometheus/exporter-toolkit/web"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/robinmarechal/asteriskk_exporter/collector"
)

var (
	listenAddress         = kingpin.Flag("web.listen-address", "The address to listen on for HTTP requests.").Default(":9795").String()
	asteriskPath          = kingpin.Flag("asterisk.path", "Path to Asterisk binary").Default("/usr/sbin/asterisk").String()
	prefix                = kingpin.Flag("metrics.prefix", "Prefix of exposed metrics").Default("asterisk").String()
	metricsPath           = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	enableExporterMetrics = kingpin.Flag("web.enable-exporter-metrics", "Include metrics about the exporter itself (process_*, go_*).").Default("false").Bool()
	enablePromHttpMetrics = kingpin.Flag("web.enable-promhttp-metrics", "Include metrics about the http server itself (promhttp_*)").Default("true").Bool()
	maxRequests           = kingpin.Flag("web.max-requests", "Maximum number of parallel scrape requests. Use 0 to disable.").Default("40").Int()
)

func main() {
	os.Exit(run())
}

func run() int {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("asteriskk_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "starting asteriskk_exporter", "version", version.Info())
	level.Info(logger).Log("build_context", version.BuildContext())

	http.Handle(*metricsPath, newHandler(*enableExporterMetrics, *enablePromHttpMetrics, *maxRequests, &logger))

	handleHealth(&logger)
	handleRoot(&logger)

	return startServer(&logger)
}

func startServer(logger *log.Logger) int {
	srv := &http.Server{Addr: *listenAddress}
	srvc := make(chan struct{})
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	go func() {
		level.Info(*logger).Log("msg", "Listening on address", "address", *listenAddress)
		if err := web.ListenAndServe(srv, "", *logger); err != http.ErrServerClosed {
			level.Error(*logger).Log("msg", "Error starting HTTP server", "err", err)
			close(srvc)
		}
	}()

	for {
		select {
		case <-term:
			level.Info(*logger).Log("msg", "Received SIGTERM, exiting gracefully...")
			return 0
		case <-srvc:
			return 1
		}
	}
}

func handleHealth(logger *log.Logger) {
	http.HandleFunc("/-/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy"))
	})
}

func handleRoot(logger *log.Logger) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html>
    <head><title>Asterisk Exporter</title></head>
    	<body>
    		<h1>Asterisk Exporter</h1>
    		<p><a href="` + *metricsPath + `">Metrics</a></p>
    		<p><a href="config">Configuration</a></p>
		</body>
    </html>`))
	})

}

// handler wraps an unfiltered http.Handler but uses a filtered handler,
// created on the fly, if filtering is requested. Create instances with
// newHandler.
type handler struct {
	unfilteredHandler http.Handler
	// exporterMetricsRegistry is a separate registry for the metrics about
	// the exporter itself.
	exporterMetricsRegistry *prometheus.Registry
	includeExporterMetrics  bool
	includePromHttpMetrics  bool
	maxRequests             int
	logger                  *log.Logger
}

func newHandler(includeExporterMetrics bool, enablePromHttpMetrics bool, maxRequests int, logger *log.Logger) *handler {
	h := &handler{
		exporterMetricsRegistry: prometheus.NewRegistry(),
		includeExporterMetrics:  includeExporterMetrics,
		includePromHttpMetrics:  enablePromHttpMetrics,
		maxRequests:             maxRequests,
		logger:                  logger,
	}

	if h.includeExporterMetrics {
		h.exporterMetricsRegistry.MustRegister(
			prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
			prometheus.NewGoCollector(),
		)
	}

	if innerHandler, err := h.innerHandler(logger); err != nil {
		panic(fmt.Sprintf("Couldn't create metrics handler: %s", err))
	} else {
		h.unfilteredHandler = innerHandler
	}
	return h
}

// ServeHTTP implements http.Handler.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// No filters, use the prepared unfiltered handler.
	h.unfilteredHandler.ServeHTTP(w, r)
}

// innerHandler is used to create both the one unfiltered http.Handler to be
// wrapped by the outer handler and also the filtered handlers created on the
// fly. The former is accomplished by calling innerHandler without any arguments
// (in which case it will log all the collectors enabled via command-line
// flags).
func (h *handler) innerHandler(logger *log.Logger) (http.Handler, error) {
	nc := collector.NewAsteriskCollector(*prefix, asteriskPath, logger)

	r := prometheus.NewRegistry()
	r.MustRegister(version.NewCollector("asteriskk_exporter"))
	if err := r.Register(nc); err != nil {
		return nil, fmt.Errorf("couldn't register asteriskk collector: %s", err)
	}
	handler := promhttp.HandlerFor(
		prometheus.Gatherers{h.exporterMetricsRegistry, r},
		promhttp.HandlerOpts{
			ErrorHandling:       promhttp.ContinueOnError,
			MaxRequestsInFlight: h.maxRequests,
			Registry:            h.exporterMetricsRegistry,
		},
	)
	if h.includePromHttpMetrics {
		// Note that we have to use h.exporterMetricsRegistry here to
		// use the same promhttp metrics for all expositions.
		handler = promhttp.InstrumentMetricHandler(
			h.exporterMetricsRegistry, handler,
		)
	}
	return handler, nil
}
