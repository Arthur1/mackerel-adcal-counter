package writer

import (
	"context"

	"github.com/mackerelio/mackerel-client-go"
)

type Writer interface {
	Write(string, []*mackerel.MetricValue) error
	WriteWithContext(context.Context, string, []*mackerel.MetricValue) error
}
