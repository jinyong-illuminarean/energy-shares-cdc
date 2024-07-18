package debezium

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/private/energy-shares-cdc/internal/auth"
)

type DebeziumClient struct {
	client    *resty.Client
	sigV4Auth *auth.SigV4Auth
	baseURL   string
}

func New(sigV4Auth *auth.SigV4Auth, baseURL string) *DebeziumClient {
	return &DebeziumClient{
		client:    resty.New(),
		sigV4Auth: sigV4Auth,
		baseURL:   baseURL,
	}
}

func (c *DebeziumClient) HealthCheck(ctx context.Context) error {
	lambdaPayload := auth.SigV4LambdaPayload{
		Method:   "GET",
		Endpoint: c.baseURL,
	}

	payload, err := json.Marshal(lambdaPayload)
	if err != nil {
		return fmt.Errorf("fail to marshal: %v", err)
	}

	headers, err := c.sigV4Auth.SignedHeaders(ctx, payload)

	if err != nil {
		return fmt.Errorf("failed to get auth headers: %v", err)
	}

	resp, err := c.client.R().
		SetHeaders(headers).
		Get(c.baseURL)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode())
	}
	return nil
}

func (c *DebeziumClient) ListConnectors(ctx context.Context) ([]string, error) {
	lambdaPayload := auth.SigV4LambdaPayload{
		Method:   "GET",
		Endpoint: c.baseURL + "/connectors",
	}

	payload, err := json.Marshal(lambdaPayload)
	if err != nil {
		return nil, fmt.Errorf("fail to marshal: %v", err)
	}

	headers, err := c.sigV4Auth.SignedHeaders(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth headers: %v", err)
	}

	resp, err := c.client.R().
		SetHeaders(headers).
		SetResult([]string{}).
		Get(c.baseURL + "/connectors")

	if err != nil {
		return nil, err
	}

	return *resp.Result().(*[]string), nil
}

func (c *DebeziumClient) GetConnectorConfig(ctx context.Context, name string) (map[string]interface{}, error) {
	lambdaPayload := auth.SigV4LambdaPayload{
		Method:   "GET",
		Endpoint: c.baseURL + "/connectors/" + name + "/config",
	}

	payload, err := json.Marshal(lambdaPayload)
	if err != nil {
		return nil, fmt.Errorf("fail to marshal: %v", err)
	}

	headers, err := c.sigV4Auth.SignedHeaders(ctx, payload)

	if err != nil {
		return nil, fmt.Errorf("failed to get auth headers: %v", err)
	}

	resp, err := c.client.R().
		SetHeaders(headers).
		SetResult(map[string]interface{}{}).
		Get(c.baseURL + "/connectors/" + name + "/config")

	if err != nil {
		return nil, err
	}

	return *resp.Result().(*map[string]interface{}), nil
}

func (c *DebeziumClient) UpdateConnectorConfig(ctx context.Context, name string, config map[string]interface{}) error {
	body, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("fail to marshal: %v", err)
	}

	lambdaPayload := auth.SigV4LambdaPayload{
		Method:   "PUT",
		Endpoint: c.baseURL + "/connectors/" + name + "/config",
		Payload:  body,
	}

	payload, err := json.Marshal(lambdaPayload)
	if err != nil {
		return fmt.Errorf("fail to marshal: %v", err)
	}

	headers, err := c.sigV4Auth.SignedHeaders(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to get auth headers: %v", err)
	}

	resp, err := c.client.R().
		SetHeaders(headers).
		SetBody(body).
		SetHeader("Content-Type", "application/json").
		Put(c.baseURL + "/connectors/" + name + "/config")

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to update connector config, status: %d, message: %s", resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

func (c *DebeziumClient) CreateConnector(ctx context.Context, name string, config map[string]interface{}) error {
	body := map[string]interface{}{
		"name":   name,
		"config": config,
	}

	lambdaPayload := auth.SigV4LambdaPayload{
		Method:   "PUT",
		Endpoint: c.baseURL + "/connectors",
		Payload:  body,
	}

	payload, err := json.Marshal(lambdaPayload)
	if err != nil {
		return fmt.Errorf("fail to marshal: %v", err)
	}

	headers, err := c.sigV4Auth.SignedHeaders(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to get auth headers: %v", err)
	}

	resp, err := c.client.R().
		SetHeaders(headers).
		SetBody(payload).
		Post(c.baseURL + "/connectors")

	if err != nil {
		return err
	}

	if resp.StatusCode() != 201 {
		return fmt.Errorf("failed to create connector, status: %d", resp.StatusCode())
	}
	return nil
}
