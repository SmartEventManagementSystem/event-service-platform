package runtime

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"backend/event-service-platform/app/internal/config"
	"backend/event-service-platform/app/pkg/db"
	"backend/event-service-platform/app/pkg/redis"
)

type Clients struct {
}

type Resource struct {
	Config     config.ApplicationConfig
	Logger     *zap.Logger
	DB         *db.DB
	Redis      redis.Redis
	HttpClient *resty.Client
	SqsClient  *sqs.Client
	Clients    Clients
}
