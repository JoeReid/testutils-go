package containers

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/ory/dockertest/v3"
)

func DynamoDB(t *testing.T) *dynamo.DB {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("failed to construct docker pool: %v", err)
	}

	if err := pool.Client.Ping(); err != nil {
		t.Fatalf("failed to connect to docker: %v", err)
	}

	dynamoContainer, err := pool.Run("amazon/dynamodb-local", "latest", nil)
	if err != nil {
		t.Fatalf("failed to start dynamo container: %v", err)
	}

	ses, err := session.NewSession()
	if err != nil {
		t.Fatalf("failed to create aws session: %v", err)
	}

	cfg := &aws.Config{
		Region:      aws.String("us-east-1"),
		Endpoint:    aws.String("http://localhost:" + dynamoContainer.GetPort("8000/tcp")),
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", ""),
	}

	if err := pool.Retry(func() error {
		db := dynamo.New(ses, cfg)

		_, err := db.ListTables().All()
		return err
	}); err != nil {
		t.Fatalf("failed to connect to dynamo container: %v", err)
	}

	t.Cleanup(func() {
		if err := pool.Purge(dynamoContainer); err != nil {
			t.Fatalf("failed to purge dynamo container: %v", err)
		}
	})

	return dynamo.New(ses, cfg)
}
