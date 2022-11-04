package mackerel

import (
	"context"

	"github.com/mackerelio/mackerel-client-go"
)

type Writer struct {
	client *mackerel.Client
}

func NewWriter(apiKey string) (*Writer, error) {
	client := mackerel.NewClient(apiKey)
	client.UserAgent = "mackerel-adcal-counter"
	return &Writer{
		client: client,
	}, nil
}

func (e *Writer) Write(service string, metrics []*mackerel.MetricValue) error {
	return e.WriteWithContext(context.Background(), service, metrics)
}

func (e *Writer) WriteWithContext(ctx context.Context, service string, values []*mackerel.MetricValue) error {
	return e.client.PostServiceMetricValues(service, values)
}
