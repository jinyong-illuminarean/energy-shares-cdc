package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/go-resty/resty/v2"
	"os"
	"strings"
)

type LambdaPayload struct {
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
	Payload  string `json:"payload"`
}

func getAuthHeaders() (map[string]string, error) {
	// AWS 설정 로드
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return nil, fmt.Errorf("AWS 설정 로드 실패: %v", err)
	}

	// Lambda 클라이언트 생성
	client := lambda.NewFromConfig(cfg)

	// Lambda 함수에 전달할 페이로드 준비
	payload := LambdaPayload{
		Method:   "GET",
		Endpoint: "https://dbzm-common-gw.illuminarean.com",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("페이로드 마샬링 실패: %v", err)
	}

	// Lambda 함수 호출
	result, err := client.Invoke(context.TODO(), &lambda.InvokeInput{
		FunctionName: aws.String("github-api-gw-token"),
		Payload:      payloadBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("Lambda 함수 호출 실패: %v", err)
	}

	// 응답 파싱
	var headers map[string]string
	err = json.Unmarshal(result.Payload, &headers)
	if err != nil {
		return nil, fmt.Errorf("응답 언마샬링 실패: %v", err)
	}

	// 헤더 키 정제 (대시를 언더스코어로 변경)
	sanitizedHeaders := make(map[string]string)
	for k, v := range headers {
		sanitizedHeaders[strings.ReplaceAll(k, "-", "_")] = v
	}

	return sanitizedHeaders, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: debezium-updater <environment>")
		fmt.Println("Available environments: dev, qa, stage, prod")
		os.Exit(1)
	}

	fmt.Printf("Debezium Updater[%s]\n", os.Args[1])

	client := resty.New()

	headers, err := getAuthHeaders()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	response, err := client.R().
		SetHeader("X-Amz-Date", headers["X-Amz-Date"]).
		SetHeader("X-Amz-Security-Token", headers["X-Amz-Security-Token"]).
		SetHeader("Authorization", headers["Authorization"]).
		Get("https://dbzm-common-gw.illuminarean.com")
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	fmt.Println("Response Body:")
	fmt.Println(string(response.Body()))
}
