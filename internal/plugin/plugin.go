package plugin

import (
	"context"

	"github.com/selectel/fluent-bit-cloudwatch-input-plugin/internal/model"
	"github.com/selectel/fluent-bit-cloudwatch-input-plugin/internal/plugin/adapter"
)

type Plugin struct {
	region        string
	endpoint      string
	logGroupName  string
	logStreamName string

	cloudwatch adapter.Cloudwatch
	storage    adapter.Storage
}

func NewPlugin(
	region, endpoint, logGroupName, logStreamName string,
	cloudwatch adapter.Cloudwatch,
	storage adapter.Storage,
) *Plugin {
	return &Plugin{
		region:        region,
		endpoint:      endpoint,
		logGroupName:  logGroupName,
		logStreamName: logStreamName,
		cloudwatch:    cloudwatch,
		storage:       storage,
	}
}

func (p *Plugin) GetLogEvents(ctx context.Context, nextToken string) ([]model.Event, string, error) {
	return p.cloudwatch.GetLogEvents(ctx, p.logGroupName, p.logStreamName, nextToken)
}

func (p *Plugin) GetNextToken(ctx context.Context) (string, error) {
	return p.storage.GetNextToken(ctx, p.region, p.logGroupName, p.logStreamName)
}

func (p *Plugin) SetNextToken(ctx context.Context, nextToken string) error {
	return p.storage.SetNextToken(ctx, p.region, p.logGroupName, p.logStreamName, nextToken)
}
