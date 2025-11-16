package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vinsec/nexttrace_exporter/collector"
	"github.com/vinsec/nexttrace_exporter/config"
	"github.com/vinsec/nexttrace_exporter/executor"
)

var (
	configFile = kingpin.Flag(
		"config.file",
		"Path to configuration file.",
	).Default("config.yml").String()

	listenAddress = kingpin.Flag(
		"web.listen-address",
		"Address to listen on for web interface and telemetry.",
	).Default("localhost:9101").String()

	metricsPath = kingpin.Flag(
		"web.telemetry-path",
		"Path under which to expose metrics.",
	).Default("/metrics").String()

	nexttraceBinary = kingpin.Flag(
		"nexttrace.binary",
		"Path to nexttrace binary.",
	).Default("nexttrace").String()

	nexttraceTimeout = kingpin.Flag(
		"nexttrace.timeout",
		"Timeout for nexttrace execution.",
	).Default("2m").Duration()

	logLevel = kingpin.Flag(
		"log.level",
		"Log level (debug, info, warn, error).",
	).Default("info").String()
)

type Server struct {
	executor  *executor.Executor
	collector *collector.Collector
	registry  *prometheus.Registry
	config    *config.Config
	logger    *slog.Logger
	ctx       context.Context
	cancel    context.CancelFunc
}

func main() {
	kingpin.Version("0.1.0")
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	// Setup logger
	logger := setupLogger(*logLevel)

	// Load initial configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger.Info("Configuration loaded successfully",
		"targets", len(cfg.Targets),
		"config_file", *configFile)

	// Create context for managing goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize server
	server := &Server{
		executor: executor.NewExecutor(*nexttraceBinary, *nexttraceTimeout, logger),
		config:   cfg,
		logger:   logger,
		ctx:      ctx,
		cancel:   cancel,
		registry: prometheus.NewRegistry(),
	}

	// Create collector
	server.collector = collector.NewCollector(server.executor, cfg.Targets, logger)

	// Register collector
	server.registry.MustRegister(server.collector)

	// Use config file values if command-line flags are at default
	if *listenAddress == "localhost:9101" && cfg.Server.ListenAddress != "" {
		listenAddress = &cfg.Server.ListenAddress
	}
	if *metricsPath == "/metrics" && cfg.Server.MetricsPath != "" {
		metricsPath = &cfg.Server.MetricsPath
	}

	// Start executor
	server.executor.Start(ctx, cfg.Targets)

	// Setup signal handling
	go server.handleSignals()

	// Start HTTP server
	if err := server.startHTTPServer(); err != nil {
		logger.Error("HTTP server failed", "error", err)
		os.Exit(1)
	}
}

func setupLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	return slog.New(handler)
}

func (s *Server) startHTTPServer() error {
	mux := http.NewServeMux()

	// Metrics endpoint
	mux.Handle(*metricsPath, promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{}))

	// Health check endpoint
	mux.HandleFunc("/-/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK\n")
	})

	// Reload endpoint
	mux.HandleFunc("/-/reload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := s.reload(); err != nil {
			s.logger.Error("Failed to reload configuration", "error", err)
			http.Error(w, fmt.Sprintf("Failed to reload: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Configuration reloaded successfully\n")
	})

	// Landing page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<html>
<head><title>NextTrace Exporter</title></head>
<body>
<h1>NextTrace Exporter</h1>
<p><a href="%s">Metrics</a></p>
<h2>Configured Targets</h2>
<ul>
`, *metricsPath)

		for _, target := range s.config.Targets {
			fmt.Fprintf(w, "<li><strong>%s</strong> (%s) - Interval: %s, Max Hops: %d</li>\n",
				target.Name, target.Host, target.Interval, target.MaxHops)
		}

		fmt.Fprintf(w, `</ul>
<h2>Endpoints</h2>
<ul>
<li><a href="/-/healthy">Health Check</a></li>
<li><a href="/-/reload">Reload Configuration</a> (POST)</li>
</ul>
</body>
</html>`)
	})

	s.logger.Info("Starting HTTP server",
		"address", *listenAddress,
		"metrics_path", *metricsPath)

	srv := &http.Server{
		Addr:         *listenAddress,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return srv.ListenAndServe()
}

func (s *Server) reload() error {
	s.logger.Info("Reloading configuration", "config_file", *configFile)

	// Load new configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Update server state
	s.config = cfg

	// Reload executor with new targets
	s.executor.Reload(s.ctx, cfg.Targets)

	// Update collector targets
	s.collector.UpdateTargets(cfg.Targets)

	s.logger.Info("Configuration reloaded successfully", "targets", len(cfg.Targets))

	return nil
}

func (s *Server) handleSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	for sig := range sigChan {
		switch sig {
		case syscall.SIGHUP:
			s.logger.Info("Received SIGHUP, reloading configuration")
			if err := s.reload(); err != nil {
				s.logger.Error("Failed to reload configuration", "error", err)
			}
		case syscall.SIGINT, syscall.SIGTERM:
			s.logger.Info("Received shutdown signal", "signal", sig)
			s.shutdown()
			os.Exit(0)
		}
	}
}

func (s *Server) shutdown() {
	s.logger.Info("Shutting down...")
	s.cancel()
	s.executor.Stop()
	s.logger.Info("Shutdown complete")
}
