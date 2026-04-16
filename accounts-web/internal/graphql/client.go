package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	url        string
	httpClient *http.Client
}

type PingResult struct {
	Message string `json:"message"`
	Service string `json:"service"`
}

type graphQLRequest struct {
	Query string `json:"query"`
}

type graphQLError struct {
	Message string `json:"message"`
}

type pingResponseEnvelope struct {
	Data struct {
		Ping PingResult `json:"ping"`
	} `json:"data"`
	Errors []graphQLError `json:"errors"`
}

func New(url string, timeout time.Duration) *Client {
	return &Client{
		url: url,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Ping(ctx context.Context) (PingResult, error) {
	payload, err := json.Marshal(graphQLRequest{
		Query: `query Ping { ping { message service } }`,
	})
	if err != nil {
		return PingResult{}, fmt.Errorf("marshal graphql request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(payload))
	if err != nil {
		return PingResult{}, fmt.Errorf("build graphql request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PingResult{}, fmt.Errorf("execute graphql request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PingResult{}, fmt.Errorf("unexpected graphql status: %s", resp.Status)
	}

	var envelope pingResponseEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return PingResult{}, fmt.Errorf("decode graphql response: %w", err)
	}

	if len(envelope.Errors) > 0 {
		return PingResult{}, fmt.Errorf("graphql error: %s", envelope.Errors[0].Message)
	}

	return envelope.Data.Ping, nil
}
