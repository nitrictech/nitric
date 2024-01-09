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

package env

// Standard system environment variables
var (
	MAX_WORKERS     = GetEnv("MAX_WORKERS", "300")
	MIN_WORKERS     = GetEnv("MIN_WORKERS", "1")
	SERVICE_ADDRESS = GetEnv("SERVICE_ADDRESS", "127.0.0.1:50051")
)

// Open telemetry environment variables
// FIXME: there are a provider level concern so should be implemented in the provider layer
var (
	OTELCOL_BIN      = GetEnv("OTELCOL_BIN", "/usr/bin/otelcol-contrib")
	OTELCOL_CONFIG   = GetEnv("OTELCOL_CONFIG", "/etc/otelcol/config.yaml")
	MEMBRANE_VERSION = GetEnv("MEMBRANE_VERSION", "unknown")
)
