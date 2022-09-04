package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	shinobiexporter "github.com/pavelgopanenko/shinobi_exporter"
	"github.com/pavelgopanenko/shinobi_exporter/config"
	"github.com/pavelgopanenko/shinobi_exporter/metric"
	"github.com/pavelgopanenko/shinobi_exporter/shinobi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app := &cli.App{
		Name:     "shinobi_exporter",
		Usage:    "Shinobi exporter to Prometheus",
		Version:  shinobiexporter.Version,
		Action:   run,
		Flags:    config.CLIFlags(),
		HideHelp: true,
	}

	err := app.RunContext(ctx, os.Args)
	if err != nil {
		fmt.Println(err) // nolint: forbidigo

		defer os.Exit(1)
	}
}

func run(cmd *cli.Context) error {
	ctx := cmd.Context

	server, err := shinobi.NewServerDefault(
		cmd.String(config.EndpointFlag),
		cmd.String(config.TokenFlag),
		cmd.Bool(config.Insecure),
	)
	if err != nil {
		return err
	}

	const v13 = shinobi.Group("v13")

	if _, err := server.Keys(ctx, v13); err != nil {
		return fmt.Errorf("shinobi communication error: %w", err)
	}

	prometheus.MustRegister(
		metric.NewServerCollector(ctx, server, v13),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Shinobi Exporter</title></head>
			<body>
			<h1>Shinobi Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:    cmd.String(config.WebListenAddress),
		Handler: mux,
	}

	errCh := make(chan error)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	return srv.Shutdown(ctx)
}
