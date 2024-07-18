package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/go-resty/resty/v2"
)

type LambdaPayload struct {
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
}

func getAuthHeaders(cfg aws.Config) (map[string]string, error) {
	client := lambda.NewFromConfig(cfg)

	payload := LambdaPayload{
		Method:   "GET",
		Endpoint: "https://dbzm-common-gw.illuminarean.com",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("페이로드 마샬링 실패: %v", err)
	}

	result, err := client.Invoke(context.TODO(), &lambda.InvokeInput{
		FunctionName: aws.String("github-api-gw-token"),
		Payload:      payloadBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("Lambda 함수 호출 실패: %v", err)
	}

	var headers map[string]string
	err = json.Unmarshal(result.Payload, &headers)
	if err != nil {
		return nil, fmt.Errorf("응답 언마샬링 실패: %v", err)
	}

	return headers, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: debezium-updater <environment>")
		fmt.Println("Available environments: dev, qa, stage, prod")
		os.Exit(1)
	}

	fmt.Printf("Debezium Updater[%s]\n", os.Args[1])

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		fmt.Printf("AWS 설정 로드 실패: %v\n", err)
		os.Exit(1)
	}

	// Get Auth Headers
	headers, err := getAuthHeaders(cfg)
	if err != nil {
		fmt.Printf("인증 헤더 가져오기 실패: %v\n", err)
		os.Exit(1)
	}

	// Print headers (similar to setting environment variables)
	for k, v := range headers {
		fmt.Printf("%s=%s\n", k, v)
	}

	// Call API
	client := resty.New()
	response, err := client.R().
		SetHeader("X-Amz-Date", headers["X-Amz-Date"]).
		SetHeader("X-Amz-Security-Token", headers["X-Amz-Security-Token"]).
		SetHeader("Authorization", headers["Authorization"]).
		Get("https://dbzm-common-gw.illuminarean.com")

	if err != nil {
		fmt.Printf("API 호출 실패: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Response Status Code:", response.StatusCode())
	fmt.Println("Response Body:")
	fmt.Println(string(response.Body()))
}
