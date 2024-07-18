package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type SigV4LambdaPayload struct {
	Method   string      `json:"method"`
	Endpoint string      `json:"endpoint"`
	Payload  interface{} `json:"payload"`
}

type SigV4Auth struct {
	client *lambda.Client
}

func NewSigV4(cfg aws.Config) *SigV4Auth {
	return &SigV4Auth{
		client: lambda.NewFromConfig(cfg),
	}
}

func (s SigV4Auth) SignedHeaders(ctx context.Context, payload SigV4LambdaPayload) (map[string]string, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("fail to marshal: %v", err)
	}

	fmt.Println(string(payloadBytes))

	result, err := s.client.Invoke(ctx, &lambda.InvokeInput{
		FunctionName: aws.String("github-api-gw-token"),
		Payload:      payloadBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("fail to invoke lambda: %v", err)
	}

	var headers map[string]string

	err = json.Unmarshal(result.Payload, &headers)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal: %v", err)
	}

	return headers, nil
}
