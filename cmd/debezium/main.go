package main

import (
	"context"
	"fmt"
	"github.com/private/energy-shares-cdc/internal/auth"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/go-resty/resty/v2"
)

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

	sigv4Client := auth.NewSigV4(cfg)

	// Get Auth Headers
	headers, err := sigv4Client.SignedHeaders(context.Background(), auth.SigV4LambdaPayload{
		Method:   "GET",
		Endpoint: "https://dbzm-common-gw.illuminarean.com",
	})

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
