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

package cors

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/imdario/mergo"
	"github.com/valyala/fasthttp"

	base_http "github.com/nitrictech/nitric/cloud/common/runtime/gateway"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/worker/pool"
)

func GetCorsConfig(vals *v1.ApiCorsDefinition) (*v1.ApiCorsDefinition, error) {
	defaultVal := &v1.ApiCorsDefinition{
		AllowCredentials: false,
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		ExposeHeaders:    []string{},
		MaxAge:           300,
	}

	if err := mergo.Merge(defaultVal, vals, mergo.WithOverride); err != nil {
		return nil, err
	}

	return defaultVal, nil
}

// Used for GCP and Local CORs with fasthttp headers
func GetCorsHeaders(config *v1.ApiCorsDefinition) (*map[string]string, error) {
	corsHeaders := map[string]string{}

	corsConfig, err := GetCorsConfig(config)
	if err != nil {
		return nil, err
	}

	corsHeaders["Access-Control-Allow-Credentials"] = strconv.FormatBool(corsConfig.GetAllowCredentials())

	if len(corsConfig.GetAllowOrigins()) > 0 {
		corsHeaders["Access-Control-Allow-Origin"] = strings.Join(corsConfig.GetAllowOrigins(), ",")
	}

	if len(corsConfig.GetAllowMethods()) > 0 {
		corsHeaders["Access-Control-Allow-Methods"] = strings.Join(corsConfig.GetAllowMethods(), ",")
	}

	if len(corsConfig.GetAllowHeaders()) > 0 {
		corsHeaders["Access-Control-Allow-Headers"] = strings.Join(corsConfig.GetAllowHeaders(), ",")
	}

	if len(corsConfig.GetExposeHeaders()) > 0 {
		corsHeaders["Access-Control-Expose-Headers"] = strings.Join(corsConfig.GetExposeHeaders(), ",")
	}

	corsHeaders["Access-Control-Max-Age"] = strconv.FormatInt(int64(corsConfig.GetMaxAge()), 10)

	return &corsHeaders, nil
}

func GetEnvKey(name string) string {
	return fmt.Sprintf("NITRIC_CORS_%s", strings.ToUpper(name))
}

func CreateCorsMiddleware(cache map[string]map[string]string) base_http.HttpMiddleware {
	return func(rc *fasthttp.RequestCtx, wp pool.WorkerPool) bool {
		api := string(rc.Request.Header.Peek("x-nitric-api"))
		method := string(rc.Request.Header.Method())

		if cache[api] != nil {
			applyCorsHeaders(rc, cache[api])
			return true
		}

		corsHeaders, err := getCorsHeadersForAPI(api)
		if err != nil {
			if method == "OPTIONS" {
				rc.Response.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			}

			return true
		}

		cache[api] = corsHeaders

		applyCorsHeaders(rc, corsHeaders)

		return true
	}
}

func getCorsHeadersForAPI(name string) (map[string]string, error) {
	env := os.Getenv(GetEnvKey(name))

	if env == "" {
		return nil, fmt.Errorf("no cors env var found for api %s", name)
	}

	headers := make(map[string]string)

	// Unmarshal the JSON string into the map
	err := json.Unmarshal([]byte(env), &headers)
	if err != nil {
		return nil, err
	}

	return headers, nil
}

func applyCorsHeaders(rc *fasthttp.RequestCtx, corsHeaders map[string]string) {
	for k, v := range corsHeaders {
		rc.Response.Header.Add(k, v)
	}
}
