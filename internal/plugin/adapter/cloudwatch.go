package adapter

import (
	"context"

	"github.com/selectel/fluent-bit-cloudwatch-input-plugin/internal/model"
)

type Cloudwatch interface {
	GetLogEvents(ctx context.Context, logGroupName, logStreamName, nextToken string) ([]model.Event, string, error)
}
