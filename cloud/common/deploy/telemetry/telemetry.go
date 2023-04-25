package telemetry

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"
)

type TelemetryConfig struct {
	Config string
	Uri string
}

type TelemetryConfigArgs struct {
	MetricName           string
	TraceName            string
	MetricExporterConfig string
	TraceExporterConfig  string
	Extensions           []string
	TraceSampling int
}

//go:embed otel-collector.yaml
var otelCollectorTemplate string

const (
	// https://github.com/open-telemetry/opentelemetry-collector/releases
	otelVersion = "v0.66.0"
)

func NewTelemetryConfig(config *TelemetryConfigArgs) (*TelemetryConfig, error) {
	yamlTemplate := template.Must(template.New("otelconfig").Parse(otelCollectorTemplate))
	builder := &strings.Builder{}

	err := yamlTemplate.Execute(builder, config)
	if err != nil {
		return nil, err
	}

	otelCollectorVersion := fmt.Sprintf(
		"https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/%s/otelcol-contrib_%s_linux_amd64.tar.gz",
		otelVersion, strings.TrimSpace(strings.TrimLeft(otelVersion, "v")),
	)

	return &TelemetryConfig {
		Config: builder.String(),
		Uri: otelCollectorVersion,
	}, nil
}