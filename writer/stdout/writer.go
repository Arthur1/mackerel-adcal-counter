package stdout

import (
	"context"
	"fmt"
	"math"

	"github.com/mackerelio/mackerel-client-go"
)

type Writer struct{}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) Write(service string, values []*mackerel.MetricValue) error {
	return w.WriteWithContext(context.Background(), service, values)
}

func (w *Writer) WriteWithContext(ctx context.Context, service string, values []*mackerel.MetricValue) error {
	fmt.Printf("service: %s\n", service)
	for _, value := range values {
		var v float64
		switch i := value.Value.(type) {
		case int32:
			v = float64(i)
		case uint32:
			v = float64(i)
		case float32:
			v = float64(i)
		case int64:
			v = float64(i)
		case uint64:
			v = float64(i)
		case float64:
			v = i
		default:
			v = math.NaN()
		}
		fmt.Printf("%s\t%f\t%d\n", value.Name, v, value.Time)
	}
	return nil
}
