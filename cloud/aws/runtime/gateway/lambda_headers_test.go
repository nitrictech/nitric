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

package gateway

import (
	"testing"
)

func TestLambdaHeaders_MultipleSetCookies(t *testing.T) {
	lh := newLambdaHeaders()
	lh.Add("Set-Cookie", "session=abc; Path=/; HttpOnly")
	lh.Add("Set-Cookie", "csrf=xyz; Path=/; Secure")

	if len(lh.Cookies) != 2 {
		t.Fatalf("expected 2 cookies, got %d", len(lh.Cookies))
	}
	if lh.Cookies[0] != "session=abc; Path=/; HttpOnly" {
		t.Errorf("unexpected cookie[0]: %s", lh.Cookies[0])
	}
	if lh.Cookies[1] != "csrf=xyz; Path=/; Secure" {
		t.Errorf("unexpected cookie[1]: %s", lh.Cookies[1])
	}
	if _, ok := lh.Headers["set-cookie"]; ok {
		t.Error("Set-Cookie should not appear in Headers map")
	}
}

func TestLambdaHeaders_MultiValueNonCookieCommaFolded(t *testing.T) {
	lh := newLambdaHeaders()
	lh.Add("Link", "</style.css>; rel=preload")
	lh.Add("Link", "</script.js>; rel=preload")

	val, ok := lh.Headers["link"]
	if !ok {
		t.Fatal("expected 'link' in Headers")
	}
	expected := "</style.css>; rel=preload, </script.js>; rel=preload"
	if val != expected {
		t.Errorf("expected %q, got %q", expected, val)
	}
}

func TestLambdaHeaders_SingleValueUnchanged(t *testing.T) {
	lh := newLambdaHeaders()
	lh.Add("Content-Type", "application/json")

	val, ok := lh.Headers["content-type"]
	if !ok {
		t.Fatal("expected 'content-type' in Headers")
	}
	if val != "application/json" {
		t.Errorf("expected 'application/json', got %q", val)
	}
}

func TestLambdaHeaders_MixedCaseNormalisedToLowercase(t *testing.T) {
	lh := newLambdaHeaders()
	lh.Add("X-Custom-Header", "value1")
	lh.Add("x-custom-header", "value2")
	lh.Add("X-CUSTOM-HEADER", "value3")

	val, ok := lh.Headers["x-custom-header"]
	if !ok {
		t.Fatal("expected 'x-custom-header' in Headers")
	}
	expected := "value1, value2, value3"
	if val != expected {
		t.Errorf("expected %q, got %q", expected, val)
	}
	// Should only have one key
	if len(lh.Headers) != 1 {
		t.Errorf("expected 1 header key, got %d", len(lh.Headers))
	}
}

func TestLambdaHeaders_SetCookieCaseInsensitive(t *testing.T) {
	lh := newLambdaHeaders()
	lh.Add("SET-COOKIE", "a=1")
	lh.Add("set-cookie", "b=2")
	lh.Add("Set-Cookie", "c=3")

	if len(lh.Cookies) != 3 {
		t.Fatalf("expected 3 cookies, got %d", len(lh.Cookies))
	}
	if len(lh.Headers) != 0 {
		t.Errorf("expected no headers, got %d", len(lh.Headers))
	}
}

func TestLambdaHeaders_EmptyInitialState(t *testing.T) {
	lh := newLambdaHeaders()

	if len(lh.Headers) != 0 {
		t.Errorf("expected empty headers, got %d", len(lh.Headers))
	}
	if len(lh.Cookies) != 0 {
		t.Errorf("expected empty cookies, got %d", len(lh.Cookies))
	}
}
