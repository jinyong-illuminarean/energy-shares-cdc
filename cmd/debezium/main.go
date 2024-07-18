package main

import (
	"context"
	"fmt"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/private/energy-shares-cdc/internal/auth"
	"github.com/private/energy-shares-cdc/internal/config"
	"github.com/private/energy-shares-cdc/internal/debezium"
	"log"
	"os"
	"time"
)

func main() {
	// 1분 타임아웃
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if len(os.Args) != 3 {
		log.Fatal("Usage: ./debezium <environment> <db_password>")
	}

	fmt.Println("pass : ", os.Args[2])

	cfg, err := config.LoadConfig(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	awscfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion("ap-northeast-2"))
	if err != nil {
		log.Fatalf("Error loading AWS config: %v", err)
	}

	signV4Client := auth.NewSigV4(awscfg)

	client := debezium.New(signV4Client, cfg.DebeziumURL)

	// 서버 헬스체크
	if err := client.HealthCheck(ctx); err != nil {
		log.Fatalf("Server health check failed: %v", err)
	}

	fmt.Println("Server is healthy")

	// connect 리스트 조회
	connectors, err := client.ListConnectors(ctx)
	if err != nil {
		log.Fatalf("Failed to list connectors: %v", err)
	}

	fmt.Println("Connectors:", connectors)

	// 해당 환경에 맞는 커넥터가 존재하는지 확인
	connectorName := fmt.Sprintf("mysql.%s.cdc", os.Args[1])
	exists := false
	for _, c := range connectors {
		if c == connectorName {
			exists = true
			break
		}
	}

	fmt.Println("Connector exists:", exists)

	if exists {
		// 커넥터가 존재하면 상세 조회 후 설정 업데이트
		currentConfig, err := client.GetConnectorConfig(ctx, connectorName)
		if err != nil {
			log.Fatalf("Failed to get connector config: %v", err)
		}

		fmt.Println("Current connector config:", currentConfig)

		if !cfg.Compare(currentConfig) {
			if err := client.UpdateConnectorConfig(ctx, connectorName, cfg.ConnectorConfig); err != nil {
				log.Fatalf("Failed to update connector config: %v", err)
			}
			fmt.Println("Connector config updated successfully")
		} else {
			fmt.Println("Connector config is up to date")
		}
	} else {
		// 커넥터가 존재하지 않으면 새로 생성
		if err := client.CreateConnector(ctx, connectorName, cfg.ConnectorConfig); err != nil {
			log.Fatalf("Failed to create connector: %v", err)
		}
		fmt.Println("Connector created successfully")
	}
}
