// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package telemetry

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"
)

type TelemetryConfig struct {
	Config string
	Uri    string
}

type TelemetryConfigArgs struct {
	MetricName           string
	TraceName            string
	MetricExporterConfig string
	TraceExporterConfig  string
	Extensions           []string
	TraceSampling        int
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

	return &TelemetryConfig{
		Config: builder.String(),
		Uri:    otelCollectorVersion,
	}, nil
}
