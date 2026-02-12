/*
Copyright © 2026 Amanda Hager Lopes de Andrade Katz amandahla@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
)

const MAX_ELAPSED_TIME = 10 * time.Second

// Call makes an HTTP request to the specified path with the given method and configuration.
// If retry is true, it will retry the request with exponential backoff in case of failure.
// If payload is provided and the method is POST, it will include the payload in the request body.
// It returns the response body as a byte slice or an error if the request fails.
// Call makes an HTTP request to the specified path with the given method and configuration.
// Accepts context.Context for cancellation and timeout propagation.
func Call(ctx context.Context, config Config, path string, method string, payload []byte, retry bool) ([]byte, error) {
	var output []byte
	client := &http.Client{Timeout: time.Duration(config.Timeout) * time.Second}
	synapseURL := fmt.Sprintf("%s%s", config.BaseURL, path)
	var sendBody io.Reader
	if method == http.MethodPost && payload != nil {
		sendBody = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, synapseURL, sendBody)
	if err != nil {
		return output, fmt.Errorf("Request to %s failed: %v", synapseURL, err)
	}
	req.Header.Set("Authorization", "Bearer "+config.AccessToken)
	if retry {
		return callWithRetry(ctx, client, req, synapseURL)
	}
	resp, err := client.Do(req)
	if err != nil {
		return output, fmt.Errorf("Request to %s failed: %v", synapseURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return output, fmt.Errorf("Request to %s returned unexpected status: %v", synapseURL, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return output, fmt.Errorf("Failed to read response body from %s: %v", synapseURL, err)
	}

	return body, nil
}

// callWithRetry performs an HTTP request with exponential backoff, accepting context.Context.
func callWithRetry(ctx context.Context, client *http.Client, req *http.Request, synapseURL string) ([]byte, error) {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = MAX_ELAPSED_TIME

	var output []byte
	err := backoff.Retry(func() error {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			return fmt.Errorf("Request to %s failed: %v", synapseURL, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("Request to %s returned unexpected status: %v", synapseURL, resp.Status)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Failed to read response body from %s: %v", synapseURL, err)
		}
		output = body
		return nil
	}, b)
	return output, err
}
